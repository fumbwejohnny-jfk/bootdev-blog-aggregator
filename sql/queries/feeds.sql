-- name: CreateFeed :one
INSERT INTO feeds (id, url, created_at, updated_at, name, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: GetFeedById :one
SELECT name FROM feeds WHERE id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds; 

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: CreateFeedFollow :one
WITH new_follow AS (
    INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT
    new_follow.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM new_follow
JOIN users
    ON users.id = new_follow.user_id
JOIN feeds
    ON feeds.id = new_follow.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.id,
    feed_follows.user_id,
    feed_follows.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users
    ON users.id = feed_follows.user_id
INNER JOIN feeds
    ON feeds.id = feed_follows.feed_id
    WHERE feed_follows.user_id = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;


-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
-- WHERE last_fetched_at IS NULL OR last_fetched_at < NOW() - INTERVAL '1 hour'
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;
