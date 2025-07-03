package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Connection parameters
	connStr := "postgres://bloguser:blogpassword@localhost:5432/blogdb?sslmode=disable"

	// Try to connect with retries
	var db *sql.DB
	var err error

	// Retry logic for connection
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
	}
	defer db.Close()

	// Check if posts table exists
	var tableExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'posts'
		)
	`).Scan(&tableExists)

	if err != nil {
		log.Fatalf("Error checking if table exists: %v", err)
	}

	if tableExists {
		fmt.Println("✅ Posts table already exists, no need to apply migration")
		return
	}

	// Read the migration file
	migrationSQL, err := ioutil.ReadFile("migrations/01_create_posts_table.sql")
	if err != nil {
		log.Fatalf("Error reading migration file: %v", err)
	}

	// Apply the migration
	fmt.Println("Applying migration to create posts table...")
	_, err = db.Exec(string(migrationSQL))
	if err != nil {
		log.Fatalf("Error applying migration: %v", err)
	}

	fmt.Println("✅ Migration applied successfully")

	// Verify the table was created
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'posts'
		)
	`).Scan(&tableExists)

	if err != nil {
		log.Fatalf("Error checking if table exists after migration: %v", err)
	}

	if tableExists {
		fmt.Println("✅ Verified posts table exists")
	} else {
		fmt.Println("❌ Failed to create posts table")
		return
	}
}
