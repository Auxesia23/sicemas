-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT uuidv7(),
    nip_index       BYTEA NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    nama_lengkap    VARCHAR(255) NOT NULL,
    nip             TEXT NOT NULL,
    jabatan         VARCHAR(255) NOT NULL,
    unit_kerja      VARCHAR(255) NOT NULL,
    email           TEXT,
    nomor_telepon   TEXT NOT NULL,
    created_at      TIMESTAMPTZ  DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
