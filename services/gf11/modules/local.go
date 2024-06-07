package modules

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	internal_models "eamold/internal/models"
	"eamold/services/gf11/db"
	"eamold/services/gf11/models"
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
			case "dataget":
				return m.gamedata_dataget(elm)
			case "autoir":
				return m.gamedata_autoir(elm)
			case "rend":
				return m.gamedata_rend(elm)
			case "cardcheck":
				return m.gamedata_cardcheck(elm)
			case "cardregist":
				return m.gamedata_cardregist(elm)
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

func (m *moduleLocal) gamedata_cardcheck(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_CardCheck
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	r := &models.Response_GameData_CardCheck{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},
	}

	profile, err := m.db.GetProfileByCardId(context.TODO(),
		db.GetProfileByCardIdParams{
			GameType: m.gameType,
			Cardid:   request.Card.Id,
		},
	)
	if err == nil {
		r.Card = &models.Response_GameData_CardCheck_Card{
			Status: constants.CardStatusSuccess,
			GdId:   profile.Gdid,
			Pass:   profile.Pass,
			Skill:  profile.Skill,
			Color:  getSkillColor(profile.Skill),
		}
	} else if err == sql.ErrNoRows {
		r.Card = &models.Response_GameData_CardCheck_Card{
			Status: constants.CardStatusNew,
		}
	} else {
		r.Card = &models.Response_GameData_CardCheck_Card{
			Status: constants.CardStatusError,
		}
	}

	return r, nil
}

func (m *moduleLocal) gamedata_cardregist(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_CardRegist
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	if request.Card == nil || request.Card.Id == "" {
		return nil, fmt.Errorf("cardregist: card info not provided")
	}

	profile, err := m.db.CreateCardProfile(context.TODO(), db.CreateCardProfileParams{
		GameType:   m.gameType,
		Gdid:       int64(rand.Int31()),
		Cardid:     request.Card.Id,
		Irid:       request.Card.IrId,
		Name:       request.Card.Name,
		Pass:       request.Card.Pass,
		Type:       int64(request.Card.Type),
		UpdateFlag: int64(request.Card.Update),
		PuzzleNo:   int64(request.Card.PuzzleNo),
		Recovery:   int64(request.Card.Recovery),
	})

	if err != nil {
		return nil, fmt.Errorf("cardregist: couldn't register profile: %v", err)
	}

	return &models.Response_GameData_CardRegist{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Card: models.Response_GameData_CardRegist_Card{
			Status: int(constants.CardStatusSuccess),
			GdId:   int(profile.Gdid),
		},
	}, nil
}
func (m *moduleLocal) gamedata_dataget(elm internal_models.MethodXmlElement) (any, error) {
	/*
		IR-related info:

		＃1 2004/04/22（水）15:00　～　2004/05/18（火）15:00
		https://web.archive.org/web/20040813013409fw_/http://www.konami.co.jp/am/gfdm/gf11dm10/inrank/inrank_kitei1.html

		＃2 2004/05/19（水）15:00　～　2004/06/15（火）15:00
		https://web.archive.org/web/20040813013413fw_/http://www.konami.co.jp/am/gfdm/gf11dm10/inrank/inrank_kitei2.html

		＃3 2004/06/16（水）15:00　～　2004/07/13（火）15:00
		https://web.archive.org/web/20040813013417fw_/http://www.konami.co.jp/am/gfdm/gf11dm10/inrank/inrank_kitei3.html

		＃4 2004/07/14（水）15:00　～　2004/08/03（火）15:00
		https://web.archive.org/web/20040813013420fw_/http://www.konami.co.jp/am/gfdm/gf11dm10/inrank/inrank_kitei4.html

		＃5 2004/08/04（水）15:00　～　2004/09/07（火）15:00
		#5 was session mode only
		https://web.archive.org/web/20040813013424fw_/http://www.konami.co.jp/am/gfdm/gf11dm10/inrank/inrank_kitei5.html
	*/

	var request models.Request_GameData_DataGet
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	favorites, err := m.db.GetFavorites(context.TODO(),
		db.GetFavoritesParams{
			GameType: m.gameType,
			Limit:    int64(request.Favorite.Count),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("dataget: GetFavorites: %v", err)
	}

	lands := make([]models.Response_GameData_DataGet_LandData_Land, 5)
	for i := range len(lands) {
		lands[i] = models.Response_GameData_DataGet_LandData_Land{
			Team:   i,
			Area:   1,  // TODO: What's this?
			Hidden: 31, // Everything unlocked
		}
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
		GameType: m.gameType,
		Sid:      request.MachineSerialId,
	})
	if err != nil {
		myShop = db.GetShopBySerialRow{
			Pref:   int64(request.ShopRank.Pref),
			Name:   "",
			Points: 0,
		}
	}

	shopMyOrder, err := m.db.GetShopRank(context.TODO(), db.GetShopRankParams{
		GameType: m.gameType,
		Sid:      request.MachineSerialId,
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
			From:   (time.Now().Year() * 10000) + (1 * 100) + 1,   // YYYY0101
			To:     (time.Now().Year() * 10000) + (12 * 100) + 31, // YYYY1231
		},
		ShopRank: models.Response_GameData_DataGet_ShopRank{
			From: (time.Now().Year() * 10000) + (1 * 100) + 1,   // YYYY0101
			To:   (time.Now().Year() * 10000) + (12 * 100) + 31, // YYYY1231
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
			All: 4,
			Com: 4,
		},
		LandData: models.Response_GameData_DataGet_LandData{
			Round: 4,
			Land:  lands,
		},
	}, nil
}

