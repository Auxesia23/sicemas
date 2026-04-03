-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ NULL;
ALTER TABLE situs_keagamaan ADD COLUMN deleted_at TIMESTAMPTZ NULL;
CREATE INDEX idx_users_active ON users(id) WHERE deleted_at IS NULL;
CREATE INDEX idx_situs_active ON situs_keagamaan(id) WHERE deleted_at IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_situs_active;
DROP INDEX IF EXISTS idx_users_active;
ALTER TABLE situs_keagamaan DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;
-- +goose StatementEnd
