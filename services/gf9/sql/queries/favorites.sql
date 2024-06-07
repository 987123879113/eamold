-- name: GetFavorites :many
SELECT musicid
FROM gf9dm8_favorites
WHERE game_Type = ?
ORDER BY count DESC
LIMIT ?;

-- name: UpdateFavoriteCount :exec
INSERT INTO gf9dm8_favorites (game_type, musicid, count)
VALUES (?, ?, 1)
ON CONFLICT(game_type, musicid) DO
UPDATE SET count = count + 1;
