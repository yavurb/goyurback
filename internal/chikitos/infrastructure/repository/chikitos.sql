-- name: CreateChikito :one
INSERT INTO chikitos (public_id, url, description) VALUES ($1, $2, $3) RETURNING *;