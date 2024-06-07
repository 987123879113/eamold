-- name: GetMessages :many
SELECT message
FROM gf9dm8_messages
WHERE enabled = 1
AND game_type = ?;
