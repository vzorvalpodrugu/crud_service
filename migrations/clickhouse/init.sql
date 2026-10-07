CREATE DATABASE IF NOT EXISTS analytics;

CREATE TABLE IF NOT EXISTS analytics.posts_events (
    event_id UInt64,
    event_type String,

    post_id UInt64,
    post_name String,
    post_text String,
    post_author_id UInt64,
    post_created_at DateTime64(3, UTC),
    post_updated_at DateTime64(3, UTC),
    received_at DateTime64(3, UTC) default now64()
)
ENGINE = MergeTree()
ORDER BY (event_type, post_id, received_at);