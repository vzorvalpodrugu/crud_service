-- +goose up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS outbox(
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    event_type VARCHAR(100) NOT NULL,
    infodata JSONB NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    sent_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_outbox_status ON outbox(status);

-- +goose StatementEnd


-- +goose down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_outbox_status;
DROP TABLE IF EXISTS outbox;
-- +goose StatementEnd