func (m *moduleLocal) gamedata_gametop(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_GameTop
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	players := make([]models.Response_GameData_GameTop_Player, len(request.Players))
	for playerIdx, player := range request.Players {
		if player.GdId == nil {
			players[playerIdx] = models.Response_GameData_GameTop_Player{
				Number: player.Number,
			}
		} else {
			playerData, err := m.db.GetProfileByGdId(context.TODO(), db.GetProfileByGdIdParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get gdid %v profile: %v", player.GdId, err)
			}

			allMaxSkillPoints, err := m.db.GetAllMaxSkillPointsByGdid(context.TODO(), db.GetAllMaxSkillPointsByGdidParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get gdid %v seq stats: %v", player.GdId, err)
			}

			allScores, err := m.db.GetSeqStatsByGdid(context.TODO(), db.GetSeqStatsByGdidParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
			})
			if err != nil {
				return nil, fmt.Errorf("gametop: could not get gdid %v seq stats: %v", player.GdId, err)
			}

			// TODO: How do these differ?
			skillSeqs := make([]byte, 0x100)
			musicSeqs := make([]byte, 0x100)

			for i := range 0x100 {
				skillSeqs[i] = 'x'
				musicSeqs[i] = 'x'
			}

			for _, v := range allMaxSkillPoints {
				skillSeqs[v.Netid] = byte(v.SeqMode + '0')
				musicSeqs[v.Netid] = byte(v.SeqMode + '0')
			}

			skillPercsStrs := make([][]string, 9)
			for seq := range len(skillPercsStrs) {
				skillPercsStrs[seq] = make([]string, 0x100)
				for i := range 0x100 {
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
			}
		}
	}

	lands := make([]models.Response_GameData_GameTop_LandData_Land, 5)
	for i := range 5 {
		lands[i] = models.Response_GameData_GameTop_LandData_Land{
			Team:   i,
			Hidden: 31, // Everything unlocked
		}
	}

	return &models.Response_GameData_GameTop{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		IRData: models.Response_GameData_GameTop_IR{
			All: 4,
			Com: 4,
		},

		LandData: models.Response_GameData_GameTop_LandData{
			Round: 4,
			Land:  lands,
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
			gdid := -1
			if player.GdId != nil {
				gdid = *player.GdId
			}

			m.db.SaveScore(context.TODO(), db.SaveScoreParams{
				GameType: m.gameType,
				Gdid:     int64(gdid),
				Netid:    int64(stageIdMap[stage.Number].NetId),
				Courseid: int64(stageIdMap[stage.Number].CourseId),
				SeqMode:  int64(stage.SeqMode),
				Flags:    int64(stage.Flags),
				Score:    int64(stage.Score),
				Clear:    int64(stage.Clear),
				Combo:    int64(stage.Combo),
				Skill:    int64(stage.Skill),
				Perc:     int64(stage.Perc),
			})
		}

		if player.GdId != nil {
			m.db.RecalculateTotalSkillPointsForGdid(context.TODO(), db.RecalculateTotalSkillPointsForGdidParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
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
			if stageIdMap[stage.Number].NetId != -1 {
				netid := stageIdMap[stage.Number].NetId

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

				if player.GdId != nil {
					best_, err := m.db.IsBestScore(context.TODO(), db.IsBestScoreParams{
						GameType: m.gameType,
						Netid:    int64(netid),
						SeqMode:  int64(stage.SeqMode),
						Gdid:     int64(*player.GdId),
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
			} else if stageIdMap[stage.Number].CourseId != -1 {
				courseid := stageIdMap[stage.Number].CourseId

				all, err := m.db.GetCourseTotalPlayCount(context.TODO(), db.GetCourseTotalPlayCountParams{
					GameType: m.gameType,
					Courseid: int64(courseid),
					SeqMode:  int64(stage.SeqMode),
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: GetCourseTotalPlayCount: %v", err)
				}

				order, err := m.db.GetCourseCurrentRank(context.TODO(), db.GetCourseCurrentRankParams{
					GameType: m.gameType,
					Courseid: int64(courseid),
					SeqMode:  int64(stage.SeqMode),
					Score:    int64(stage.Score),
				})
				if err != nil {
					return nil, fmt.Errorf("gameend: GetCourseCurrentRank: %v", err)
				}

				best := 0

				if player.GdId != nil {
					best_, err := m.db.GetCourseBestRank(context.TODO(), db.GetCourseBestRankParams{
						GameType: m.gameType,
						Courseid: int64(courseid),
						SeqMode:  int64(stage.SeqMode),
						Gdid:     int64(*player.GdId),
					})
					if err != nil {
						return nil, fmt.Errorf("gameend: GetCourseBestRank: %v", err)
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

		playerCount, err := m.db.GetPlayerCount(context.TODO(), m.gameType)
		if err != nil {
			return nil, fmt.Errorf("gameend: GetPlayerCount: %v", err)
		}

		if player.GdId != nil {
			playerSkill_, err := m.db.GetPlayerSkill(context.TODO(), db.GetPlayerSkillParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
			})
			if err != nil {
				return nil, fmt.Errorf("gameend: GetPlayerSkill: %v", err)
			}

			playerRank_, err := m.db.GetPlayerRank(context.TODO(), db.GetPlayerRankParams{
				GameType: m.gameType,
				Gdid:     int64(*player.GdId),
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
		LandData: models.Response_GameData_GameEnd_LandData{
			From:    20040101, // Don't set the year too high or it breaks the display graphics in-game
			To:      20041231,
			Session: 0,
		},
		Players: playerResponses,
	}, nil
}

func (m *moduleLocal) gamedata_rend(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GameData_Rend
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	for _, player := range request.Players {
		m.db.UpsertPuzzle(context.TODO(), db.UpsertPuzzleParams{
			GameType: m.gameType,
			Gdid:     int64(player.GdId),
			PuzzleNo: int64(player.Puzzle.Number),
			Flags:    int64(player.Puzzle.Flags),
			Hidden:   int64(player.Puzzle.Hidden),
		})
	}

	return &models.Response_GameData_Rend{
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
