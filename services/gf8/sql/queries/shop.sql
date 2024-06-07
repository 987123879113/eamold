-- name: GetRankedShops :many
SELECT name, points, pref
FROM gf8dm7_shops
WHERE game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetRankedShopsByPref :many
SELECT name, points
FROM gf8dm7_shops
WHERE pref = ?
AND game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: UpdateShopPoints :exec
INSERT INTO gf8dm7_shops (game_type, pref, name, points)
VALUES (?1, ?2, ?3, ?4)
ON CONFLICT(game_type, pref, name) DO
UPDATE SET points = points + ?4;
