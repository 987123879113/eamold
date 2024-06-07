-- name: AddScore :exec
INSERT INTO gf8dm7_scores
(game_type, musicid, seq, score)
VALUES
(?,?,?,?);

-- name: AddRankedScore :exec
INSERT INTO gf8dm7_ranked_scores
(game_type, musicid, seq, score, flags, name)
VALUES
(?,?,?,?,?,?);

-- name: GetTotalRankedScoreCounts :many
SELECT musicid, COUNT(*) AS count
FROM gf8dm7_scores
WHERE musicid IN (sqlc.slice('musicids'))
AND game_type = ?
GROUP BY musicid;


-- name: GetScoreRank :one
SELECT COUNT(*) AS count
FROM gf8dm7_scores
WHERE musicid = ?
AND game_type = ?
AND seq = ?
AND score >= ?;