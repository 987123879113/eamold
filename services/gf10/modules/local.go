package modules

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"strings"

	internal_models "eamold/internal/models"
	"eamold/utils"

	"eamold/services/gf10/db"
	"eamold/services/gf10/models"
	"eamold/services/gfdm_common/constants"
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
			case "autoir":
				return m.gamedata_autoir(elm)
			case "roulettend":
				return m.gamedata_roulettend(elm)
			case "cardchk":
				return m.gamedata_cardchk(elm)
			case "cardreg":
				return m.gamedata_cardreg(elm)
			case "gametop":
				return m.gamedata_gametop(elm)
			case "gameend":
				return m.gamedata_gameend(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func getSkillColor(skill int64) int64 {
	switch skill / 10000 {
	case 0:
		return 0
	case 1:
		return 0
	case 2:
		return 1
	case 3:
		return 1
	case 4:
		return 2
	case 5:
		return 2
	case 6:
		return 3
	case 7:
		return 3
	case 8:
		return 4
	case 9:
		return 5
	case 10:
		return 6
	default:
		return 0
	}
}

func (m *moduleLocal) gamedata_cardchk(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_Gamedata_CardChk

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	r := &models.Response_Gamedata_CardChk{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},
	}

	profile, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
		Cardid:   request.Card.Id,
		GameType: m.gameType,
	})
	if err == nil {
		r.Card = &models.Response_Gamedata_CardChk_Card{
			Status: constants.CardStatusSuccess,
			Pass:   profile.Pass,
			Skill:  profile.Skill,
			Color:  getSkillColor(profile.Skill),
		}
	} else if err == sql.ErrNoRows {
		r.Card = &models.Response_Gamedata_CardChk_Card{
			Status: constants.CardStatusNew,
		}
	} else {
		r.Card = &models.Response_Gamedata_CardChk_Card{
			Status: constants.CardStatusError,
		}
	}

	return r, nil
}

func (m *moduleLocal) gamedata_cardreg(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_CardReg

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	if request.Card == nil || request.Card.Id == "" {
		return nil, fmt.Errorf("cardregist: card info not provided")
	}

	_, err := m.db.CreateCardProfile(context.TODO(), db.CreateCardProfileParams{
		GameType:   m.gameType,
		Cardid:     request.Card.Id,
		Name:       request.Card.Name,
		Pass:       request.Card.Pass,
		Type:       int64(request.Card.Type),
		UpdateFlag: int64(request.Card.Update),
		Recovery:   int64(request.Card.Recovery),
	})

	if err != nil {
		return nil, fmt.Errorf("cardregist: couldn't register profile: %v", err)
	}

	return &models.Response_GameData_CardReg{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},
	}, nil
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

	myShop, err := m.db.GetShopBySerial(context.TODO(), db.GetShopBySerialParams{
		Sid:      request.MachineSerialId,
		GameType: m.gameType,
	})
	if err != nil {
		myShop = db.GetShopBySerialRow{
			Pref:   int64(request.ShopRank.Pref),
			Name:   "",
			Points: 0,
		}
	}

	shopMyOrder, err := m.db.GetShopRank(context.TODO(), db.GetShopRankParams{
		Sid:      request.MachineSerialId,
		GameType: m.gameType,
	})
	if err != nil {
		shopMyOrder = 0
	}

	prefShopMyOrder, err := m.db.GetShopRankByPref(context.TODO(), db.GetShopRankByPrefParams{
		GameType: m.gameType,
		Sid:      request.MachineSerialId,
		Pref:     int64(request.ShopRank.Pref),
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
			Count:  len(favorites),
			NetIDs: utils.GenerateListStringInt64(favorites),
			Round:  1,
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
		IR: models.Response_GameData_DataGet_IR{
			Round: 1,
		},
	}, nil
}

