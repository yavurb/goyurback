-- name: CreateAPIKey :one
INSERT INTO apikeys (public_id, name, key) VALUES ($1, $2, $3) RETURNING *;

-- name: RevokeAPIKey :exec
UPDATE apikeys SET revoked_at = now(), updated_at = now(), revoked = true WHERE public_id = $1 RETURNING *;

-- name: GetAPIKeyByValue :one
SELECT * from apikeys WHERE key = $1;
