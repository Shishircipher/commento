package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

var db *pgxpool.Pool

// Initialize database connection
func initDB() {
	databaseURL := os.Getenv("COMMENTO_DB_URL") 
	var err error
	db, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
}

// Comment struct
type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	ParentID *int   `json:"parent_id"`
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
}

// insertComment inserts a new comment and updates CommentClosure
func insertComment(ctx context.Context, comment Comment) (int, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // Rollback if any error occurs

	var newCommentID int
	err = tx.QueryRow(ctx, `
		INSERT INTO Comment (post_id, parent_id, content, author_id) 
		VALUES ($1, $2, $3, $4) RETURNING id
	`, comment.PostID, comment.ParentID, comment.Content, comment.AuthorID).Scan(&newCommentID)
	if err != nil {
		return 0, err
	}

	// Insert self-reference in CommentClosure
	_, err = tx.Exec(ctx, `
		INSERT INTO CommentClosure (ancestor_id, descendant_id, depth) 
		VALUES ($1, $1, 0)
	`, newCommentID)
	if err != nil {
		return 0, err
	}

	// Insert ancestor paths if the comment has a parent
	if comment.ParentID != nil {
		_, err = tx.Exec(ctx, `
			INSERT INTO CommentClosure (ancestor_id, descendant_id, depth) 
			SELECT ancestor_id, $1, depth + 1
			FROM CommentClosure
			WHERE descendant_id = $2
		`, newCommentID, *comment.ParentID)
		if err != nil {
			return 0, err
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return newCommentID, nil
}

// getComments retrieves all comments for a given post_id
func getComments(ctx context.Context, postID int) ([]Comment, error) {
	rows, err := db.Query(ctx, `
		SELECT id, post_id, parent_id, content, author_id 
		FROM Comment 
		WHERE post_id = $1 AND parent_id IS NULL
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.AuthorID); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// getReplies retrieves all replies to a specific comment (post_id and comment_id)
func getReplies(ctx context.Context, postID, commentID int) ([]Comment, error) {
	rows, err := db.Query(ctx, `
		SELECT c.id, c.post_id, c.parent_id, c.content, c.author_id
		FROM Comment c
		JOIN CommentClosure cc ON c.id = cc.descendant_id
		WHERE cc.ancestor_id = $1 AND c.post_id = $2 AND cc.depth > 0
		ORDER BY cc.depth ASC
	`, commentID, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Comment
	for rows.Next() {
		var reply Comment
		if err := rows.Scan(&reply.ID, &reply.PostID, &reply.ParentID, &reply.Content, &reply.AuthorID); err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}

// handleInsertComment - HTTP handler for inserting a comment
func handleInsertComment(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	commentID, err := insertComment(ctx, comment)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting comment: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":    "Comment added successfully",
		"comment_id": commentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetComments - HTTP handler for retrieving comments by postID
func handleGetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	comments, err := getComments(ctx, postID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving comments: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// handleGetReplies - HTTP handler for retrieving replies of a comment
func handleGetReplies(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	commentID, err2 := strconv.Atoi(vars["commentID"])
	if err != nil || err2 != nil {
		http.Error(w, "Invalid post ID or comment ID", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	replies, err := getReplies(ctx, postID, commentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving replies: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(replies)
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/comments", handleInsertComment).Methods("POST")
	r.HandleFunc("/comments/{postID}", handleGetComments).Methods("GET")   // Get comments of a post
	r.HandleFunc("/comments/{postID}/{commentID}", handleGetReplies).Methods("GET") // Get replies to a comment

	// Enable CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"*"}, // Allow all origins. Change this to your frontend URL for security.
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)
	fmt.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}

