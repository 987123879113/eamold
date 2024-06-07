package providers

import (
	"context"
	"fmt"

	"eamold/services/gf11/db"
	common "eamold/services/gfdm_common/modules"
)

type userIdDataProvider struct {
	db       *db.Queries
	gameType int64
}

func NewUserIdDataProvider(db *db.Queries, gameType int) common.UserIdDataProvider {
	return &userIdDataProvider{
		db:       db,
		gameType: int64(gameType),
	}
}

func (m *userIdDataProvider) GetUserInfoFromCardId(ctx context.Context, cardId string) (common.UserIdUserInfo, error) {
	profile, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
		GameType: m.gameType,
		Cardid:   cardId,
	})
	if err != nil {
		return common.UserIdUserInfo{}, fmt.Errorf("GetUserInfoFromCardId: %v", err)
	}

	active := 1
	if profile.Expired == 1 {
		active = 0
	}

	return common.UserIdUserInfo{
		Name:   profile.Name,
		Active: active,
	}, nil
}
