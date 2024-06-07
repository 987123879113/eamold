-- name: GetMessages :many
SELECT message
FROM gf8dm7_messages
WHERE enabled = 1
AND game_type = ?;
