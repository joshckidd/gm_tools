-- +goose Up
CREATE TABLE instances (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,item_id UUID NOT NULL REFERENCES items ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE custom_field_instance_values (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,custom_field_value TEXT NOT NULL
    ,instance_id UUID NOT NULL REFERENCES instances ON DELETE CASCADE
    ,custom_field_id UUID NOT NULL REFERENCES custom_fields ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE custom_field_instance_values;
DROP TABLE instances;