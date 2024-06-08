// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: profile.sql

package db

import (
	"context"
)

const getProfileByCardId = `-- name: GetProfileByCardId :one
SELECT game_type, cardid, name, color, recovery, styles, hidden, expired
FROM gf8dm7puv_card_profile
WHERE cardid = ?
AND game_type = ?
`

type GetProfileByCardIdParams struct {
	Cardid   string
	GameType int64
}

func (q *Queries) GetProfileByCardId(ctx context.Context, arg GetProfileByCardIdParams) (Gf8dm7puvCardProfile, error) {
	row := q.db.QueryRowContext(ctx, getProfileByCardId, arg.Cardid, arg.GameType)
	var i Gf8dm7puvCardProfile
	err := row.Scan(
		&i.GameType,
		&i.Cardid,
		&i.Name,
		&i.Color,
		&i.Recovery,
		&i.Styles,
		&i.Hidden,
		&i.Expired,
	)
	return i, err
}

const updateCardProfile = `-- name: UpdateCardProfile :exec
INSERT INTO gf8dm7puv_card_profile (game_type, cardid, name, color, recovery, styles, hidden)
VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)
ON CONFLICT (game_type, cardid) DO
UPDATE SET color=?4, recovery=?5, styles=?6, hidden=?7
`

type UpdateCardProfileParams struct {
	GameType int64
	Cardid   string
	Name     string
	Color    int64
	Recovery int64
	Styles   int64
	Hidden   int64
}

func (q *Queries) UpdateCardProfile(ctx context.Context, arg UpdateCardProfileParams) error {
	_, err := q.db.ExecContext(ctx, updateCardProfile,
		arg.GameType,
		arg.Cardid,
		arg.Name,
		arg.Color,
		arg.Recovery,
		arg.Styles,
		arg.Hidden,
	)
	return err
}

const updatePuzzleProgress = `-- name: UpdatePuzzleProgress :exec
INSERT INTO gf8dm7puv_puzzles (game_type, cardid, number, flags, out)
VALUES (?1, ?2, ?3, ?4, ?5)
ON CONFLICT (game_type, cardid, number) DO
UPDATE SET number=?3, flags=?4, out=?5
`

type UpdatePuzzleProgressParams struct {
	GameType int64
	Cardid   string
	Number   int64
	Flags    int64
	Out      int64
}

func (q *Queries) UpdatePuzzleProgress(ctx context.Context, arg UpdatePuzzleProgressParams) error {
	_, err := q.db.ExecContext(ctx, updatePuzzleProgress,
		arg.GameType,
		arg.Cardid,
		arg.Number,
		arg.Flags,
		arg.Out,
	)
	return err
}
