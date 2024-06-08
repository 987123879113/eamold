-- name: SaveScore :exec
INSERT INTO gf10dm9_scores (game_type, cardid, netid, courseid, seq_mode, flags, score, clear, combo, skill, perc, irall, ircom)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: RecalculateTotalSkillPointsForCardId :exec
UPDATE gf10dm9_card_profile
SET skill = (SELECT SUM(t.max_skill)
FROM (
    SELECT gs.netid, MAX(gs.skill) AS max_skill
    FROM gf10dm9_scores AS gs
    WHERE gs.cardid = ?1
    AND gs.game_type = ?2
    AND clear > 0
    GROUP BY gs.netid
    ORDER BY gs.skill DESC
    LIMIT 30
) as t)
WHERE gf10dm9_card_profile.cardid = ?1;

-- name: GetSongTotalPlayCount :one
SELECT COUNT(*)
FROM gf10dm9_scores
WHERE netid = ?
AND game_type = ?
AND seq_mode = ?;

-- name: GetSongCurrentRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_scores
WHERE netid = ?
AND game_type = ?
AND seq_mode = ?
AND skill > ?;

-- name: IsBestScore :one
SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 0 ELSE 1 END AS INTEGER)
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid = ?
AND seq_mode = ?
AND score >= ?;

-- name: GetPlayerCount :one
SELECT COUNT(*)
FROM gf10dm9_card_profile
WHERE game_type = ?;

-- name: GetPlayerSkill :one
SELECT skill
FROM gf10dm9_card_profile
WHERE cardid = ?
AND game_type = ?;

-- name: GetPlayerRank :one
WITH sorted_skills AS (
    SELECT cardid, skill
    FROM gf10dm9_card_profile
    WHERE game_type = ?
    ORDER BY skill DESC
), ranked_skills AS (
    SELECT cardid, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_skills
)
SELECT CAST(rank AS INTEGER)
FROM ranked_skills
WHERE cardid = ?;

-- name: GetSeqStatsByCardId :many
SELECT netid, seq_mode, CAST(MAX(skill) AS INTEGER) AS `skill`, perc
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid != -1
GROUP BY cardid, netid, seq_mode;

-- name: GetAllMaxSkillPointsByCardId :many
SELECT netid, seq_mode, CAST(MAX(skill) AS INTEGER) AS `skill`, perc
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid != -1
GROUP BY cardid, netid;

-- name: GetFavorites :many
SELECT netid
FROM gf10dm9_scores
WHERE netid != -1
AND game_type = ?
GROUP BY netid
ORDER BY COUNT(netid) DESC
LIMIT ?;

-- name: GetShopBySerial :one
SELECT pref, name, points
FROM gf10dm9_shops
WHERE sid = ?
AND game_type = ?;

-- name: GetRankedShops :many
SELECT pref, name, points
FROM gf10dm9_shops
WHERE game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetRankedShopsByPref :many
SELECT name, points
FROM gf10dm9_shops
WHERE pref = ?
AND game_type = ?
ORDER BY points DESC
LIMIT ?;

-- name: GetShopRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_shops
WHERE points > (
    SELECT t.points
    FROM gf10dm9_shops AS t
    WHERE t.sid = ?
    AND t.game_type = ?
    LIMIT 1
);

-- name: GetShopRankByPref :one
SELECT COUNT(*) + 1
FROM gf10dm9_shops
WHERE points > (
    SELECT t.points
    FROM gf10dm9_shops AS t
    WHERE t.sid = ?
    AND t.game_type = ?
    AND t.pref = ?
    LIMIT 1
);

-- name: UpdateShopPoints :exec
INSERT INTO gf10dm9_shops
(game_type, sid, pref, name, points)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, sid) DO
UPDATE SET points=points + ?5;

-- name: GetCourseTotalPlayCount :one
SELECT COUNT(*)
FROM gf10dm9_scores
WHERE courseid = ?
AND game_type = ?
AND seq_mode = ?;

-- name: GetCourseCurrentRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_scores
WHERE courseid = ?
AND game_type = ?
AND seq_mode = ?
AND score > ?;

-- name: GetCourseBestRank :one
WITH sorted_scores AS (
    SELECT cardid, score
    FROM gf10dm9_scores
    WHERE courseid = ?
    AND game_type = ?
    AND seq_mode = ?
    ORDER BY score DESC
), ranked_scores AS (
    SELECT cardid, ROW_NUMBER() OVER() AS `rank`
    FROM sorted_scores
)
SELECT CAST(MIN(rank) AS INTEGER)
FROM ranked_scores
WHERE cardid = ?;

-- name: UpdateLandPoints :exec
INSERT INTO gf10dm9_shops
(game_type, sid, pref, name, points)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, sid) DO
UPDATE SET points=points + ?5;