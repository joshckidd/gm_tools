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

-- name: GetAggregateRolls :many
SELECT *
FROM aggregate_rolls
WHERE username = $1;

-- name: GetRolls :many
SELECT *
FROM rolls
WHERE aggregate_roll_id = $1;

-- name: DeleteOldRolls :exec
DELETE FROM aggregate_rolls
WHERE created_at < $1;
