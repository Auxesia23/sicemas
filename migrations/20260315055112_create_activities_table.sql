-- +goose Up
-- +goose StatementBegin
CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    target_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Bikin index di kolom created_at karena dashboard akan selalu nge-query data terbaru (ORDER BY created_at DESC LIMIT 5)
CREATE INDEX idx_activities_created_at ON activities(created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS activities;
-- +goose StatementEnd
