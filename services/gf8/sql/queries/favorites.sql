-- name: GetFavorites :many
SELECT musicid
FROM gf8dm7_favorites
WHERE game_type = ?
ORDER BY count DESC
LIMIT ?;

-- name: UpdateFavoriteCount :exec
INSERT INTO gf8dm7_favorites (game_type, musicid, count)
VALUES (?, ?, 1)
ON CONFLICT(game_type, musicid) DO
UPDATE SET count = count + 1;
