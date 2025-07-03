# Blog API

A simple API for interacting with a blog, built with Go and PostgreSQL.

## Project Structure

- `Dockerfile` - Configures PostgreSQL database
- `docker-compose.yml` - Sets up the database service with migrations
- `migrations/` - Contains SQL migration files for database schema
- `test_migration.go` - Script to verify migrations worked correctly
- `models/` - Data structures used in the application
- `db/` - Database access layer
- `handlers/` - HTTP request handlers
- `main.go` - Application entry point and server configuration

## Setup Instructions

### 1. Start the Database

If this is your first time running the database, or if you want to reset it:

```bash
# Remove any existing containers and volumes
docker-compose down -v

# Start the database
docker-compose up -d
```

This will:
- Build and start the PostgreSQL container
- Apply the migrations in the `migrations` directory
- Create the `posts` table with the required fields

**Note:** PostgreSQL only runs initialization scripts when the database is first created. If you've already started the container before and the volume persists, the migrations won't run again. This is why we recommend using `docker-compose down -v` to remove volumes before starting if you're having issues.

#### Alternative: Apply Migration Manually

If you don't want to reset your database but need to create the post-table, you can use the provided script:

```bash
go run apply_migration.go
```

This script will check if the post-table exists and create it if it doesn't, without affecting any existing data.

### 2. Install Go Dependencies

```bash
go mod tidy
```

This will download the required Go dependencies, including the PostgreSQL driver.

### 3. Verify Migrations

```bash
go run test_migration.go
```

This script will connect to the database and verify that the `posts` table exists with the correct structure.

## Database Schema

### Posts Table

| Column       | Type                     | Description                   |
|--------------|--------------------------|-------------------------------|
| id           | SERIAL                   | Primary key                   |
| title        | VARCHAR(255)             | Title of the blog post        |
| content      | TEXT                     | Main content of the blog post |
| date_created | TIMESTAMP WITH TIME ZONE | When the post was created     |
| created_by   | VARCHAR(100)             | Author of the post            |

## Connection String

For Go applications:
```
postgres://bloguser:blogpassword@localhost:5432/blogdb?sslmode=disable
```

## API Endpoints

The API server runs on port 8080 by default. To start the server:

```bash
go run main.go
```

### Testing the API

A test script is provided to verify that all API endpoints are working correctly:

```bash
# In a separate terminal (after starting the server)
go run test_api.go
```

This script will:
1. Create a test post
2. Retrieve all posts
3. Retrieve the specific post
4. Update the post
5. Delete the post
6. Verify the post was deleted

If all tests pass, you'll see "All tests passed successfully!" in the console.

### Available Endpoints

#### GET /posts
Returns a list of all blog posts.

**Response:**
```json
[
  {
    "id": 1,
    "title": "First Post",
    "content": "This is my first blog post",
    "date_created": "2023-05-01T12:00:00Z",
    "created_by": "john"
  },
  {
    "id": 2,
    "title": "Second Post",
    "content": "This is another post",
    "date_created": "2023-05-02T14:30:00Z",
    "created_by": "jane"
  }
]
```

#### GET /posts/{id}
Returns a single blog post by ID.

**Response:**
```json
{
  "id": 1,
  "title": "First Post",
  "content": "This is my first blog post",
  "date_created": "2023-05-01T12:00:00Z",
  "created_by": "john"
}
```

#### POST /posts
Create a new blog post.

**Request:**
```json
{
  "title": "New Post",
  "content": "This is a new blog post",
  "created_by": "alice"
}
```

**Response:**
```json
{
  "id": 3,
  "title": "New Post",
  "content": "This is a new blog post",
  "date_created": "2023-05-03T10:15:00Z",
  "created_by": "alice"
}
```

#### PUT /posts/{id}
Updates an existing blog post.

**Request:**
```json
{
  "title": "Updated Post",
  "content": "This post has been updated"
}
```

**Response:**
```json
{
  "id": 1,
  "title": "Updated Post",
  "content": "This post has been updated",
  "date_created": "2023-05-01T12:00:00Z",
  "created_by": "john"
}
```

#### DELETE /posts/{id}
Delete a blog post.

**Response:** No content (204)

## Example Usage with cURL

### Get all posts
```bash
curl -X GET http://localhost:8080/posts
```

### Get a single post
```bash
curl -X GET http://localhost:8080/posts/1
```

### Create a new post
```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"New Post","content":"This is a new blog post","created_by":"alice"}'
```

### Update a post
```bash
curl -X PUT http://localhost:8080/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Post","content":"This post has been updated"}'
```

### Delete a post
```bash
curl -X DELETE http://localhost:8080/posts/1
```

## Stopping the Database

```bash
docker-compose down
```

To remove volumes as well:
```bash
docker-compose down -v
```
