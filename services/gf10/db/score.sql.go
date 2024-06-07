// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: score.sql

package db

import (
	"context"
)

const getAllMaxSkillPointsByCardId = `-- name: GetAllMaxSkillPointsByCardId :many
SELECT netid, seq_mode, CAST(MAX(skill) AS INTEGER) AS ` + "`" + `skill` + "`" + `, perc
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid != -1
GROUP BY cardid, netid
`

type GetAllMaxSkillPointsByCardIdParams struct {
	Cardid   string
	GameType int64
}

type GetAllMaxSkillPointsByCardIdRow struct {
	Netid   int64
	SeqMode int64
	Skill   int64
	Perc    int64
}

func (q *Queries) GetAllMaxSkillPointsByCardId(ctx context.Context, arg GetAllMaxSkillPointsByCardIdParams) ([]GetAllMaxSkillPointsByCardIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMaxSkillPointsByCardId, arg.Cardid, arg.GameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMaxSkillPointsByCardIdRow
	for rows.Next() {
		var i GetAllMaxSkillPointsByCardIdRow
		if err := rows.Scan(
			&i.Netid,
			&i.SeqMode,
			&i.Skill,
			&i.Perc,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCourseBestRank = `-- name: GetCourseBestRank :one
WITH sorted_scores AS (
    SELECT cardid, score
    FROM gf10dm9_scores
    WHERE courseid = ?
    AND game_type = ?
    AND seq_mode = ?
    ORDER BY score DESC
), ranked_scores AS (
    SELECT cardid, ROW_NUMBER() OVER() AS ` + "`" + `rank` + "`" + `
    FROM sorted_scores
)
SELECT CAST(MIN(rank) AS INTEGER)
FROM ranked_scores
WHERE cardid = ?
`

type GetCourseBestRankParams struct {
	Courseid int64
	GameType int64
	SeqMode  int64
	Cardid   string
}

func (q *Queries) GetCourseBestRank(ctx context.Context, arg GetCourseBestRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCourseBestRank,
		arg.Courseid,
		arg.GameType,
		arg.SeqMode,
		arg.Cardid,
	)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getCourseCurrentRank = `-- name: GetCourseCurrentRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_scores
WHERE courseid = ?
AND game_type = ?
AND seq_mode = ?
AND score > ?
`

type GetCourseCurrentRankParams struct {
	Courseid int64
	GameType int64
	SeqMode  int64
	Score    int64
}

func (q *Queries) GetCourseCurrentRank(ctx context.Context, arg GetCourseCurrentRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCourseCurrentRank,
		arg.Courseid,
		arg.GameType,
		arg.SeqMode,
		arg.Score,
	)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getCourseTotalPlayCount = `-- name: GetCourseTotalPlayCount :one
SELECT COUNT(*)
FROM gf10dm9_scores
WHERE courseid = ?
AND game_type = ?
AND seq_mode = ?
`

type GetCourseTotalPlayCountParams struct {
	Courseid int64
	GameType int64
	SeqMode  int64
}

func (q *Queries) GetCourseTotalPlayCount(ctx context.Context, arg GetCourseTotalPlayCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getCourseTotalPlayCount, arg.Courseid, arg.GameType, arg.SeqMode)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getFavorites = `-- name: GetFavorites :many
SELECT netid
FROM gf10dm9_scores
WHERE netid != -1
AND game_type = ?
GROUP BY netid
ORDER BY COUNT(netid) DESC
LIMIT ?
`

type GetFavoritesParams struct {
	GameType int64
	Limit    int64
}

func (q *Queries) GetFavorites(ctx context.Context, arg GetFavoritesParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getFavorites, arg.GameType, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var netid int64
		if err := rows.Scan(&netid); err != nil {
			return nil, err
		}
		items = append(items, netid)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPlayerCount = `-- name: GetPlayerCount :one
SELECT COUNT(*)
FROM gf10dm9_card_profile
WHERE game_type = ?
`

func (q *Queries) GetPlayerCount(ctx context.Context, gameType int64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPlayerCount, gameType)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getPlayerRank = `-- name: GetPlayerRank :one
WITH sorted_skills AS (
    SELECT cardid, skill
    FROM gf10dm9_card_profile
    WHERE game_type = ?
    ORDER BY skill DESC
), ranked_skills AS (
    SELECT cardid, ROW_NUMBER() OVER() AS ` + "`" + `rank` + "`" + `
    FROM sorted_skills
)
SELECT CAST(rank AS INTEGER)
FROM ranked_skills
WHERE cardid = ?
`

type GetPlayerRankParams struct {
	GameType int64
	Cardid   string
}

func (q *Queries) GetPlayerRank(ctx context.Context, arg GetPlayerRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPlayerRank, arg.GameType, arg.Cardid)
	var rank int64
	err := row.Scan(&rank)
	return rank, err
}

const getPlayerSkill = `-- name: GetPlayerSkill :one
SELECT skill
FROM gf10dm9_card_profile
WHERE cardid = ?
AND game_type = ?
`

type GetPlayerSkillParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) GetPlayerSkill(ctx context.Context, arg GetPlayerSkillParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getPlayerSkill, arg.Cardid, arg.GameType)
	var skill int64
	err := row.Scan(&skill)
	return skill, err
}

const getRankedShops = `-- name: GetRankedShops :many
SELECT pref, name, points
FROM gf10dm9_shops
WHERE game_type = ?
ORDER BY points DESC
LIMIT ?
`

type GetRankedShopsParams struct {
	GameType int64
	Limit    int64
}

type GetRankedShopsRow struct {
	Pref   int64
	Name   string
	Points int64
}

func (q *Queries) GetRankedShops(ctx context.Context, arg GetRankedShopsParams) ([]GetRankedShopsRow, error) {
	rows, err := q.db.QueryContext(ctx, getRankedShops, arg.GameType, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRankedShopsRow
	for rows.Next() {
		var i GetRankedShopsRow
		if err := rows.Scan(&i.Pref, &i.Name, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRankedShopsByPref = `-- name: GetRankedShopsByPref :many
SELECT name, points
FROM gf10dm9_shops
WHERE pref = ?
AND game_type = ?
ORDER BY points DESC
LIMIT ?
`

type GetRankedShopsByPrefParams struct {
	Pref     int64
	GameType int64
	Limit    int64
}

type GetRankedShopsByPrefRow struct {
	Name   string
	Points int64
}

func (q *Queries) GetRankedShopsByPref(ctx context.Context, arg GetRankedShopsByPrefParams) ([]GetRankedShopsByPrefRow, error) {
	rows, err := q.db.QueryContext(ctx, getRankedShopsByPref, arg.Pref, arg.GameType, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRankedShopsByPrefRow
	for rows.Next() {
		var i GetRankedShopsByPrefRow
		if err := rows.Scan(&i.Name, &i.Points); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSeqStatsByCardId = `-- name: GetSeqStatsByCardId :many
SELECT netid, seq_mode, CAST(MAX(skill) AS INTEGER) AS ` + "`" + `skill` + "`" + `, perc
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid != -1
GROUP BY cardid, netid, seq_mode
`

type GetSeqStatsByCardIdParams struct {
	Cardid   string
	GameType int64
}

type GetSeqStatsByCardIdRow struct {
	Netid   int64
	SeqMode int64
	Skill   int64
	Perc    int64
}

func (q *Queries) GetSeqStatsByCardId(ctx context.Context, arg GetSeqStatsByCardIdParams) ([]GetSeqStatsByCardIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getSeqStatsByCardId, arg.Cardid, arg.GameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetSeqStatsByCardIdRow
	for rows.Next() {
		var i GetSeqStatsByCardIdRow
		if err := rows.Scan(
			&i.Netid,
			&i.SeqMode,
			&i.Skill,
			&i.Perc,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getShopBySerial = `-- name: GetShopBySerial :one
SELECT pref, name, points
FROM gf10dm9_shops
WHERE sid = ?
AND game_type = ?
`

type GetShopBySerialParams struct {
	Sid      string
	GameType int64
}

type GetShopBySerialRow struct {
	Pref   int64
	Name   string
	Points int64
}

func (q *Queries) GetShopBySerial(ctx context.Context, arg GetShopBySerialParams) (GetShopBySerialRow, error) {
	row := q.db.QueryRowContext(ctx, getShopBySerial, arg.Sid, arg.GameType)
	var i GetShopBySerialRow
	err := row.Scan(&i.Pref, &i.Name, &i.Points)
	return i, err
}

const getShopRank = `-- name: GetShopRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_shops
WHERE points > (
    SELECT t.points
    FROM gf10dm9_shops AS t
    WHERE t.sid = ?
    AND t.game_type = ?
    LIMIT 1
)
`

type GetShopRankParams struct {
	Sid      string
	GameType int64
}

func (q *Queries) GetShopRank(ctx context.Context, arg GetShopRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getShopRank, arg.Sid, arg.GameType)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getShopRankByPref = `-- name: GetShopRankByPref :one
SELECT COUNT(*) + 1
FROM gf10dm9_shops
WHERE points > (
    SELECT t.points
    FROM gf10dm9_shops AS t
    WHERE t.sid = ?
    AND t.game_type = ?
    AND t.pref = ?
    LIMIT 1
)
`

type GetShopRankByPrefParams struct {
	Sid      string
	GameType int64
	Pref     int64
}

func (q *Queries) GetShopRankByPref(ctx context.Context, arg GetShopRankByPrefParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getShopRankByPref, arg.Sid, arg.GameType, arg.Pref)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getSongCurrentRank = `-- name: GetSongCurrentRank :one
SELECT COUNT(*) + 1
FROM gf10dm9_scores
WHERE netid = ?
AND game_type = ?
AND seq_mode = ?
AND skill > ?
`

type GetSongCurrentRankParams struct {
	Netid    int64
	GameType int64
	SeqMode  int64
	Skill    int64
}

func (q *Queries) GetSongCurrentRank(ctx context.Context, arg GetSongCurrentRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSongCurrentRank,
		arg.Netid,
		arg.GameType,
		arg.SeqMode,
		arg.Skill,
	)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const getSongTotalPlayCount = `-- name: GetSongTotalPlayCount :one
SELECT COUNT(*)
FROM gf10dm9_scores
WHERE netid = ?
AND game_type = ?
AND seq_mode = ?
`

type GetSongTotalPlayCountParams struct {
	Netid    int64
	GameType int64
	SeqMode  int64
}

func (q *Queries) GetSongTotalPlayCount(ctx context.Context, arg GetSongTotalPlayCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSongTotalPlayCount, arg.Netid, arg.GameType, arg.SeqMode)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const isBestScore = `-- name: IsBestScore :one
SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 0 ELSE 1 END AS INTEGER)
FROM gf10dm9_scores
WHERE cardid = ?
AND game_type = ?
AND netid = ?
AND seq_mode = ?
AND score >= ?
`

type IsBestScoreParams struct {
	Cardid   string
	GameType int64
	Netid    int64
	SeqMode  int64
	Score    int64
}

func (q *Queries) IsBestScore(ctx context.Context, arg IsBestScoreParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, isBestScore,
		arg.Cardid,
		arg.GameType,
		arg.Netid,
		arg.SeqMode,
		arg.Score,
	)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const recalculateTotalSkillPointsForCardId = `-- name: RecalculateTotalSkillPointsForCardId :exec
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
WHERE gf10dm9_card_profile.cardid = ?1
`

type RecalculateTotalSkillPointsForCardIdParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) RecalculateTotalSkillPointsForCardId(ctx context.Context, arg RecalculateTotalSkillPointsForCardIdParams) error {
	_, err := q.db.ExecContext(ctx, recalculateTotalSkillPointsForCardId, arg.Cardid, arg.GameType)
	return err
}

const saveScore = `-- name: SaveScore :exec
INSERT INTO gf10dm9_scores (game_type, cardid, netid, courseid, seq_mode, flags, score, clear, combo, skill, perc, irall, ircom)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

type SaveScoreParams struct {
	GameType int64
	Cardid   string
	Netid    int64
	Courseid int64
	SeqMode  int64
	Flags    int64
	Score    int64
	Clear    int64
	Combo    int64
	Skill    int64
	Perc     int64
	Irall    int64
	Ircom    int64
}

func (q *Queries) SaveScore(ctx context.Context, arg SaveScoreParams) error {
	_, err := q.db.ExecContext(ctx, saveScore,
		arg.GameType,
		arg.Cardid,
		arg.Netid,
		arg.Courseid,
		arg.SeqMode,
		arg.Flags,
		arg.Score,
		arg.Clear,
		arg.Combo,
		arg.Skill,
		arg.Perc,
		arg.Irall,
		arg.Ircom,
	)
	return err
}

const updateLandPoints = `-- name: UpdateLandPoints :exec
INSERT INTO gf10dm9_shops
(game_type, sid, pref, name, points)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (sid) DO
UPDATE SET points=points + ?5
`

type UpdateLandPointsParams struct {
	GameType int64
	Sid      string
	Pref     int64
	Name     string
	Points   int64
}

func (q *Queries) UpdateLandPoints(ctx context.Context, arg UpdateLandPointsParams) error {
	_, err := q.db.ExecContext(ctx, updateLandPoints,
		arg.GameType,
		arg.Sid,
		arg.Pref,
		arg.Name,
		arg.Points,
	)
	return err
}

const updateShopPoints = `-- name: UpdateShopPoints :exec
INSERT INTO gf10dm9_shops
(game_type, sid, pref, name, points)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (sid) DO
UPDATE SET points=points + ?5
`

type UpdateShopPointsParams struct {
	GameType int64
	Sid      string
	Pref     int64
	Name     string
	Points   int64
}

func (q *Queries) UpdateShopPoints(ctx context.Context, arg UpdateShopPointsParams) error {
	_, err := q.db.ExecContext(ctx, updateShopPoints,
		arg.GameType,
		arg.Sid,
		arg.Pref,
		arg.Name,
		arg.Points,
	)
	return err
}