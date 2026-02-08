-- name: CreateAggregateRoll :one
INSERT INTO aggregate_rolls (id, created_at, updated_at, string, result, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
)
RETURNING id, created_at, updated_at, string, result, username;

-- name: CreateRoll :one
INSERT INTO rolls (id, created_at, updated_at, string, result, individual_rolls, aggregate_roll_id, username)
VALUES (
    gen_random_uuid()
    ,NOW()
    ,NOW()
    ,$1
    ,$2
    ,$3
    ,$4
    ,$5
)
RETURNING id, created_at, updated_at, string, result, individual_rolls, aggregate_roll_id, username;