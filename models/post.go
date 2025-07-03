package models

import (
	"time"
)

// Post represents a blog post in the system
type Post struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	DateCreated time.Time `json:"date_created"`
	CreatedBy   string    `json:"created_by"`
}

// NewPost is used when creating a post (ID and DateCreated are handled by the database)
type NewPost struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedBy string `json:"created_by"`
}

// UpdatePost is used when updating a post
type UpdatePost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
