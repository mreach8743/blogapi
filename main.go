package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog2/auth"
	"blog2/db"
	"blog2/handlers"
)

const (
	dbConnStr  = "postgres://bloguser:blogpassword@localhost:5432/blogdb?sslmode=disable"
	serverAddr = ":8080"
)

func main() {
	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Blog API server...")

	// Connect to the database
	database, err := db.NewDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Connected to database successfully")

	// Set up JWT configuration
	jwtConfig := auth.DefaultJWTConfig()

	// Create handlers
	postsHandler := handlers.NewPostsHandler(database)
	usersHandler := handlers.NewUsersHandler(database, jwtConfig)

	// Set up routes
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.Handle("/users/register", usersHandler)
	mux.Handle("/users/login", usersHandler)

	// Protected routes (authentication required)
	protectedHandler := auth.RequireAuth(jwtConfig)(postsHandler)
	mux.Handle("/posts", protectedHandler)
	mux.Handle("/posts/", protectedHandler)

	// Protected user routes
	protectedUserHandler := auth.RequireAuth(jwtConfig)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/me" {
			usersHandler.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	mux.Handle("/users/me", protectedUserHandler)

	// Add middleware for logging
	handler := logMiddleware(mux)

	// Configure the server
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// logMiddleware logs all HTTP requests
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}
