-- name: GetProfileByCardId :one
SELECT *
FROM gf10dm9_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: CreateCardProfile :one
INSERT INTO gf10dm9_card_profile (game_type, cardid, name, pass, type, update_flag, recovery)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRecoveryCount :exec
UPDATE gf10dm9_card_profile
SET recovery = ?
WHERE cardid = ?
AND game_type = ?;

-- name: GetRecentPuzzle :one
SELECT *
FROM gf10dm9_puzzle AS gp
INNER JOIN gf10dm9_card_profile AS gcp ON (gp.cardid = gcp.cardid)
WHERE gp.cardid = ?
AND gp.game_type = ?
AND gp.puzzle_no = gcp.puzzle_no;

-- name: GetPuzzleByNumber :one
SELECT *
FROM gf10dm9_puzzle
WHERE cardid = ?
AND game_type = ?
AND puzzle_no = ?;

-- name: UpsertPuzzle :exec
INSERT INTO gf10dm9_puzzle
(game_type, cardid, puzzle_no, flags, hidden)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, cardid, puzzle_no) DO
UPDATE SET flags=?4, hidden=?5;