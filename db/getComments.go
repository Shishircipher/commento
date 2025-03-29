package db

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)
// getComments retrieves all comments for a given post_id
func GetComments(ctx context.Context, db *pgxpool.Pool, postID, limit, offset int) ([]Comment, error) {
        rows, err := db.Query(ctx, `
                SELECT id, post_id, parent_id, content, author_id 
                FROM Comment 
                WHERE post_id = $1 AND parent_id IS NULL
		ORDER BY id ASC
                LIMIT $2 OFFSET $3
        `, postID, limit, offset)
        if err != nil {
		log.Printf("%v ", err)
                return nil, err
        }
        defer rows.Close()

        var comments []Comment
        for rows.Next() {
                var comment Comment
                if err := rows.Scan(&comment.ID, &comment.PostID, &comment.ParentID, &comment.Content, &comment.AuthorID); err != nil {
                        log.Printf("%v ", err)
			return nil, err

                }
                comments = append(comments, comment)
        }
        return comments, nil
}
