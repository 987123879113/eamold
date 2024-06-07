-- name: GetRankedShops :many
SELECT name, points, pref
FROM gf9dm8_shops
WHERE game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetRankedShopsByPref :many
SELECT name, points
FROM gf9dm8_shops
WHERE pref = ?
AND game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetShopRank :one
WITH sorted_shops AS (
    SELECT *
    FROM gf9dm8_shops
    WHERE game_type = ?
    ORDER BY points DESC
), ranked_shops AS (
    SELECT *, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_shops
)
SELECT CAST(rank AS INTEGER)
FROM ranked_shops
WHERE name = ?
AND pref = ?;

-- name: GetShopRankByPref :one
WITH sorted_shops AS (
    SELECT *
    FROM gf9dm8_shops
    WHERE pref = ?
    AND game_type = ?
    ORDER BY points DESC
), ranked_shops AS (
    SELECT *, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_shops
)
SELECT CAST(rank AS INTEGER)
FROM ranked_shops
WHERE name = ?;

-- name: UpdateShopPoints :exec
INSERT INTO gf9dm8_shops (game_type, pref, name, points)
VALUES (?1, ?2, ?3, ?4)
ON CONFLICT(game_type, pref, name) DO
UPDATE SET points = points + ?4;

-- name: GetShopByPcbid :one
SELECT s.*
FROM gf9dm8_shops AS s
INNER JOIN gf9dm8_shop_machines AS sm ON sm.pcbid = ?
WHERE s.name = sm.name
AND s.game_type = ?
AND s.pref = sm.pref;

-- name: AddMachineToShop :exec
INSERT INTO gf9dm8_shop_machines
(game_type, pcbid, name, pref)
VALUES
(?1, ?2, ?3, ?4)
ON CONFLICT(pcbid) DO
UPDATE SET name = ?3, pref = ?4;