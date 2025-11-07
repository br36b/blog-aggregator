-- name: CreateFeed :one
INSERT INTO feeds (
    id,
    created_at,
    updated_at,
    name,
    url,
    user_id
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE feeds.url = $1;

-- name: GetFeedCreatorName :one
SELECT
    users.name
FROM
    feeds
JOIN
    users ON users.id = feeds.user_id
WHERE
    $1 = users.id;
