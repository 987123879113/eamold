package modules

import (
	"context"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gf8/db"
	"eamold/services/gf8/models"
	"eamold/utils"
)

type moduleLocal struct {
	name     string
	db       *db.Queries
	gameType int64
}

func NewModuleLocal(db *db.Queries, gameType int) *moduleLocal {
	return &moduleLocal{
		name:     "local",
		db:       db,
		gameType: int64(gameType),
	}
}

func (m *moduleLocal) Name() string {
	return m.name
}

func (m *moduleLocal) Url() *string {
	return nil
}

func (m *moduleLocal) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "gamedata":
		{
			switch elm.Method {
			case "dataget":
				return m.gamedata_dataget(elm)
			case "gameend":
				return m.gamedata_gameend(elm)
			case "rankset":
				return m.gamedata_rankset(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *moduleLocal) gamedata_dataget(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_DataGet

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	favorites, err := m.db.GetFavorites(context.TODO(), db.GetFavoritesParams{
		GameType: m.gameType,
		Limit:    int64(request.Favorite.Count),
	})
	if err != nil {
		return nil, fmt.Errorf("dataget: GetFavorites: %v", err)
	}

	shops, err := m.db.GetRankedShops(context.TODO(), db.GetRankedShopsParams{
		GameType: m.gameType,
		Limit:    int64(request.ShopRank.Count),
	})
	if err != nil {
		return nil, fmt.Errorf("dataget: GetRankedShops: %v", err)
	}

	prefShops, err := m.db.GetRankedShopsByPref(context.TODO(), db.GetRankedShopsByPrefParams{
		GameType: m.gameType,
		Pref:     int64(request.ShopRank.Pref),
		Limit:    int64(request.ShopRank.PrefCount),
	})
	if err != nil {
		return nil, fmt.Errorf("dataget: GetRankedShopsByPref: %v", err)
	}

	shopNames := make([]string, len(shops))
	shopPrefs := make([]int64, len(shops))
	shopPoints := make([]int64, len(shops))
	for i, v := range shops {
		shopNames[i] = v.Name
		shopPrefs[i] = v.Pref
		shopPoints[i] = v.Points
	}

	prefShopNames := make([]string, 0, len(prefShops))
	prefShopPoints := make([]int64, 0, len(prefShops))
	for _, v := range prefShops {
		prefShopNames = append(prefShopNames, v.Name)
		prefShopPoints = append(prefShopPoints, v.Points)
	}

	return &models.Response_GameData_DataGet{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		Favorite: models.Response_GameData_DataGet_Favorite{
			Count:    len(favorites),
			MusicIDs: utils.GenerateListStringInt64(favorites),
			Round:    1,
		},
		ShopRank: models.Response_GameData_DataGet_ShopRank{
			Round: 1,
			Shop: models.Response_GameData_DataGet_ShopRank_Shop{
				Count:  len(shops),
				Names:  utils.GenerateListString(shopNames),
				Prefs:  utils.GenerateListStringInt64(shopPrefs),
				Points: utils.GenerateListStringInt64(shopPoints),
			},
			PrefShop: models.Response_GameData_DataGet_ShopRank_PrefShop{
				Count:  len(prefShops),
				Names:  utils.GenerateListString(prefShopNames),
				Points: utils.GenerateListStringInt64(prefShopPoints),
			},
		},
		Prize: models.Response_GameData_DataGet_Prize{
			Flag: 1, // Some kind of unlock flag but I don't understand what it unlocks
		},
	}, nil
}

func (m *moduleLocal) gamedata_gameend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameEnd

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	r := &models.Response_GameData_GameEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}

	if request.Favorite != nil {
		musicIds := utils.SplitListStringInt64(request.Favorite.MusicIds)

		if len(musicIds) != request.Favorite.Count {
			return nil, fmt.Errorf("gameend: invalid favorite music IDs count")
		}

		for _, musicid := range musicIds {
			m.db.UpdateFavoriteCount(context.TODO(), db.UpdateFavoriteCountParams{
				GameType: m.gameType,
				Musicid:  musicid,
			})
		}
	}

	if request.Rank != nil {
		musicIds := utils.SplitListStringInt64(request.Rank.MusicIDs)

		if len(musicIds) != request.Rank.StageCount {
			return nil, fmt.Errorf("gameend: invalid rank music IDs count")
		}

		for _, player := range request.Rank.Player {
			seqs := utils.SplitListStringInt64(player.Seqs)
			if len(seqs) != request.Rank.StageCount {
				return nil, fmt.Errorf("gameend: invalid rank seqs")
			}

			scores := utils.SplitListStringInt64(player.Scores)
			if len(scores) != request.Rank.StageCount {
				return nil, fmt.Errorf("gameend: invalid rank scores")
			}

			for i := range len(musicIds) {
				m.db.AddScore(context.TODO(), db.AddScoreParams{
					GameType: m.gameType,
					Musicid:  musicIds[i],
					Seq:      seqs[i],
					Score:    scores[i],
				})
			}
		}

		// This has to be done after all of the player data has been saved so the counts are accurate
		// TODO: Does this really not take seq into account?
		stageRanks, err := m.db.GetTotalRankedScoreCounts(context.TODO(), db.GetTotalRankedScoreCountsParams{
			GameType: m.gameType,
			Musicids: musicIds,
		})
		if err != nil {
			return nil, fmt.Errorf("gameend: GetTotalRankedScoreCounts: %v", err)
		}

		stageRanksMap := make(map[int64]int64, len(stageRanks))
		for _, v := range stageRanks {
			stageRanksMap[v.Musicid] = v.Count
		}

		stageRanksOrdered := make([]int64, len(stageRanks))
		for i, v := range musicIds {
			stageRanksOrdered[i] = stageRanksMap[v]
		}

		r.Rank.All = utils.GenerateListStringInt64(stageRanksOrdered)
		for _, player := range request.Rank.Player {
			seqs := utils.SplitListStringInt64(player.Seqs)
			if len(seqs) != request.Rank.StageCount {
				return nil, fmt.Errorf("gameend: invalid rank seqs")
			}

			scores := utils.SplitListStringInt64(player.Scores)
			if len(scores) != request.Rank.StageCount {
				return nil, fmt.Errorf("gameend: invalid rank scores")
			}

			order := make([]int64, len(musicIds))

			for i := range musicIds {
				scoreRank, err := m.db.GetScoreRank(context.TODO(), db.GetScoreRankParams{
					GameType: m.gameType,
					Musicid:  musicIds[i],
					Seq:      seqs[i],
					Score:    scores[i],
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: GetScoreRank: %v", err)
				}

				order[i] = scoreRank
			}

			r.Rank.Rank = append(r.Rank.Rank, models.Response_GameData_GameEnd_RankData_Rank{
				Number: player.Number,
				Order:  utils.GenerateListStringInt64(order),
			})
		}
	}

	if request.ShopRank != nil {
		m.db.UpdateShopPoints(context.TODO(), db.UpdateShopPointsParams{
			GameType: m.gameType,
			Pref:     int64(request.ShopRank.Pref),
			Name:     request.ShopRank.ShopName,
			Points:   int64(request.ShopRank.Point),
		})
	}

	return r, nil
}

func (m *moduleLocal) gamedata_rankset(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_RankSet

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	musicIds := utils.SplitListStringInt64(request.Rank.MusicIDs)

	if len(musicIds) != request.Rank.StageCount {
		return nil, fmt.Errorf("rankset: invalid rank music IDs count")
	}

	for _, player := range request.Rank.Player {
		seqs := utils.SplitListStringInt64(player.Seqs)
		if len(seqs) != request.Rank.StageCount {
			return nil, fmt.Errorf("rankset: invalid rank seqs")
		}

		scores := utils.SplitListStringInt64(player.Scores)
		if len(scores) != request.Rank.StageCount {
			return nil, fmt.Errorf("rankset: invalid rank scores")
		}

		flags := utils.SplitListStringInt64(player.Scores)
		if len(flags) != request.Rank.StageCount {
			return nil, fmt.Errorf("rankset: invalid rank flags")
		}

		for i := range len(seqs) {
			err := m.db.AddRankedScore(context.TODO(), db.AddRankedScoreParams{
				GameType: m.gameType,
				Musicid:  musicIds[i],
				Seq:      seqs[i],
				Score:    scores[i],
				Flags:    flags[i],
				Name:     player.Name,
			})
			if err != nil {
				return nil, fmt.Errorf("rankset: AddRankedScore: %v", err)
			}
		}
	}

	return &models.Response_GameData_RankSet{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}
