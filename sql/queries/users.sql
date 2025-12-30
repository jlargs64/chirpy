-- name: CreateUser :one
INSERT INTO users (
    id,
    created_at,
    updated_at,
    email,
    hashed_password,
    is_chirpy_red
) VALUES (gen_random_uuid(), now(), now(), $1, $2, false)
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: UpdateUserById :one
UPDATE users
SET email = $1, hashed_password = $2, updated_at = now()
WHERE id = $3
RETURNING *;

-- name: UpgradeUserToChirpyRed :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING *;
