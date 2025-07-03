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

	// Test creating a post
	log.Println("Testing POST /posts")
	post := createPost()
	log.Printf("Created post with ID: %d\n", post.ID)

	// Test getting all posts
	log.Println("Testing GET /posts")
	posts := getAllPosts()
	log.Printf("Found %d posts\n", len(posts))

	// Test getting a single post
	log.Println("Testing GET /posts/{id}")
	singlePost := getPost(post.ID)
	log.Printf("Retrieved post: %s\n", singlePost.Title)

	// Test updating a post
	log.Println("Testing PUT /posts/{id}")
	updatedPost := updatePost(post.ID)
	log.Printf("Updated post title: %s\n", updatedPost.Title)

	// Test deleting a post
	log.Println("Testing DELETE /posts/{id}")
	deletePost(post.ID)
	log.Println("Post deleted successfully")

	// Verify the post was deleted
	log.Println("Verifying post was deleted")
	posts = getAllPosts()
	for _, p := range posts {
		if p.ID == post.ID {
			log.Fatalf("Post with ID %d still exists after deletion", post.ID)
		}
	}

	log.Println("All tests passed successfully!")
}

func createPost() models.Post {
	newPost := models.NewPost{
		Title:     "Test Post",
		Content:   "This is a test post created by the API test script",
		CreatedBy: "test_script",
	}

	jsonData, err := json.Marshal(newPost)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(baseURL+"/posts", "application/json", bytes.NewBuffer(jsonData))
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

func getAllPosts() []models.Post {
	resp, err := http.Get(baseURL + "/posts")
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

func getPost(id int) models.Post {
	resp, err := http.Get(fmt.Sprintf("%s/posts/%d", baseURL, id))
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

func updatePost(id int) models.Post {
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

func deletePost(id int) {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/posts/%d", baseURL, id), nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

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
