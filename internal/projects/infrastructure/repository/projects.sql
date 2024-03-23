-- name: CreateProject :one
INSERT INTO projects (public_id, name, description, tags, thumbnail_url, website_url, live, post_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;