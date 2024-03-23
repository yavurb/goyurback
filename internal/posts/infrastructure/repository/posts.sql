-- name: CreatePost :one
INSERT INTO posts (title, author, slug, description, content) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC;

-- name: GetPostBySlug :one
SELECT * FROM posts WHERE slug = $1;

-- name: UpdatePost :one
UPDATE posts SET title = $1, author = $2, slug = $3, description = $4, content = $5, status = $6, published_at = $7, updated_at = now() WHERE id = $8 RETURNING *;
