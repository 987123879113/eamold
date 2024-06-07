-- name: AddScore :exec
INSERT INTO gf8dm7puv_scores
(game_type, cardid, musicid, musicnum, seq, flags, encore, extra, score, clear, skill, combo)
VALUES
(?,?,?,?,?,?,?,?,?,?,?,?);

-- name: GetTotalPlayedCount :one
SELECT COUNT(*) AS count
FROM gf8dm7puv_scores
WHERE musicid = ?
AND game_type = ?
AND seq = ?
GROUP BY musicid;

-- name: GetScoreRank :one
SELECT COUNT(*) AS count
FROM gf8dm7puv_scores
WHERE musicid = ?
AND game_type = ?
AND seq = ?
AND score >= ?;

-- name: IsBestScore :one
SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 0 ELSE 1 END AS INTEGER)
FROM gf8dm7puv_scores
WHERE cardid = ?
AND game_type = ?
AND musicid = ?
AND seq = ?
AND score >= ?;

-- name: GetSkillScoresByCardId :many
SELECT musicid
FROM gf8dm7puv_scores
WHERE cardid = ?
AND game_type = ?
ORDER BY score DESC
LIMIT ?;

-- name: GetSkillPointsByCardId :one
SELECT CAST(IFNULL(SUM(t.skill), 0) AS INTEGER) AS skill_points
FROM (
    SELECT skill
    FROM gf8dm7puv_scores
    WHERE cardid = ?
    AND game_type = ?
    ORDER BY score DESC
    LIMIT 30
) AS t;
