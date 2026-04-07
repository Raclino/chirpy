-- name: CreateRefreshToken :one
INSERT INTO
    refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES
    ($1, $2, $3, $4, $5, $6) 
RETURNING token;

-- name: GetTokenByRefreshToken :one
SELECT
    *
FROM
    refresh_tokens
WHERE
    token = $1;