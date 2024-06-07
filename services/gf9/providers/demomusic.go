package providers

import (
	"context"
	"fmt"

	"eamold/services/gf9/db"
	common "eamold/services/gfdm_common/modules"
)

type demoMusicDataProvider struct {
	db       *db.Queries
	gameType int64
}

func NewDemoMusicDataProvider(db *db.Queries, gameType int) common.DemoMusicDataProvider {
	return &demoMusicDataProvider{
		db:       db,
		gameType: int64(gameType),
	}
}

func (m *demoMusicDataProvider) GetDemoMusic(ctx context.Context, limit int) ([]int64, error) {
	musicIds, err := m.db.GetDemoMusic(context.TODO(), db.GetDemoMusicParams{
		GameType: m.gameType,
		Limit:    int64(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("GetDemoMusic: %v", err)
	}

	if len(musicIds) < limit {
		suppMusicIds, err := m.db.GetFavorites(context.TODO(), db.GetFavoritesParams{
			GameType: m.gameType,
			Limit:    int64(limit),
		})

		if err != nil {
			return nil, fmt.Errorf("GetFavorites: %v", err)
		}

		existingMusicIds := make(map[int64]bool, len(musicIds))
		for _, mid := range musicIds {
			existingMusicIds[mid] = true
		}

		for _, mid := range suppMusicIds {
			if _, ok := existingMusicIds[mid]; !ok {
				musicIds = append(musicIds, mid)
			}
		}
	}

	return musicIds, nil
}
