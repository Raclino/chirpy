-- name: CreateUser :one
INSERT INTO
    users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES
    ($1, $2, $3, $4, $5, $6) 
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: DeleteAllUsers :exec
DELETE FROM
    users RETURNING *;

-- name: GetUsers :many
SELECT
    *
FROM
    users;

-- name: UpdateUserPwdEmail :one
UPDATE users 
SET email = $2, hashed_password = $3, updated_at = $4
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpdateUserChirpyRedMembership :one
UPDATE users
SET is_chirpy_red = true
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;