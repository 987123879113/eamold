package modules

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gf9/db"
	"eamold/services/gf9/models"
	"eamold/services/gfdm_common/constants"

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
			case "cardchk":
				return m.gamedata_cardchk(elm)
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

func (m *moduleLocal) gamedata_cardchk(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_CardChk

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	r := &models.Response_GameData_CardChk{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},
	}

	_, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
		GameType: m.gameType,
		Cardid:   request.Card.Id,
	})
	if err == nil {
		r.Card = &models.Response_GameData_CardChk_Card{
			Status: constants.CardStatusSuccess,
		}
	} else if err == sql.ErrNoRows {
		r.Card = &models.Response_GameData_CardChk_Card{
			Status: constants.CardStatusNew,
		}
	} else {
		r.Card = &models.Response_GameData_CardChk_Card{
			Status: constants.CardStatusError,
		}
	}

	return r, nil
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

	shopNames := make([]string, len(shops))
	shopPrefs := make([]int64, len(shops))
	shopPoints := make([]int64, len(shops))
	for i, v := range shops {
		shopNames[i] = v.Name
		shopPrefs[i] = v.Pref
		shopPoints[i] = v.Points
	}

	prefShops, err := m.db.GetRankedShopsByPref(context.TODO(), db.GetRankedShopsByPrefParams{
		GameType: m.gameType,
		Pref:     int64(request.ShopRank.Pref),
		Limit:    int64(request.ShopRank.PrefCount),
	})
	if err != nil {
		return nil, fmt.Errorf("dataget: GetRankedShopsByPref: %v", err)
	}

	prefShopNames := make([]string, 0, len(prefShops))
	prefShopPoints := make([]int64, 0, len(prefShops))
	for _, v := range prefShops {
		prefShopNames = append(prefShopNames, v.Name)
		prefShopPoints = append(prefShopPoints, v.Points)
	}

	myShop, err := m.db.GetShopByPcbid(context.TODO(), db.GetShopByPcbidParams{
		Pcbid:    elm.SourceId,
		GameType: m.gameType,
	})
	if err != nil {
		myShop = db.Gf9dm8Shop{
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

	return &models.Response_GameData_DataGet{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Favorite: models.Response_GameData_DataGet_Favorite{
			Count:    len(favorites),
			MusicIDs: utils.GenerateListStringInt64(favorites),
			Round:    1,
		},
		ShopRank: models.Response_GameData_DataGet_ShopRank{
			Round:    1,
			MyHidden: 0, // TODO: What's this for?
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
	}, nil
}

func (m *moduleLocal) gamedata_gametop(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameTop

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	players := make([]models.Response_GameData_GameTop_Player, len(request.Players))
	for i, player := range request.Players {
		if player.CardId != "" {
			_, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
				Cardid:   player.CardId,
				GameType: m.gameType,
			})
			if err == sql.ErrNoRows {
				players[i] = models.Response_GameData_GameTop_Player{
					Number: player.Number,
					Status: 0,
				}
			} else if err != nil {
				return nil, fmt.Errorf("gametop: GetSkillScoresByCardId: %v", err)
			} else {
				allMaxSkillPoints, err := m.db.GetAllMaxSkillPointsByCardId(context.TODO(), db.GetAllMaxSkillPointsByCardIdParams{
					Cardid:   player.CardId,
					GameType: m.gameType,
				})
				if err != nil {
					return nil, fmt.Errorf("gametop: could not get gdid %v seq stats: %v", player.CardId, err)
				}

				// TODO: How do these differ?
				skillSeqs := make([]byte, 0x80)
				musicSeqs := make([]byte, 0x80)

				for i := range 0x80 {
					skillSeqs[i] = 'x'
					musicSeqs[i] = 'x'
				}

				for _, v := range allMaxSkillPoints {
					skillSeqs[v.MusicNum] = byte(v.SeqMode + '0')
					musicSeqs[v.MusicNum] = byte(v.SeqMode + '0')
				}

				players[i] = models.Response_GameData_GameTop_Player{
					Number:    player.Number,
					Status:    1, // must be 1 for values to be read
					Recovery:  int(player.Recovery),
					SkillSeqs: string(skillSeqs),
					MusicSeqs: string(musicSeqs),
				}
			}
		} else {
			players[i] = models.Response_GameData_GameTop_Player{
				Number: player.Number,
				Status: 0,
			}
		}
	}

	return &models.Response_GameData_GameTop{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Players: players,
	}, nil
}

