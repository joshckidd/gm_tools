-- +goose Up
CREATE TABLE types (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,type_name TEXT UNIQUE NOT NULL
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE items (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,item_name TEXT NOT NULL
    ,item_description TEXT NOT NULL
    ,type_id UUID NOT NULL REFERENCES types ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE custom_fields (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,custom_field_name TEXT NOT NULL
    ,custom_field_type TEXT NOT NULL
    ,type_id UUID NOT NULL REFERENCES types ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

CREATE TABLE custom_field_values (
    id UUID NOT NULL PRIMARY KEY
    ,created_at TIMESTAMP NOT NULL
    ,updated_at TIMESTAMP NOT NULL
    ,custom_field_value TEXT NOT NULL
    ,item_id UUID NOT NULL REFERENCES items ON DELETE CASCADE
    ,custom_field_id UUID NOT NULL REFERENCES custom_fields ON DELETE CASCADE
    ,username TEXT NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE custom_field_values;
DROP TABLE custom_fields;
DROP TABLE items;
DROP TABLE types;