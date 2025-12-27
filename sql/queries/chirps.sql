-- name: CreateChirp :one
INSERT INTO chirps (
    id,
    body,
    created_at,
    updated_at,
    user_id
) VALUES (gen_random_uuid(), $1, now(), now(), $2)
RETURNING *;

-- name: ResetChirps :exec
DELETE FROM chirps;

-- name: GetChirps :many
SELECT *
FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT *
FROM chirps
WHERE id = $1;
