package db

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
	"fmt"
)


func InsertComment(ctx context.Context, db *pgxpool.Pool, comment NewComment) (int, error) {

        tx, err := db.Begin(ctx)
	if err != nil {
        log.Println("Failed to begin transaction:", err)
        return 0, err
         }

     //   if err != nil {
       //         return 0, err
      //  }
      //  defer tx.Rollback(ctx)

        var authorID int
        log.Println("Checking if user exists:", comment.Author.Email)
        err = tx.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", comment.Author.Email).Scan(&authorID)
	log.Printf("Query users : %v\n", err)
        log.Println("User not found, inserting new user:", comment.Author.Email)

	// Fix: If user does not exist (ID is 0), insert new user
        if err != nil || authorID == 0 {
        	log.Println("User does not exist, inserting new user:", comment.Author.Email)
        	err = tx.QueryRow(ctx, `
            	INSERT INTO users (name, email, picture)
            	VALUES ($1, $2, $3) RETURNING id
        	`, comment.Author.Name, comment.Author.Email, comment.Author.Picture).Scan(&authorID)

                if err != nil {
                log.Printf("Error inserting user: %v\n", err)
                return 0, fmt.Errorf("failed to insert user: %w", err)
                }
         }
        var newCommentID int
        err = tx.QueryRow(ctx, `
                INSERT INTO Comment (post_id, parent_id, content, author_id) 
                VALUES ($1, $2, $3, $4) RETURNING id
        `, comment.PostID, comment.ParentID, comment.Content, authorID).Scan(&newCommentID)
        if err != nil {
                log.Printf("Error inserting comment: %v\n", err)
                return 0, err
	}
        _, err = tx.Exec(ctx, `
                INSERT INTO CommentClosure (ancestor_id, descendant_id, depth)
                VALUES ($1, $1, 0)
        `, newCommentID)
        if err != nil {
                return 0, err
        }

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

        err = tx.Commit(ctx)
        if err != nil {
		log.Println("Transaction commit failed:", err)
                return 0, err
        }



	return newCommentID, nil
}
