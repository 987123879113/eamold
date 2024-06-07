-- name: GetShop :one
SELECT *
FROM gf8dm7puv_shops
WHERE name = ?
AND game_type = ?
AND pref = ?;

-- name: GetRankedShops :many
SELECT name, points, pref
FROM gf8dm7puv_shops
WHERE game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetRankedShopsByPref :many
SELECT name, points
FROM gf8dm7puv_shops
WHERE pref = ?
AND game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetShopRank :one
WITH sorted_shops AS (
    SELECT *
    FROM gf8dm7puv_shops
    ORDER BY points DESC
), ranked_shops AS (
    SELECT *, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_shops
)
SELECT CAST(rank AS INTEGER)
FROM ranked_shops
WHERE name = ?
AND game_type = ?
AND pref = ?;

-- name: GetShopRankByPref :one
WITH sorted_shops AS (
    SELECT *
    FROM gf8dm7puv_shops
    WHERE pref = ?
    ORDER BY points DESC
), ranked_shops AS (
    SELECT *, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_shops
)
SELECT CAST(rank AS INTEGER)
FROM ranked_shops
WHERE name = ?
AND game_type = ?;

-- name: UpdateShopPoints :exec
INSERT INTO gf8dm7puv_shops (game_type, pref, name, points)
VALUES (?1, ?2, ?3, ?4)
ON CONFLICT(game_type, pref, name) DO
UPDATE SET points = points + ?4;

-- name: GetShopByPcbid :one
SELECT s.*
FROM gf8dm7puv_shops AS s
INNER JOIN gf8dm7puv_shop_machines AS sm ON sm.pcbid = ?
WHERE s.name = sm.name
AND s.game_type = ?
AND s.pref = sm.pref
LIMIT 1;

-- name: AddMachineToShop :exec
INSERT INTO gf8dm7puv_shop_machines
(game_type, pcbid, name, pref)
VALUES
(?1, ?2, ?3, ?4)
ON CONFLICT(pcbid) DO
UPDATE SET name = ?3, pref = ?4;