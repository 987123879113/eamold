-- name: GetExChallengeProgress :many
SELECT *
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?;

-- name: GetExChallengeSerial :one
SELECT serial
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?
AND exid = ?;

-- name: GetExChallengeClearCounts :many
SELECT exid, COUNT(exid) AS clear_count
FROM gf10dm9_ex_progress
WHERE game_type = ?
AND clear = 1
GROUP BY exid;

-- name: GetLatestClearedExChallenge :one
SELECT exid
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?
AND clear = 1
ORDER BY exid DESC;

-- name: UpsertExChallenge :exec
INSERT INTO gf10dm9_ex_progress
(game_type, cardid, exid, seen, clear)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, cardid, exid) DO
UPDATE SET seen=?4, clear=?5;

-- name: AwardSerialForExChallenge :exec
UPDATE gf10dm9_ex_progress
SET serial = (SELECT COUNT(*)
    FROM gf10dm9_ex_progress AS t1
    WHERE t1.game_type = ?1
    AND t1.exid = ?3
    AND t1.clear = 1
    AND t1.serial != 0
    AND t1.cardid != ?2) + 1
WHERE gf10dm9_ex_progress.game_type = ?1
AND gf10dm9_ex_progress.cardid = ?2
AND gf10dm9_ex_progress.exid = ?3
AND gf10dm9_ex_progress.serial = 0;
