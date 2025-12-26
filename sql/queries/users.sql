-- name: CreateUser :one
INSERT INTO users (
    id,
    createdat,
    updatedat,
    email
) VALUES (get_random_uuid(), now(), now(), $1)
RETURNING *;
