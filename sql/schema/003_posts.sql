-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    description TEXT,
    published_at TIMESTAMP NOT NULL default NOW(),
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE IF EXISTS posts CASCADE;