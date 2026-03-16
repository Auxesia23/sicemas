-- +goose Up
-- +goose StatementBegin
CREATE TABLE foto_situs (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    situs_id UUID REFERENCES situs_keagamaan(id) ON DELETE CASCADE,
    image_url VARCHAR(255) NOT NULL,
    public_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS foto_situs;
-- +goose StatementEnd
