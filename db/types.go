package db

// Comment struct
type Comment struct {
        ID       int    `json:"id"`
        PostID   int    `json:"post_id"`
        ParentID *int   `json:"parent_id"`
        Content  string `json:"content"`
        AuthorID int    `json:"author_id"`
}
