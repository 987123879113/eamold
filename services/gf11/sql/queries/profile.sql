-- name: GetProfileByCardId :one
SELECT *
FROM gf11dm10_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: GetProfileByGdId :one
SELECT *
FROM gf11dm10_card_profile
WHERE gdid = ?
AND game_type = ?;

-- name: CreateCardProfile :one
INSERT INTO gf11dm10_card_profile (game_type, gdid, cardid, irid, name, pass, type, update_flag, puzzle_no, recovery)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRecoveryCount :exec
UPDATE gf11dm10_card_profile
SET recovery = ?
WHERE gdid = ?
AND game_type = ?;

-- name: GetRecentPuzzle :one
SELECT gp.*
FROM gf11dm10_puzzle AS gp
INNER JOIN gf11dm10_card_profile AS gcp ON (gp.game_type = gcp.game_type AND gp.gdid = gcp.gdid)
WHERE gp.gdid = ?
AND gp.game_type = ?
AND gp.puzzle_no = gcp.puzzle_no;

-- name: GetPuzzleByNumber :one
SELECT *
FROM gf11dm10_puzzle
WHERE gdid = ?
AND game_type = ?
AND puzzle_no = ?;

-- name: UpsertPuzzle :exec
INSERT INTO gf11dm10_puzzle
(game_type, gdid, puzzle_no, flags, hidden)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, gdid, puzzle_no) DO
UPDATE SET flags=?4, hidden=?5;