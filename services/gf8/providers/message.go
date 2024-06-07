package providers

import (
	"context"

	"eamold/services/gf8/db"
	common "eamold/services/gfdm_common/modules"
)

type messageDataProvider struct {
	db       *db.Queries
	gameType int64
}

func NewMessageDataProvider(db *db.Queries, gameType int) common.MessageDataProvider {
	return &messageDataProvider{
		db:       db,
		gameType: int64(gameType),
	}
}

func (m *messageDataProvider) GetMessages(ctx context.Context) ([]string, error) {
	return m.db.GetMessages(context.TODO(), m.gameType)
}
