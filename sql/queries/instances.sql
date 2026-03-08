-- name: CreateInstance :one
INSERT INTO instances (id, created_at, updated_at, item_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
)
RETURNING id, created_at, updated_at, item_id, username;

-- name: CreateCustomFieldInstanceValue :one
INSERT INTO custom_field_instance_values (id, created_at, updated_at, custom_field_value, instance_id, custom_field_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
    ,$4
)
RETURNING id, created_at, updated_at, custom_field_value, instance_id, custom_field_id, username;

-- name: GetInstances :many
SELECT 
    instances.id
    ,instances.created_at
    ,instances.updated_at
    ,items.item_name
    ,items.item_description
    ,items.type_id
    ,instances.username
FROM instances
JOIN items ON items.id = instances.item_id
WHERE instances.username = $1
ORDER BY instances.created_at;

-- name: GetCustomFieldInstanceValues :many
SELECT 
    custom_field_instance_values.custom_field_value
    ,custom_fields.custom_field_name
FROM custom_field_instance_values
JOIN custom_fields ON custom_fields.id = custom_field_instance_values.custom_field_id
WHERE custom_field_instance_values.instance_id = $1;

-- name: DeleteOldInstances :exec
DELETE FROM instances
WHERE created_at < $1;