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