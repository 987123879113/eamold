-- name: GetDemoMusic :many
SELECT musicid
FROM gf8dm7_demomusic
WHERE game_type = ?
LIMIT ?;
