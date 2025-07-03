-- Create posts table
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    date_created TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(100) NOT NULL
);

-- Add indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_posts_date_created ON posts(date_created);
CREATE INDEX IF NOT EXISTS idx_posts_created_by ON posts(created_by);

-- Add comments to document the table
COMMENT ON TABLE posts IS 'Stores blog post content and metadata';
COMMENT ON COLUMN posts.id IS 'Unique identifier for each post';
COMMENT ON COLUMN posts.title IS 'Title of the blog post';
COMMENT ON COLUMN posts.content IS 'Main content of the blog post';
COMMENT ON COLUMN posts.date_created IS 'Timestamp when the post was created';
COMMENT ON COLUMN posts.created_by IS 'Username or identifier of the post author';