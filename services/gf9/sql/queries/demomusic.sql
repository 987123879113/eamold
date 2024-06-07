-- name: GetDemoMusic :many
SELECT musicid
FROM gf9dm8_demomusic
WHERE game_type = ?
LIMIT ?;
