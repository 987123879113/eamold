-- name: GetProfileByCardId :one
SELECT *
FROM gf9dm8_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: UpdateCardProfile :exec
INSERT INTO gf9dm8_card_profile (game_type, cardid, name, recovery, styles, hidden)
VALUES (?1, ?2, ?3, ?4, ?5, ?6)
ON CONFLICT (cardid) DO
UPDATE SET recovery=?4, styles=?5, hidden=?6;

-- name: UpdatePuzzleProgress :exec
INSERT INTO gf9dm8_puzzles (game_type, cardid, number, flags, out)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (cardid) DO
UPDATE SET number=?3, flags=?4, out=?5;
