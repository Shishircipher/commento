package db

import (
	"context"
	"log"
//	"github.com/jackc/pgx/v5/pgxpool"
)
// getComments retrieves all comments for a given post_id
func GetComments(ctx context.Context, postID int) ([]Comment, error) {
        rows, err := DB.Query(ctx, `
                SELECT id, post_id, parent_id, content, author_id 
                FROM Comment 
                WHERE post_id = $1 AND parent_id IS NULL
        `, postID)
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
