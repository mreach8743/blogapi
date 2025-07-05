package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"blog2/models"
)

const baseURL = "http://localhost:8080"

func main() {
	// Wait for the server to start
	log.Println("Waiting for server to start...")
	time.Sleep(2 * time.Second)

	// Test user registration
	log.Println("Testing POST /users/register")
	token, user := registerUser()
	log.Printf("Registered user: %s with ID: %d\n", user.Username, user.ID)

	// Test user login
	log.Println("Testing POST /users/login")
	loginToken, loginUser := loginUser(user.Username)
	log.Printf("Logged in user: %s with ID: %d\n", loginUser.Username, loginUser.ID)

	// Test getting current user
	log.Println("Testing GET /users/me")
	currentUser := getCurrentUser(token)
	log.Printf("Current user: %s with ID: %d\n", currentUser.Username, currentUser.ID)

	// Test creating a post with authentication
	log.Println("Testing POST /posts with authentication")
	post := createPost(token)
	log.Printf("Created post with ID: %d\n", post.ID)

	// Test getting all posts with authentication
	log.Println("Testing GET /posts with authentication")
	posts := getAllPosts(token)
	log.Printf("Found %d posts\n", len(posts))

	// Test getting a single post with authentication
	log.Println("Testing GET /posts/{id} with authentication")
	singlePost := getPost(token, post.ID)
	log.Printf("Retrieved post: %s\n", singlePost.Title)

	// Test updating a post with authentication
	log.Println("Testing PUT /posts/{id} with authentication")
	updatedPost := updatePost(token, post.ID)
	log.Printf("Updated post title: %s\n", updatedPost.Title)

	// Test deleting a post with authentication
	log.Println("Testing DELETE /posts/{id} with authentication")
	deletePost(token, post.ID)
	log.Println("Post deleted successfully")

	// Verify the post was deleted
	log.Println("Verifying post was deleted")
	posts = getAllPosts(token)
	for _, p := range posts {
		if p.ID == post.ID {
			log.Fatalf("Post with ID %d still exists after deletion", post.ID)
		}
	}

	log.Println("All authentication tests passed successfully!")
}

func registerUser() (string, models.User) {
	// Generate a unique username to avoid conflicts
	username := fmt.Sprintf("testuser_%d", time.Now().Unix())

	newUser := models.NewUser{
		Username: username,
		Email:    username + "@example.com",
		Password: "password123",
	}

	jsonData, err := json.Marshal(newUser)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(baseURL+"/users/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to register user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to register user. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var loginResponse models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return loginResponse.Token, loginResponse.User
}

func loginUser(username string) (string, models.User) {
	loginRequest := models.LoginRequest{
		Username: username,
		Password: "password123",
	}

	jsonData, err := json.Marshal(loginRequest)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to login. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var loginResponse models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return loginResponse.Token, loginResponse.User
}

func getCurrentUser(token string) models.User {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/users/me", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to get current user. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return user
}

func createPost(token string) models.Post {
	newPost := models.NewPost{
		Title:     "Test Post",
		Content:   "This is a test post created by the API test script",
		CreatedBy: "test_script",
	}

	jsonData, err := json.Marshal(newPost)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/posts", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to create post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to create post. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var post models.Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return post
}

func getAllPosts(token string) []models.Post {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/posts", nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get posts: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to get posts. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var posts []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return posts
}

func getPost(token string, id int) models.Post {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/posts/%d", baseURL, id), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to get post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to get post. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var post models.Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return post
}

func updatePost(token string, id int) models.Post {
	updatePost := models.UpdatePost{
		Title:   "Updated Test Post",
		Content: "This post has been updated by the test script",
	}

	jsonData, err := json.Marshal(updatePost)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/posts/%d", baseURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to update post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to update post. Status: %d, Response: %s", resp.StatusCode, string(body))
	}

	var post models.Post
	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return post
}

func deletePost(token string, id int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/posts/%d", baseURL, id), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to delete post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatalf("Failed to delete post. Status: %d, Response: %s", resp.StatusCode, string(body))
	}
}