func (m *moduleLocal) gamedata_gametop(elm internal_models.MethodXmlElement) (any, error) {
	/*
		EX Challenge info
		https://web.archive.org/web/20040810043410fw_/http://www.konami.co.jp/am/gfdm/gf10dm9/howto/index.html
		https://plaza.rakuten.co.jp/thm/42652/
		https://nickjager.hatenablog.com/entry/20121202/1354441854

		#1: Clear a song with an S rank or better while getting a full combo
		DM: 24000 available, GF: 16000 available
		Date: 2003/10/22 to 2003/11/4
		Date: 2003/11/28 to 2003/12/4

		#2: Clear a song that is level 40 or higher with a combo of 300 or higher
		DM: 20000 available, GF: 12000 available
		Date: 2003/11/5 to 2003/11/16
		Date: 2003/12/18 to 2003/12/23

		#3: Clear a song that is level 40 or higher with less than 10 (for DM)/20 (for GF) misses
		DM: 12000 available, GF: 6000 available
		Date: 2003/11/17 to 2003/12/4

		#4: Clear a song on Extreme difficulty in Extra mode with a full combo
		DM: 10000 available, GF: 5000 available
		Date: 2003/12/5 to 2003/12/17

		#5: Clear a song that is level 70 or higher with 700 (for DM)/350 (for GF) or more perfects
		DM: 9000 available?, GF: 4500 available?
		Date: 2003/12/24 to 2004/1/4

		#6: Clear a song that is level 65 or higher with an SS
		DM: 8000 available, GF: 4000 available?
		Date: 2004/1/11 to 2004/1/21
		Date: 2004/2/24 to 2004/3/?

		#7: Clear a song that is level 75 (for DM)/level 70 (for GF) or higher with less than 10 (for DM)/20 (for GF) misses
		DM: 4000 available, GF: 2000 available
		Date: 2004/1/27 to 2004/2/7

		#8: Clear Agnus Dei on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
		DM: 3000 available, GF: 1500 available
		Date: 2004/2/13 to 2004/2/23
		Date: 2004/3/31 to 2004/4/4

		#9: Clear MODEL DD5 on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
		DM: 1500 available, GF: 750 available
		Date: 2004/2/29 to 2004/3/11

		#10: Clear Timepiece Phase II on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
		DM: 750 available??, GF: 375 available??
		Date: 2004/?/? to 2004/4/21
	*/

	var request models.Request_GameData_GameTop

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	players := make([]models.Response_GameData_GameTop_Player, len(request.Players))
	for playerIdx, player := range request.Players {
		if player.CardId != "" {
			playerData, err := m.db.GetProfileByCardId(context.TODO(), db.GetProfileByCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get cardid %v profile: %v", player.CardId, err)
			}

			allMaxSkillPoints, err := m.db.GetAllMaxSkillPointsByCardId(context.TODO(), db.GetAllMaxSkillPointsByCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get cardid %v seq stats: %v", player.CardId, err)
			}

			// TODO: How do these differ?
			skillSeqs := make([]byte, 192)
			musicSeqs := make([]byte, 192)

			for i := range len(skillSeqs) {
				skillSeqs[i] = 'x'
				musicSeqs[i] = 'x'
			}

			for _, v := range allMaxSkillPoints {
				skillSeqs[v.Netid] = byte(v.SeqMode + '0')
				musicSeqs[v.Netid] = byte(v.SeqMode + '0')
			}

			allScores, err := m.db.GetSeqStatsByCardId(context.TODO(), db.GetSeqStatsByCardIdParams{
				GameType: m.gameType,
				Cardid:   playerData.Cardid,
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get cardid %v seq stats: %v", player.CardId, err)
			}

			skillPercsStrs := make([][]string, 9)
			for seq := range len(skillPercsStrs) {
				skillPercsStrs[seq] = make([]string, 192)
				for i := range len(skillPercsStrs) {
					skillPercsStrs[seq][i] = "00"
				}
			}

			for _, v := range allScores {
				if v.Perc >= 100 {
					skillPercsStrs[v.SeqMode-1][v.Netid] = "A0"
				} else if v.Perc >= 0 {
					skillPercsStrs[v.SeqMode-1][v.Netid] = fmt.Sprintf("%02d", v.Perc)
				}
			}

			skillPercs := make([]models.Response_GameData_GameTop_Player_SkillPerc, 9)
			for seq := range len(skillPercs) {
				skillPercs[seq].SeqMode = seq + 1
				skillPercs[seq].Values = strings.Join(skillPercsStrs[seq], "")
			}

			players[playerIdx] = models.Response_GameData_GameTop_Player{
				Number:     player.Number,
				Recovery:   int(playerData.Recovery),
				SkillSeqs:  string(skillSeqs),
				MusicSeqs:  string(musicSeqs),
				SkillPercs: skillPercs,

				Ex: models.Response_GameData_GameTop_Player_Ex{
					New: 1,
				},

				Ir: models.Response_GameData_GameTop_Player_IR{
					New: 1,
				},
			}
		} else {
			players[playerIdx] = models.Response_GameData_GameTop_Player{
				Number: player.Number,
			}
		}
	}

	return &models.Response_GameData_GameTop{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		IRData: models.Response_GameData_GameTop_IRData{
			Round: 1,
		},

		Players: players,
	}, nil
}

func (m *moduleLocal) gamedata_gameend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameEnd

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	type stageInfo struct {
		NetId    int
		CourseId int
	}

	stageIdMap := make(map[int]stageInfo, len(request.StageData.Stages))
	for _, v := range request.StageData.Stages {
		if v.CourseID != nil {
			stageIdMap[v.Number] = stageInfo{
				NetId:    -1,
				CourseId: *v.CourseID,
			}
		} else {
			stageIdMap[v.Number] = stageInfo{
				NetId:    v.NetID,
				CourseId: -1,
			}
		}
	}

	for _, player := range request.Players {
		for _, stage := range player.Play.Stages {
			m.db.SaveScore(context.TODO(), db.SaveScoreParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Netid:    int64(stage.NetId),
				SeqMode:  int64(stage.SeqMode),
				Flags:    int64(stage.Flags),
				Score:    int64(stage.Score),
				Clear:    int64(stage.Clear),
				Combo:    int64(stage.Combo),
				Skill:    int64(stage.Skill),
				Perc:     int64(stage.Perc),
			})
		}

		if player.CardId != "" {
			m.db.RecalculateTotalSkillPointsForCardId(context.TODO(), db.RecalculateTotalSkillPointsForCardIdParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
		}
	}

	// Modify shop name to meet constraints expected by game
	shopName := strings.Replace(request.ShopRank.ShopName, "?", "？", -1) // ASCII ? will break the dataget network response so disallow it
	if len([]rune(shopName)) > 16 {
		shopName = string([]rune(shopName)[:16])
	}

	err := m.db.UpdateShopPoints(context.TODO(), db.UpdateShopPointsParams{
		GameType: m.gameType,
		Sid:      request.Id.MachineSerialId,
		Pref:     int64(request.ShopRank.Pref),
		Name:     shopName,
		Points:   int64(request.ShopRank.Point),
	})
	if err != nil {
		log.Printf("shop err: %v\n", err)
	}

	// The responses should be generated after all of the score re-ratings are done in case there's more than 1 player
	playerResponses := []models.Response_GameData_GameEnd_Player{}
	for _, player := range request.Players {
		playerStageResponses := []models.Response_GameData_GameEnd_Player_Stage{}

		for _, stage := range player.Play.Stages {
			netid := stage.NetId

			all, err := m.db.GetSongTotalPlayCount(context.TODO(), db.GetSongTotalPlayCountParams{
				GameType: m.gameType,
				Netid:    int64(netid),
				SeqMode:  int64(stage.SeqMode),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetSongTotalPlayCount: %v", err)
			}

			order, err := m.db.GetSongCurrentRank(context.TODO(), db.GetSongCurrentRankParams{
				GameType: m.gameType,
				Netid:    int64(netid),
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
					Netid:    int64(netid),
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

		playerSkill := 0
		playerRank := 0

		playerCount, err := m.db.GetPlayerCount(context.TODO(), m.gameType)
		if err != nil {
			return nil, fmt.Errorf("gameend: GetPlayerCount: %v", err)
		}

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

		// TODO: Expose this somewhere
		const FORCE_UNLOCK = false
		if FORCE_UNLOCK {
			playerSkill += 120000
		}

		playerResponses = append(playerResponses, models.Response_GameData_GameEnd_Player{
			Number:       player.Number,
			Skill:        int(playerSkill),
			SkillOrder:   int(playerRank),
			SkillOrderNr: int(playerCount),
			Stages:       playerStageResponses,
		})
	}

	return &models.Response_GameData_GameEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Players: playerResponses,
	}, nil
}

func (m *moduleLocal) gamedata_roulettend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_RoulettEnd

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	for _, player := range request.Players {
		m.db.UpsertPuzzle(context.TODO(), db.UpsertPuzzleParams{
			GameType: m.gameType,
			Cardid:   player.CardId,
			PuzzleNo: int64(player.Puzzle.Number),
			Flags:    int64(player.Puzzle.Flags),
			Hidden:   int64(player.Puzzle.Hidden),
		})
	}

	return &models.Response_GameData_RoulettEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}

func (m *moduleLocal) gamedata_autoir(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_AutoIR

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	return &models.Response_GameData_AutoIR{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}
