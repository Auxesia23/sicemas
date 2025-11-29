-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    nip_index       BYTEA NOT NULL UNIQUE,
    nama_lengkap    VARCHAR(150) NOT NULL,
    peran           VARCHAR(150) NOT NULL,
    nip             VARCHAR(255) NOT NULL,
    email           VARCHAR(255) NOT NULL,
    nomor_telepon   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMPTZ  DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
