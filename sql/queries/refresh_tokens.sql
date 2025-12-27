-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    token,
    created_at,
    updated_at,
    user_id,
    expires_at
) VALUES (
    $1,
    now(),
    now(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT
    refresh_tokens.token,
    refresh_tokens.expires_at,
    refresh_tokens.revoked_at,
    users.id AS user_id
FROM refresh_tokens
INNER JOIN users
    ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = now(), revoked_at = now()
WHERE token = $1
RETURNING *;
