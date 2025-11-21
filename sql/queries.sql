-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
-- name: GetUserByEmail :one
SELECT * FROM users WHERE LOWER(email) = $1 LIMIT 1;
-- name: CreateUser :exec
INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3);
