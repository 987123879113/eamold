-- name: AddScore :exec
INSERT INTO gf9dm8_scores
(game_type, cardid, music_num, seq_mode, flags, encore, extra, score, clear, skill, combo)
VALUES
(?,?,?,?,?,?,?,?,?,?,?);

-- name: GetSkillScoresByCardId :many
SELECT music_num
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
ORDER BY score DESC
LIMIT ?;

-- name: GetPlayerCount :one
SELECT COUNT(*)
FROM gf9dm8_card_profile
WHERE game_type = ?;

-- name: GetPlayerSkill :one
SELECT skill
FROM gf9dm8_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: GetPlayerRank :one
WITH sorted_skills AS (
    SELECT cardid, skill
    FROM gf9dm8_card_profile
    WHERE game_type = ?
    ORDER BY skill DESC
), ranked_skills AS (
    SELECT cardid, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_skills
)
SELECT CAST(rank AS INTEGER)
FROM ranked_skills
WHERE cardid = ?;

-- name: RecalculateTotalSkillPointsForCardId :exec
UPDATE gf9dm8_card_profile
SET skill = (SELECT SUM(t.max_skill)
FROM (
    SELECT gs.music_num, MAX(gs.skill) AS max_skill
    FROM gf9dm8_scores AS gs
    WHERE gs.cardid = ?1
    AND gs.game_type = ?2
    AND clear > 0
    GROUP BY gs.music_num
    ORDER BY gs.skill DESC
    LIMIT 30
) as t)
WHERE gf9dm8_card_profile.cardid = ?1;

-- name: GetSongTotalPlayCount :one
SELECT COUNT(*)
FROM gf9dm8_scores
WHERE music_num = ?
AND game_type = ?
AND seq_mode = ?;

-- name: GetSongCurrentRank :one
SELECT COUNT(*) + 1
FROM gf9dm8_scores
WHERE music_num = ?
AND game_type = ?
AND seq_mode = ?
AND skill > ?;

-- name: IsBestScore :one
SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 0 ELSE 1 END AS INTEGER)
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
AND music_num = ?
AND seq_mode = ?
AND score >= ?;

-- name: GetAllMaxSkillPointsByCardId :many
SELECT music_num, seq_mode, CAST(MAX(skill) AS INTEGER) AS `skill`
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
AND music_num != -1
GROUP BY cardid, music_num;
