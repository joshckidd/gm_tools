-- +goose Up
CREATE TABLE aggregate_rolls (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,string TEXT NOT NULL
    ,result INT NOT NULL
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE rolls (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,string TEXT NOT NULL
    ,result INT NOT NULL
    ,individual_rolls TEXT NOT NULL
    ,aggregate_roll_id UUID NOT NULL REFERENCES aggregate_rolls ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE rolls;
DROP TABLE aggregate_rolls;