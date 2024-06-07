package modules

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gf8puv/db"
	"eamold/services/gf8puv/models"
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
			case "gametop":
				return m.gamedata_gametop(elm)
			case "gameend":
				return m.gamedata_gameend(elm)
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

	myShop, err := m.db.GetShopByPcbid(context.TODO(), db.GetShopByPcbidParams{
		GameType: m.gameType,
		Pcbid:    elm.SourceId,
	})
	if err != nil {
		myShop = db.Gf8dm7puvShop{
			GameType: m.gameType,
			Pref:     int64(request.ShopRank.Pref),
			Name:     "",
			Points:   0,
		}
	}

	shopMyOrder, err := m.db.GetShopRank(context.TODO(), db.GetShopRankParams{
		GameType: m.gameType,
		Name:     myShop.Name,
		Pref:     myShop.Pref,
	})
	if err != nil {
		shopMyOrder = 0
	}

	prefShopMyOrder, err := m.db.GetShopRankByPref(context.TODO(), db.GetShopRankByPrefParams{
		GameType: m.gameType,
		Name:     myShop.Name,
		Pref:     myShop.Pref,
	})
	if err != nil {
		prefShopMyOrder = 0
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
				Count:    len(shops),
				Names:    utils.GenerateListString(shopNames),
				Prefs:    utils.GenerateListStringInt64(shopPrefs),
				Points:   utils.GenerateListStringInt64(shopPoints),
				MyOrder:  uint(shopMyOrder),
				MyPoints: uint(myShop.Points),
			},
			PrefShop: models.Response_GameData_DataGet_ShopRank_PrefShop{
				Count:   len(prefShops),
				Names:   utils.GenerateListString(prefShopNames),
				Points:  utils.GenerateListStringInt64(prefShopPoints),
				MyOrder: uint(prefShopMyOrder),
			},
		},
		Prize: models.Response_GameData_DataGet_Prize{
			Flag: 1, // Some kind of unlock flag but I don't understand what it unlocks
		},
	}, nil
}

func (m *moduleLocal) gamedata_gametop(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameTop

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	players := make([]models.Response_GameData_GameTop_Player, len(request.Players))
	for i, v := range request.Players {
		_, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
			GameType: m.gameType,
			Cardid:   v.CardId,
		})
		if err == sql.ErrNoRows {
			players[i] = models.Response_GameData_GameTop_Player{
				Number: v.Number,
				Status: 0,
			}
		} else if err != nil {
			return nil, fmt.Errorf("gametop: GetProfileByCardId: %v", err)
		} else {
			skillMusicIds, err := m.db.GetSkillScoresByCardId(context.TODO(), db.GetSkillScoresByCardIdParams{
				GameType: m.gameType,
				Cardid:   v.CardId,
				Limit:    30,
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: GetSkillScoresByCardId: %v", err)
			}

			players[i] = models.Response_GameData_GameTop_Player{
				Number:   v.Number,
				Status:   1,
				Recovery: 2, // TODO: Does this change anything?
				SkillMid: utils.GenerateListStringInt64(skillMusicIds),
			}
		}
	}

	return &models.Response_GameData_GameTop{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		Players: players,
	}, nil
}

