package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"blog2/db"
	"blog2/models"
)

// PostsHandler handles all post-related HTTP requests
type PostsHandler struct {
	DB *db.DB
}

// NewPostsHandler creates a new PostsHandler
func NewPostsHandler(db *db.DB) *PostsHandler {
	return &PostsHandler{DB: db}
}

// ServeHTTP handles all HTTP requests for posts
func (h *PostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL if present
	path := strings.TrimPrefix(r.URL.Path, "/posts")
	path = strings.TrimSuffix(path, "/")
	
	// Route based on HTTP method and path
	switch {
	case r.Method == http.MethodGet && path == "":
		h.getPosts(w, r)
	case r.Method == http.MethodGet && path != "":
		id, err := strconv.Atoi(path[1:]) // Remove leading slash
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		h.getPost(w, r, id)
	case r.Method == http.MethodPost && path == "":
		h.createPost(w, r)
	case r.Method == http.MethodPut && path != "":
		id, err := strconv.Atoi(path[1:]) // Remove leading slash
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		h.updatePost(w, r, id)
	case r.Method == http.MethodDelete && path != "":
		id, err := strconv.Atoi(path[1:]) // Remove leading slash
		if err != nil {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		h.deletePost(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getPosts returns all posts
func (h *PostsHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.DB.GetPosts()
	if err != nil {
		http.Error(w, "Error retrieving posts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

// getPost returns a single post by ID
func (h *PostsHandler) getPost(w http.ResponseWriter, r *http.Request, id int) {
	post, err := h.DB.GetPost(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving post: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// createPost adds a new post
func (h *PostsHandler) createPost(w http.ResponseWriter, r *http.Request) {
	var newPost models.NewPost
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if newPost.Title == "" || newPost.Content == "" || newPost.CreatedBy == "" {
		http.Error(w, "Title, content, and created_by are required fields", http.StatusBadRequest)
		return
	}
	
	post, err := h.DB.CreatePost(newPost)
	if err != nil {
		http.Error(w, "Error creating post: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

// updatePost modifies an existing post
func (h *PostsHandler) updatePost(w http.ResponseWriter, r *http.Request, id int) {
	var updatePost models.UpdatePost
	if err := json.NewDecoder(r.Body).Decode(&updatePost); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Validate required fields
	if updatePost.Title == "" || updatePost.Content == "" {
		http.Error(w, "Title and content are required fields", http.StatusBadRequest)
		return
	}
	
	post, err := h.DB.UpdatePost(id, updatePost)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating post: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

// deletePost removes a post
func (h *PostsHandler) deletePost(w http.ResponseWriter, r *http.Request, id int) {
	err := h.DB.DeletePost(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting post: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}