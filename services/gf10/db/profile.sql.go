// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: profile.sql

package db

import (
	"context"
)

const createCardProfile = `-- name: CreateCardProfile :one
INSERT INTO gf10dm9_card_profile (game_type, cardid, name, pass, type, update_flag, recovery)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING game_type, cardid, name, pass, type, update_flag, recovery, skill, expired
`

type CreateCardProfileParams struct {
	GameType   int64
	Cardid     string
	Name       string
	Pass       string
	Type       int64
	UpdateFlag int64
	Recovery   int64
}

func (q *Queries) CreateCardProfile(ctx context.Context, arg CreateCardProfileParams) (Gf10dm9CardProfile, error) {
	row := q.db.QueryRowContext(ctx, createCardProfile,
		arg.GameType,
		arg.Cardid,
		arg.Name,
		arg.Pass,
		arg.Type,
		arg.UpdateFlag,
		arg.Recovery,
	)
	var i Gf10dm9CardProfile
	err := row.Scan(
		&i.GameType,
		&i.Cardid,
		&i.Name,
		&i.Pass,
		&i.Type,
		&i.UpdateFlag,
		&i.Recovery,
		&i.Skill,
		&i.Expired,
	)
	return i, err
}

const getProfileByCardId = `-- name: GetProfileByCardId :one
SELECT game_type, cardid, name, pass, type, update_flag, recovery, skill, expired
FROM gf10dm9_card_profile
WHERE cardid = ?
AND game_type = ?
`

type GetProfileByCardIdParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) GetProfileByCardId(ctx context.Context, arg GetProfileByCardIdParams) (Gf10dm9CardProfile, error) {
	row := q.db.QueryRowContext(ctx, getProfileByCardId, arg.Cardid, arg.GameType)
	var i Gf10dm9CardProfile
	err := row.Scan(
		&i.GameType,
		&i.Cardid,
		&i.Name,
		&i.Pass,
		&i.Type,
		&i.UpdateFlag,
		&i.Recovery,
		&i.Skill,
		&i.Expired,
	)
	return i, err
}

const getPuzzleByNumber = `-- name: GetPuzzleByNumber :one
SELECT game_type, cardid, puzzle_no, flags, hidden
FROM gf10dm9_puzzle
WHERE cardid = ?
AND game_type = ?
AND puzzle_no = ?
`

type GetPuzzleByNumberParams struct {
	Cardid   string
	GameType int64
	PuzzleNo int64
}

func (q *Queries) GetPuzzleByNumber(ctx context.Context, arg GetPuzzleByNumberParams) (Gf10dm9Puzzle, error) {
	row := q.db.QueryRowContext(ctx, getPuzzleByNumber, arg.Cardid, arg.GameType, arg.PuzzleNo)
	var i Gf10dm9Puzzle
	err := row.Scan(
		&i.GameType,
		&i.Cardid,
		&i.PuzzleNo,
		&i.Flags,
		&i.Hidden,
	)
	return i, err
}

const getRecentPuzzle = `-- name: GetRecentPuzzle :one
SELECT gp.game_type, gp.cardid, puzzle_no, flags, hidden, gcp.game_type, gcp.cardid, name, pass, type, update_flag, recovery, skill, expired
FROM gf10dm9_puzzle AS gp
INNER JOIN gf10dm9_card_profile AS gcp ON (gp.cardid = gcp.cardid)
WHERE gp.cardid = ?
AND gp.game_type = ?
AND gp.puzzle_no = gcp.puzzle_no
`

type GetRecentPuzzleParams struct {
	Cardid   string
	GameType int64
}

type GetRecentPuzzleRow struct {
	GameType   int64
	Cardid     string
	PuzzleNo   int64
	Flags      int64
	Hidden     int64
	GameType_2 int64
	Cardid_2   string
	Name       string
	Pass       string
	Type       int64
	UpdateFlag int64
	Recovery   int64
	Skill      int64
	Expired    int64
}

func (q *Queries) GetRecentPuzzle(ctx context.Context, arg GetRecentPuzzleParams) (GetRecentPuzzleRow, error) {
	row := q.db.QueryRowContext(ctx, getRecentPuzzle, arg.Cardid, arg.GameType)
	var i GetRecentPuzzleRow
	err := row.Scan(
		&i.GameType,
		&i.Cardid,
		&i.PuzzleNo,
		&i.Flags,
		&i.Hidden,
		&i.GameType_2,
		&i.Cardid_2,
		&i.Name,
		&i.Pass,
		&i.Type,
		&i.UpdateFlag,
		&i.Recovery,
		&i.Skill,
		&i.Expired,
	)
	return i, err
}

const updateRecoveryCount = `-- name: UpdateRecoveryCount :exec
UPDATE gf10dm9_card_profile
SET recovery = ?
WHERE cardid = ?
AND game_type = ?
`

type UpdateRecoveryCountParams struct {
	Recovery int64
	Cardid   string
	GameType int64
}

func (q *Queries) UpdateRecoveryCount(ctx context.Context, arg UpdateRecoveryCountParams) error {
	_, err := q.db.ExecContext(ctx, updateRecoveryCount, arg.Recovery, arg.Cardid, arg.GameType)
	return err
}

const upsertPuzzle = `-- name: UpsertPuzzle :exec
INSERT INTO gf10dm9_puzzle
(game_type, cardid, puzzle_no, flags, hidden)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (cardid, puzzle_no) DO
UPDATE SET flags=?4, hidden=?5
`

type UpsertPuzzleParams struct {
	GameType int64
	Cardid   string
	PuzzleNo int64
	Flags    int64
	Hidden   int64
}

func (q *Queries) UpsertPuzzle(ctx context.Context, arg UpsertPuzzleParams) error {
	_, err := q.db.ExecContext(ctx, upsertPuzzle,
		arg.GameType,
		arg.Cardid,
		arg.PuzzleNo,
		arg.Flags,
		arg.Hidden,
	)
	return err
}
