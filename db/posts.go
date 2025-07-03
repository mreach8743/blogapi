package db

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"blog2/models"
)

var (
	ErrNotFound = errors.New("post not found")
)

// DB represents a database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// GetPosts retrieves all posts from the database
func (db *DB) GetPosts() ([]models.Post, error) {
	rows, err := db.Query(`
		SELECT id, title, content, date_created, created_by 
		FROM posts 
		ORDER BY date_created DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.DateCreated, &p.CreatedBy); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPost retrieves a single post by ID
func (db *DB) GetPost(id int) (models.Post, error) {
	var p models.Post
	err := db.QueryRow(`
		SELECT id, title, content, date_created, created_by 
		FROM posts 
		WHERE id = $1
	`, id).Scan(&p.ID, &p.Title, &p.Content, &p.DateCreated, &p.CreatedBy)

	if err == sql.ErrNoRows {
		return models.Post{}, ErrNotFound
	}

	if err != nil {
		return models.Post{}, err
	}

	return p, nil
}

// CreatePost adds a new post to the database
func (db *DB) CreatePost(np models.NewPost) (models.Post, error) {
	var p models.Post
	err := db.QueryRow(`
		INSERT INTO posts (title, content, created_by) 
		VALUES ($1, $2, $3) 
		RETURNING id, title, content, date_created, created_by
	`, np.Title, np.Content, np.CreatedBy).Scan(
		&p.ID, &p.Title, &p.Content, &p.DateCreated, &p.CreatedBy,
	)

	if err != nil {
		return models.Post{}, err
	}

	return p, nil
}

// UpdatePost modifies an existing post
func (db *DB) UpdatePost(id int, up models.UpdatePost) (models.Post, error) {
	var p models.Post
	err := db.QueryRow(`
		UPDATE posts 
		SET title = $1, content = $2 
		WHERE id = $3 
		RETURNING id, title, content, date_created, created_by
	`, up.Title, up.Content, id).Scan(
		&p.ID, &p.Title, &p.Content, &p.DateCreated, &p.CreatedBy,
	)

	if err == sql.ErrNoRows {
		return models.Post{}, ErrNotFound
	}

	if err != nil {
		return models.Post{}, err
	}

	return p, nil
}

// DeletePost removes a post from the database
func (db *DB) DeletePost(id int) error {
	result, err := db.Exec(`DELETE FROM posts WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