func (m *moduleLocal) gamedata_gameend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameEnd

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	if request.Favorite != nil {
		musicNums := utils.SplitListStringInt64(request.Favorite.MusicNums)
		for _, musicNum := range musicNums {
			m.db.UpdateFavoriteCount(context.TODO(), db.UpdateFavoriteCountParams{
				GameType: m.gameType,
				Musicid:  int64(musicNum),
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

	for _, player := range request.Players {
		for _, stage := range player.Play.Stages {
			err := m.db.AddScore(context.TODO(), db.AddScoreParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				MusicNum: int64(stage.MusicNum),
				SeqMode:  int64(stage.SeqMode),
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

		if player.CardId != "" {
			// Create new card profile if one doesn't exist for the provided card ID, or just update the values
			err := m.db.UpdateCardProfile(context.TODO(), db.UpdateCardProfileParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Name:     player.Name,
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

			m.db.RecalculateTotalSkillPointsForCardId(context.TODO(), db.RecalculateTotalSkillPointsForCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
		}
	}

	// The responses should be generated after all of the score re-ratings are done in case there's more than 1 player
	playerResponses := []models.Response_GameData_GameEnd_Player{}
	for _, player := range request.Players {
		playerStageResponses := []models.Response_GameData_GameEnd_Player_Stage{}

		for _, stage := range player.Play.Stages {
			if stage.MusicNum != -1 {
				all, err := m.db.GetSongTotalPlayCount(context.TODO(), db.GetSongTotalPlayCountParams{
					GameType: m.gameType,
					MusicNum: int64(stage.MusicNum),
					SeqMode:  int64(stage.SeqMode),
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: GetSongTotalPlayCount: %v", err)
				}

				order, err := m.db.GetSongCurrentRank(context.TODO(), db.GetSongCurrentRankParams{
					GameType: m.gameType,
					MusicNum: int64(stage.MusicNum),
					SeqMode:  int64(stage.SeqMode),
					Skill:    int64(stage.Skill),
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: GetSongCurrentRank: %v", err)
				}

				best := 1

				if player.CardId != "" {
					best_, err := m.db.IsBestScore(context.TODO(), db.IsBestScoreParams{
						GameType: m.gameType,
						MusicNum: int64(stage.MusicNum),
						SeqMode:  int64(stage.SeqMode),
						Cardid:   player.CardId,
						Score:    int64(stage.Score),
					})
					if err != nil {
						return nil, fmt.Errorf("gameend: IsBestScore: %v", err)
					}

					best = int(best_)
				}

				playerStageResponses = append(playerStageResponses, models.Response_GameData_GameEnd_Player_Stage{
					All:   int(all),
					Order: int(order),
					Best:  int(best),
				})
			}
		}

		playerSkill := 0
		playerRank := 0

		if player.CardId != "" {
			playerSkill_, err := m.db.GetPlayerSkill(context.TODO(), db.GetPlayerSkillParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetPlayerSkill: %v", err)
			}

			playerRank_, err := m.db.GetPlayerRank(context.TODO(), db.GetPlayerRankParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetPlayerRank: %v", err)
			}

			playerSkill = int(playerSkill_)
			playerRank = int(playerRank_)
		}

		playerCount, err := m.db.GetPlayerCount(context.TODO(), m.gameType)
		if err != nil {
			return nil, fmt.Errorf("gameend: GetPlayerCount: %v", err)
		}

		playerResponses = append(playerResponses, models.Response_GameData_GameEnd_Player{
			Number:     player.Number,
			Skill:      int(playerSkill),
			SkillAll:   int(playerCount),
			SkillOrder: int(playerRank),
			Stages:     playerStageResponses,
		})
	}

	return &models.Response_GameData_GameEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Players: playerResponses,

		ShopRank: models.Response_GameData_GameEnd_ShopRank{
			Hidden: 0,
		},
	}, nil
}
