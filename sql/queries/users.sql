-- name: CreateUser :one
INSERT INTO users (username, created_at, updated_at, hashed_password)
VALUES (
    $1
    ,NOW()
    ,NOW()
    ,$2
)
RETURNING username, created_at, updated_at;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUserWithUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: UpdateUser :one
UPDATE users 
SET (updated_at, username, hashed_password) = (NOW(), $1 ,$2)
WHERE username = $3
RETURNING username, created_at, updated_at;