-- name: GetProfileByCardId :one
SELECT *
FROM gf8dm7puv_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: UpdateCardProfile :exec
INSERT INTO gf8dm7puv_card_profile (game_type, cardid, name, color, recovery, styles, hidden)
VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)
ON CONFLICT (cardid) DO
UPDATE SET color=?4, recovery=?5, styles=?6, hidden=?7;

-- name: UpdatePuzzleProgress :exec
INSERT INTO gf8dm7puv_puzzles (game_type, cardid, number, flags, out)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (cardid) DO
UPDATE SET number=?3, flags=?4, out=?5;
