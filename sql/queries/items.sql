-- name: CreateType :one
INSERT INTO types (id, created_at, updated_at, type_name, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
)
RETURNING id, created_at, updated_at, type_name, username;

-- name: GetTypes :many
SELECT *
FROM types
ORDER BY created_at;