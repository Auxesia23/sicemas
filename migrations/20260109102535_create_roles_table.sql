-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    name            VARCHAR(50) UNIQUE NOT NULL,
    created_at      TIMESTAMPTZ  DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
