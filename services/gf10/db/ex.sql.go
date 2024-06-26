// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: ex.sql

package db

import (
	"context"
)

const awardSerialForExChallenge = `-- name: AwardSerialForExChallenge :exec
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
AND gf10dm9_ex_progress.serial = 0
`

type AwardSerialForExChallengeParams struct {
	GameType int64
	Cardid   string
	Exid     int64
}

func (q *Queries) AwardSerialForExChallenge(ctx context.Context, arg AwardSerialForExChallengeParams) error {
	_, err := q.db.ExecContext(ctx, awardSerialForExChallenge, arg.GameType, arg.Cardid, arg.Exid)
	return err
}

const getExChallengeClearCounts = `-- name: GetExChallengeClearCounts :many
SELECT exid, COUNT(exid) AS clear_count
FROM gf10dm9_ex_progress
WHERE game_type = ?
AND clear = 1
GROUP BY exid
`

type GetExChallengeClearCountsRow struct {
	Exid       int64
	ClearCount int64
}

func (q *Queries) GetExChallengeClearCounts(ctx context.Context, gameType int64) ([]GetExChallengeClearCountsRow, error) {
	rows, err := q.db.QueryContext(ctx, getExChallengeClearCounts, gameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExChallengeClearCountsRow
	for rows.Next() {
		var i GetExChallengeClearCountsRow
		if err := rows.Scan(&i.Exid, &i.ClearCount); err != nil {
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

const getExChallengeProgress = `-- name: GetExChallengeProgress :many
SELECT game_type, cardid, exid, clear, seen, serial
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?
`

type GetExChallengeProgressParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) GetExChallengeProgress(ctx context.Context, arg GetExChallengeProgressParams) ([]Gf10dm9ExProgress, error) {
	rows, err := q.db.QueryContext(ctx, getExChallengeProgress, arg.Cardid, arg.GameType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Gf10dm9ExProgress
	for rows.Next() {
		var i Gf10dm9ExProgress
		if err := rows.Scan(
			&i.GameType,
			&i.Cardid,
			&i.Exid,
			&i.Clear,
			&i.Seen,
			&i.Serial,
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

const getExChallengeSerial = `-- name: GetExChallengeSerial :one
SELECT serial
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?
AND exid = ?
`

type GetExChallengeSerialParams struct {
	Cardid   string
	GameType int64
	Exid     int64
}

func (q *Queries) GetExChallengeSerial(ctx context.Context, arg GetExChallengeSerialParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getExChallengeSerial, arg.Cardid, arg.GameType, arg.Exid)
	var serial int64
	err := row.Scan(&serial)
	return serial, err
}

const getLatestClearedExChallenge = `-- name: GetLatestClearedExChallenge :one
SELECT exid
FROM gf10dm9_ex_progress
WHERE cardid = ?
AND game_type = ?
AND clear = 1
ORDER BY exid DESC
`

type GetLatestClearedExChallengeParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) GetLatestClearedExChallenge(ctx context.Context, arg GetLatestClearedExChallengeParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, getLatestClearedExChallenge, arg.Cardid, arg.GameType)
	var exid int64
	err := row.Scan(&exid)
	return exid, err
}

const upsertExChallenge = `-- name: UpsertExChallenge :exec
INSERT INTO gf10dm9_ex_progress
(game_type, cardid, exid, seen, clear)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, cardid, exid) DO
UPDATE SET seen=?4, clear=?5
`

type UpsertExChallengeParams struct {
	GameType int64
	Cardid   string
	Exid     int64
	Seen     int64
	Clear    int64
}

func (q *Queries) UpsertExChallenge(ctx context.Context, arg UpsertExChallengeParams) error {
	_, err := q.db.ExecContext(ctx, upsertExChallenge,
		arg.GameType,
		arg.Cardid,
		arg.Exid,
		arg.Seen,
		arg.Clear,
	)
	return err
}
