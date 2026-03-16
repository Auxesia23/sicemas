-- +goose Up
-- +goose StatementBegin
CREATE TABLE jenis_situs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    nama_jenis VARCHAR(50) NOT NULL UNIQUE,
    deskripsi TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jenis_situs;
-- +goose StatementEnd
