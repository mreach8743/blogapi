package main

import (
	"database/sql"
	"fmt"
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
		fmt.Println("✅ Posts table exists")
	} else {
		fmt.Println("❌ Posts table does not exist")
		return
	}

	// Check table structure
	rows, err := db.Query(`
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_schema = 'public' 
		AND table_name = 'posts'
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Fatalf("Error querying table structure: %v", err)
	}
	defer rows.Close()

	fmt.Println("\nTable structure:")
	fmt.Println("----------------")

	for rows.Next() {
		var columnName, dataType string
		if err := rows.Scan(&columnName, &dataType); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("- %s (%s)\n", columnName, dataType)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}

	fmt.Println("\n✅ Migration test completed successfully")
}
