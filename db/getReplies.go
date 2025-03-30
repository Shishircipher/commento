package db

import (
	"context"
	"log"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
)
// getComments retrieves all comments for a given post_id

func GetReplies(ctx context.Context, db *pgxpool.Pool, postID, commentID, depthLimit, limit, offset int) ([]Comment, error) {
        rows, err := db.Query(ctx, `
                SELECT c.id, c.post_id, c.parent_id, c.content, c.author_id, cc.depth, u.name, u.picture
                FROM Comment c
		JOIN users u ON c.author_id = u.id
                JOIN CommentClosure cc ON c.id = cc.descendant_id
                WHERE cc.ancestor_id = $1 
                  AND c.post_id = $2 
                  AND cc.depth > 0 
                  AND cc.depth <= $3
                ORDER BY cc.depth ASC
                LIMIT $4 OFFSET $5
        `, commentID, postID, depthLimit, limit, offset)
        if err != nil {
		log.Printf("%v ", err)
                return nil, err
        }
        defer rows.Close()

        var replies []Comment
        for rows.Next() {
                var reply Comment
		var authorName string
		var authorPicture sql.NullString // Handles NULL values
                var depth int
                if err := rows.Scan(&reply.ID, &reply.PostID, &reply.ParentID, &reply.Content, &reply.AuthorID, &depth, &authorName, &authorPicture); err != nil {
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
                //replies = append(replies, reply)
		// Append author details to the response
		replies = append(replies, Comment{
			ID:       reply.ID,
			PostID:   reply.PostID,
			ParentID: reply.ParentID,
			Content:  reply.Content,
			AuthorID: reply.AuthorID,
			Author: User_public{
				Name:    authorName,
				Picture: picture,
			},
		})
        }
        return replies, nil
}
