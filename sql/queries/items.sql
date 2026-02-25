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

-- name: GetTypeByName :one
SELECT *
FROM types
WHERE type_name = $1;


-- name: CreateCustomFields :one
INSERT INTO custom_fields (id, created_at, updated_at, custom_field_name, custom_field_type, type_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
    ,$4
)
RETURNING id, created_at, updated_at, custom_field_name, custom_field_type, type_id, username;

-- name: GetCustomFields :many
SELECT *
FROM custom_fields
ORDER BY created_at;