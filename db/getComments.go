package db

import (
	"context"
	"log"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)
// getComments retrieves all comments for a given post_id
func GetComments(ctx context.Context, db *pgxpool.Pool, postID, limit, offset int) ([]Comment, error) {
	rows, err := db.Query(ctx, `
		SELECT c.id, c.post_id, c.parent_id, c.content, c.author_id, u.name, u.picture
		FROM Comment c
		JOIN users u ON c.author_id = u.id
		WHERE c.post_id = $1 AND c.parent_id IS NULL
	`, postID)
        if err != nil {
		log.Printf("%v ", err)
                return nil, err
        }
        defer rows.Close()

        var comments []Comment
        for rows.Next() {
		var comment Comment
		var authorName string
		var authorPicture sql.NullString // Handles NULL values
		if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.AuthorID, &authorName, &authorPicture); err != nil {
			return nil, err
		}
		// Convert `sql.NullString` to `string`
		//var picture string
		// Convert `sql.NullString` to `string`
		picture := ""
		if authorPicture.Valid {
			picture = authorPicture.String
		} else {
			picture = "" // Or set a default profile picture URL if needed
		}
        // Append author details to the response
		comments = append(comments, Comment{
			ID:       comment.ID,
			PostID:   comment.PostID,
			ParentID: comment.ParentID,
			Content:  comment.Content,
			AuthorID: comment.AuthorID,
			Author: User_public{
				Name:    authorName,
				Picture: picture,
			},
		})
	}
	return comments, nil
}
