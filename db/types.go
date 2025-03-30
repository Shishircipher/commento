package db

// Comment struct
type Comment struct {
	ID       int    `json:"id"`
	PostID   int    `json:"post_id"`
	ParentID *int   `json:"parent_id"`
	Content  string `json:"content"`
	AuthorID int    `json:"author_id"`
	Author   User_public `json:"author"` // Add author details
}

type User struct {
        Name    string `json:"name"`
        Email   string `json:"email"`
        Picture string `json:"picture"`
}

type NewComment struct {
        Author   User  `json:"author"`
        PostID   int   `json:"post_id"`
        ParentID *int  `json:"parent_id"`
        Content  string `json:"content"`
}

type User_public struct {
	Name    string `json:"name"`
        Picture string `json:"picture"`
}
