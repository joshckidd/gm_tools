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

-- name: GetCustomFieldForType :many
SELECT *
FROM custom_fields
WHERE type_id = $1;

-- name: CreateItem :one
INSERT INTO items (id, created_at, updated_at, item_name, item_description, type_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
    ,$4
)
RETURNING id, created_at, updated_at, item_name, item_description, type_id, username;

-- name: CreateCustomFieldValue :one
INSERT INTO custom_field_values (id, created_at, updated_at, custom_field_value, custom_field_id, item_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
    ,$4
)
RETURNING id, created_at, updated_at, custom_field_value, custom_field_id, item_id, username;

-- name: GetItems :many
SELECT *
FROM items
ORDER BY created_at;

-- name: GetCustomFieldValues :many
SELECT 
    custom_fields.custom_field_name
    ,custom_field_values.custom_field_value
FROM custom_field_values
JOIN custom_fields ON custom_fields.id = custom_field_values.custom_field_id
WHERE item_id = $1
ORDER BY custom_fields.created_at;

-- name: DeleteType :exec
DELETE FROM types
WHERE id = $1;