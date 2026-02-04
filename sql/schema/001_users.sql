-- +goose Up
CREATE TABLE users (
    username TEXT NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,hashed_password TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;