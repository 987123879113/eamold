-- name: GetDemoMusic :many
SELECT musicid
FROM gf8dm7puv_demomusic
WHERE game_type = ?
LIMIT ?;
