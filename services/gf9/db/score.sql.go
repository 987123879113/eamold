// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: score.sql

package db

import (
	"context"
)

const addScore = `-- name: AddScore :exec
INSERT INTO gf9dm8_scores
(game_type, cardid, music_num, seq_mode, flags, encore, extra, score, clear, skill, combo)
VALUES
(?,?,?,?,?,?,?,?,?,?,?)
`

type AddScoreParams struct {
	GameType int64
	Cardid   string
	MusicNum int64
	SeqMode  int64
	Flags    int64
	Encore   int64
	Extra    int64
	Score    int64
	Clear    int64
	Skill    int64
	Combo    int64
}

func (q *Queries) AddScore(ctx context.Context, arg AddScoreParams) error {
	_, err := q.db.ExecContext(ctx, addScore,
		arg.GameType,
		arg.Cardid,
		arg.MusicNum,
		arg.SeqMode,
		arg.Flags,
		arg.Encore,
		arg.Extra,
		arg.Score,
		arg.Clear,
		arg.Skill,
		arg.Combo,
	)
	return err
}

const getAllMaxSkillPointsByCardId = `-- name: GetAllMaxSkillPointsByCardId :many
SELECT music_num, seq_mode, CAST(MAX(skill) AS INTEGER) AS ` + "`" + `skill` + "`" + `
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
AND music_num != -1
AND skill > 0
AND clear > 0
GROUP BY cardid, music_num
`

type GetAllMaxSkillPointsByCardIdParams struct {
	Cardid   string
	GameType int64
}

type GetAllMaxSkillPointsByCardIdRow struct {
	MusicNum int64
	SeqMode  int64
	Skill    int64
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
		if err := rows.Scan(&i.MusicNum, &i.SeqMode, &i.Skill); err != nil {
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

const getPlayerCount = `-- name: GetPlayerCount :one
SELECT COUNT(*)
FROM gf9dm8_card_profile
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
    FROM gf9dm8_card_profile
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
FROM gf9dm8_card_profile
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

const getSkillScoresByCardId = `-- name: GetSkillScoresByCardId :many
SELECT music_num
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
ORDER BY score DESC
LIMIT ?
`

type GetSkillScoresByCardIdParams struct {
	Cardid   string
	GameType int64
	Limit    int64
}

func (q *Queries) GetSkillScoresByCardId(ctx context.Context, arg GetSkillScoresByCardIdParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getSkillScoresByCardId, arg.Cardid, arg.GameType, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var music_num int64
		if err := rows.Scan(&music_num); err != nil {
			return nil, err
		}
		items = append(items, music_num)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSongCurrentRank = `-- name: GetSongCurrentRank :one
SELECT COUNT(*) + 1
FROM gf9dm8_scores
WHERE music_num = ?
AND game_type = ?
AND seq_mode = ?
AND skill > ?
`

type GetSongCurrentRankParams struct {
	MusicNum int64
	GameType int64
	SeqMode  int64
	Skill    int64
}

func (q *Queries) GetSongCurrentRank(ctx context.Context, arg GetSongCurrentRankParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSongCurrentRank,
		arg.MusicNum,
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
FROM gf9dm8_scores
WHERE music_num = ?
AND game_type = ?
AND seq_mode = ?
`

type GetSongTotalPlayCountParams struct {
	MusicNum int64
	GameType int64
	SeqMode  int64
}

func (q *Queries) GetSongTotalPlayCount(ctx context.Context, arg GetSongTotalPlayCountParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getSongTotalPlayCount, arg.MusicNum, arg.GameType, arg.SeqMode)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const isBestScore = `-- name: IsBestScore :one
SELECT CAST(CASE WHEN COUNT(*) > 0 THEN 0 ELSE 1 END AS INTEGER)
FROM gf9dm8_scores
WHERE cardid = ?
AND game_type = ?
AND music_num = ?
AND seq_mode = ?
AND score >= ?
`

type IsBestScoreParams struct {
	Cardid   string
	GameType int64
	MusicNum int64
	SeqMode  int64
	Score    int64
}

func (q *Queries) IsBestScore(ctx context.Context, arg IsBestScoreParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, isBestScore,
		arg.Cardid,
		arg.GameType,
		arg.MusicNum,
		arg.SeqMode,
		arg.Score,
	)
	var column_1 int64
	err := row.Scan(&column_1)
	return column_1, err
}

const recalculateTotalSkillPointsForCardId = `-- name: RecalculateTotalSkillPointsForCardId :exec
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
WHERE gf9dm8_card_profile.cardid = ?1
`

type RecalculateTotalSkillPointsForCardIdParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) RecalculateTotalSkillPointsForCardId(ctx context.Context, arg RecalculateTotalSkillPointsForCardIdParams) error {
	_, err := q.db.ExecContext(ctx, recalculateTotalSkillPointsForCardId, arg.Cardid, arg.GameType)
	return err
}