func (m *moduleLocal) gamedata_gameend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameEnd

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	if request.Favorite != nil {
		for _, music := range request.Favorite.Music {
			m.db.UpdateFavoriteCount(context.TODO(), db.UpdateFavoriteCountParams{
				GameType: m.gameType,
				Musicid:  int64(music.Id),
			})
		}
	}

	if request.ShopRank != nil {
		m.db.AddMachineToShop(context.TODO(), db.AddMachineToShopParams{
			GameType: m.gameType,
			Pcbid:    elm.SourceId,
			Pref:     int64(request.ShopRank.Pref),
			Name:     request.ShopRank.ShopName,
		})

		m.db.UpdateShopPoints(context.TODO(), db.UpdateShopPointsParams{
			GameType: m.gameType,
			Pref:     int64(request.ShopRank.Pref),
			Name:     request.ShopRank.ShopName,
			Points:   int64(request.ShopRank.Point),
		})
	}

	prevSkillByPlayer := make(map[int]int, len(request.Players))

	for _, player := range request.Players {
		if player.CardId != "" {
			skillPoints, err := m.db.GetSkillPointsByCardId(context.TODO(), db.GetSkillPointsByCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetSkillPointsByCardId: %v", err)
			}

			prevSkillByPlayer[player.Number] = int(skillPoints)

			// Create new card profile if one doesn't exist for the provided card ID, or just update the values
			err = m.db.UpdateCardProfile(context.TODO(), db.UpdateCardProfileParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Name:     player.Name,
				Color:    int64(player.Color),
				Recovery: int64(player.Condition.Recovery),
				Styles:   int64(player.Condition.Styles),
				Hidden:   int64(player.Condition.Hidden),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: UpdateCardProfile: %v", err)
			}

			err = m.db.UpdatePuzzleProgress(context.TODO(), db.UpdatePuzzleProgressParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Number:   int64(player.Puzzle.Number),
				Flags:    int64(player.Puzzle.Flags),
				Out:      int64(player.Puzzle.Out),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: UpdatePuzzleProgress: %v", err)
			}
		}

		for _, stage := range player.Play.Stages {
			err := m.db.AddScore(context.TODO(), db.AddScoreParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Musicid:  int64(stage.MusicId),
				Musicnum: int64(stage.MusicNum),
				Seq:      int64(stage.Seq),
				Flags:    int64(stage.Flags),
				Encore:   int64(stage.Encore),
				Extra:    int64(stage.Extra),
				Score:    int64(stage.Score),
				Clear:    int64(stage.Clear),
				Skill:    int64(stage.Skill),
				Combo:    int64(stage.Combo),
			})

			if err != nil {
				return nil, fmt.Errorf("gameend: AddScore: %v", err)
			}
		}
	}

	players := make([]models.Response_GameData_GameEnd_Player, len(request.Players))
	for i, player := range request.Players {
		stages := make([]models.Response_GameData_GameEnd_Player_Stage, len(player.Play.Stages))
		for i, stage := range player.Play.Stages {
			order, err := m.db.GetScoreRank(context.TODO(), db.GetScoreRankParams{
				GameType: m.gameType,
				Musicid:  int64(stage.MusicId),
				Seq:      int64(stage.Seq),
				Score:    int64(stage.Score),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetScoreRank: %v", err)
			}

			all, err := m.db.GetTotalPlayedCount(context.TODO(), db.GetTotalPlayedCountParams{
				GameType: m.gameType,
				Musicid:  int64(stage.MusicId),
				Seq:      int64(stage.Seq),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetTotalPlayedCount: %v", err)
			}

			best := 1
			if player.CardId != "" {
				best_, err := m.db.IsBestScore(context.TODO(), db.IsBestScoreParams{
					GameType: m.gameType,
					Cardid:   player.CardId,
					Musicid:  int64(stage.MusicId),
					Seq:      int64(stage.Seq),
					Score:    int64(stage.Score),
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: IsBestScore: %v", err)
				}

				best = int(best_)
			}

			stages[i] = models.Response_GameData_GameEnd_Player_Stage{
				All:   int(all),
				Order: int(order),
				Best:  best,
			}
		}

		players[i] = models.Response_GameData_GameEnd_Player{
			Number:    player.Number,
			Status:    0,
			SkillPrev: 0,
			Skill:     0,
			Stages:    stages,
		}

		if player.CardId != "" {
			skillPoints, err := m.db.GetSkillPointsByCardId(context.TODO(), db.GetSkillPointsByCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetSkillPointsByCardId: %v", err)
			}

			players[i].Status = 1
			players[i].SkillPrev = prevSkillByPlayer[player.Number]
			players[i].Skill = int(skillPoints)
		}
	}

	shopInfo, err := m.db.GetShop(context.TODO(), db.GetShopParams{
		GameType: m.gameType,
		Name:     request.ShopRank.ShopName,
		Pref:     int64(request.ShopRank.Pref),
	})
	if err != nil {
		return nil, fmt.Errorf("gameend: GetShop: %v", err)
	}

	shopRank, err := m.db.GetShopRank(context.TODO(), db.GetShopRankParams{
		GameType: m.gameType,
		Name:     request.ShopRank.ShopName,
		Pref:     int64(request.ShopRank.Pref),
	})
	if err != nil {
		return nil, fmt.Errorf("gameend: GetShopRank: %v", err)
	}

	return &models.Response_GameData_GameEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		Players: players,
		ShopRank: models.Response_GameData_GameEnd_ShopRank{
			Point: int(shopInfo.Points),
			Order: int(shopRank),
		},
	}, nil
}
