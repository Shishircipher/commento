package db

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)
// getComments retrieves all comments for a given post_id

func getReplies(ctx context.Context, postID, db *pgxpool.Pool, commentID, depthLimit, limit, offset int) ([]Comment, error) {
        rows, err := db.Query(ctx, `
                SELECT c.id, c.post_id, c.parent_id, c.content, c.author_id, cc.depth
                FROM Comment c
                JOIN CommentClosure cc ON c.id = cc.descendant_id
                WHERE cc.ancestor_id = $1 
                  AND c.post_id = $2 
                  AND cc.depth > 0 
                  AND cc.depth <= $3
                ORDER BY cc.depth ASC
                LIMIT $4 OFFSET $5
        `, commentID, postID, depthLimit, limit, offset)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

        var replies []Comment
        for rows.Next() {
                var reply Comment
                var depth int
                if err := rows.Scan(&reply.ID, &reply.PostID, &reply.ParentID, &reply.Content, &reply.AuthorID, &depth); err != nil {
                        return nil, err
                }
                replies = append(replies, reply)
        }
        return replies, nil
}
