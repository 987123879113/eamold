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
			Round: 10,
		},
	}, nil
}

func (m *moduleLocal) getExChallengeDatas(cardIds []string) ([]models.Response_GameData_GameTop_ExData, error) {
	/*
		EX Challenge info
		https://web.archive.org/web/20040810043410fw_/http://www.konami.co.jp/am/gfdm/gf10dm9/howto/index.html
		https://plaza.rakuten.co.jp/thm/42652/
		https://nickjager.hatenablog.com/entry/20121202/1354441854
	*/
	exChallengeClearCounts, err := m.db.GetExChallengeClearCounts(context.TODO(), m.gameType)
	if err != nil {
		return nil, err
	}

	exChallengeClearCountMap := make([]int, 10)
	for _, v := range exChallengeClearCounts {
		exChallengeClearCountMap[v.Exid] = int(v.ClearCount)
	}

	exdatas := []models.Response_GameData_GameTop_ExData{
		{
			/*
				#1: Clear a song with an S rank or better while getting a full combo
				DM: 24000 available, GF: 16000 available
				Date: 2003/10/22 to 2003/11/4
				Date: 2003/11/28 to 2003/12/4
			*/
			Round:      1,
			ExId:       0,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndFullCombo(models.ExRankS),
			Vacant:     []int{16000, 24000}[m.gameType] - exChallengeClearCountMap[0],
			Open:       2003102200,
			Close:      2003110400,
		},
		{
			/*
				#2: Clear a song that is level 40 or higher with a combo of 300 or higher
				DM: 20000 available, GF: 12000 available
				Date: 2003/11/5 to 2003/11/16
				Date: 2003/12/18 to 2003/12/23
			*/
			Round:      1,
			ExId:       1,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndMinCombo(40, 300),
			Vacant:     []int{12000, 20000}[m.gameType] - exChallengeClearCountMap[1],
			Open:       2003110500,
			Close:      2003111600,
		},
		{
			/*
				#3: Clear a song that is level 40 or higher with less than 10 (for DM)/20 (for GF) misses
				DM: 12000 available, GF: 6000 available
				Date: 2003/11/17 to 2003/12/4
			*/
			Round:      1,
			ExId:       2,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndMaxJudgementCount(40, models.ExJudgementMiss, []int{20, 10}[m.gameType]),
			Vacant:     []int{6000, 12000}[m.gameType] - exChallengeClearCountMap[2],
			Open:       2003111700,
			Close:      2003120400,
		},
		{
			/*
				#4: Clear a song on Extreme difficulty in Extra mode with a full combo
				DM: 10000 available, GF: 5000 available
				Date: 2003/12/5 to 2003/12/17
			*/
			Round:      1,
			ExId:       3,
			Skill:      0,
			Parameters: models.NewExDataChallengeOnStageAndDifficultyWithFullCombo(models.ExStageExtra, models.ExDifficultyExtreme),
			Vacant:     []int{5000, 10000}[m.gameType] - exChallengeClearCountMap[3],
			Open:       2003120500,
			Close:      2003121700,
		},
		{
			/*
				#5: Clear a song that is level 70 or higher with 700 (for DM)/350 (for GF) or more perfects
				DM: 9000 available?, GF: 4500 available?
				Date: 2003/12/24 to 2004/1/4
			*/
			Round:      1,
			ExId:       4,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndMinJudgementCount(70, models.ExJudgementPerfect, []int{350, 700}[m.gameType]),
			Vacant:     []int{4500, 9000}[m.gameType] - exChallengeClearCountMap[4],
			Open:       2003122400,
			Close:      2004010400,
		},
		{
			/*
				#6: Clear a song that is level 65 or higher with an SS
				DM: 8000 available, GF: 4000 available?
				Date: 2004/1/11 to 2004/1/21
				Date: 2004/2/24 to 2004/3/?
			*/
			Round:      1,
			ExId:       5,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndMinRank(65, models.ExRankSS),
			Vacant:     []int{4000, 8000}[m.gameType] - exChallengeClearCountMap[5],
			Open:       2004011100,
			Close:      2004012100,
		},
		{
			/*
				#7: Clear a song that is level 75 (for DM)/level 70 (for GF) or higher with less than 10 (for DM)/20 (for GF) misses
				DM: 4000 available, GF: 2000 available
				Date: 2004/1/27 to 2004/2/7
			*/
			Round:      1,
			ExId:       6,
			Skill:      0,
			Parameters: models.NewExDataChallengeMinLevelAndMaxJudgementCount([]int{70, 75}[m.gameType], models.ExJudgementMiss, []int{20, 10}[m.gameType]),
			Vacant:     []int{2000, 4000}[m.gameType] - exChallengeClearCountMap[6],
			Open:       2004012700,
			Close:      2004020700,
		},
		{
			/*
				#8: Clear extra stage on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
				DM: 3000 available, GF: 1500 available
				Date: 2004/2/13 to 2004/2/23
				Date: 2004/3/31 to 2004/4/4
			*/
			Round:      1,
			ExId:       7,
			Skill:      0,
			Parameters: models.NewExDataChallengeOnStageAndDifficultyWithMinPercent(models.ExStageExtra, models.ExDifficultyExtreme, []int{95, 93}[m.gameType]),
			Vacant:     []int{1500, 3000}[m.gameType] - exChallengeClearCountMap[7],
			Open:       2004021300,
			Close:      2004022300,
		},
		{
			/*
				#9: Clear encore stage on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
				DM: 1500 available, GF: 750 available
				Date: 2004/2/29 to 2004/3/11
			*/
			Round:      1,
			ExId:       8,
			Skill:      0,
			Parameters: models.NewExDataChallengeOnStageAndDifficultyWithMinPercent(models.ExStageEncore, models.ExDifficultyExtreme, []int{95, 93}[m.gameType]),
			Vacant:     []int{750, 1500}[m.gameType] - exChallengeClearCountMap[8],
			Open:       2004022900,
			Close:      2004031100,
		},
		{
			/*
				#10: Clear premium encore stage on Extreme difficulty with 93% (for DM)/95% (for GF) or higher perfects
				DM: 750 available??, GF: 375 available??
				Date: 2004/?/? to 2004/4/21
			*/
			Round:      1,
			ExId:       9,
			Skill:      0,
			Parameters: models.NewExDataChallengeOnStageAndDifficultyWithMinPercent(models.ExStagePremiumEncore, models.ExDifficultyExtreme, []int{95, 93}[m.gameType]),
			Vacant:     []int{375, 750}[m.gameType] - exChallengeClearCountMap[9], // ?
			Open:       2004040100,                                                // ?
			Close:      2004042100,
		},
	}

	// This is something custom for this server. The server will give all players the first ex challenge that hasn't been completed betwen the two players.
	// This event was originally a time-based event so all players always saw the same challenges, but that's not practical for a local server so this
	// solution was implemented.

	minExChallenge := 0

	if len(cardIds) > 0 {
		exChallengeClearMapTotal := make(map[int64]int64, 10)
		for i := range 10 {
			exChallengeClearMapTotal[int64(i)] = -1
		}

		for _, cardid := range cardIds {
			exChallengeProgress, err := m.db.GetExChallengeProgress(context.TODO(), db.GetExChallengeProgressParams{
				GameType: m.gameType,
				Cardid:   cardid,
			})
			if err != nil {
				return nil, err
			}

			exChallengeClearMap := make(map[int64]int64, 10)
			for _, v := range exChallengeProgress {
				exChallengeClearMap[v.Exid] = v.Clear
			}

			for i := range 10 {
				if exChallengeClearMapTotal[int64(i)] == -1 || exChallengeClearMap[int64(i)] < exChallengeClearMapTotal[int64(i)] {
					exChallengeClearMapTotal[int64(i)] = exChallengeClearMap[int64(i)]
				}
			}
		}

		for i := range 10 {
			if exChallengeClearMapTotal[int64(i)] != 1 {
				minExChallenge = i
				break
			}
		}
	}

	if minExChallenge+1 < len(exdatas) {
		return exdatas[:minExChallenge+1], nil
	}

	return exdatas, nil
}

func (m *moduleLocal) gamedata_gametop(elm internal_models.MethodXmlElement) (any, error) {
	// IR #10 course list
	// ref: https://plaza.rakuten.co.jp/kisekiyuki/diary/200401220000/
	// Not sure how accurate this really is but it's better than blank courses in-game
	ircourses := []models.Response_GameData_GameTop_Player_Course{
		{
			Class:    0,
			MusicIds: utils.GenerateListStringInt64([]int64{934, 922, 827, 910}),
			Seqs:     utils.GenerateListStringInt64([]int64{0, 0, 0, 0}),
			Diffs:    utils.GenerateListStringInt64([]int64{1, 2, 3, 4}),
		},
		{
			Class:    1,
			MusicIds: utils.GenerateListStringInt64([]int64{921, 913, 932, 510}),
			Seqs:     utils.GenerateListStringInt64([]int64{1, 1, 1, 1}),
			Diffs:    utils.GenerateListStringInt64([]int64{1, 2, 3, 4}),
		},
		{
			Class:    2,
			MusicIds: utils.GenerateListStringInt64([]int64{928, 15, 504, 619}),
			Seqs:     utils.GenerateListStringInt64([]int64{1, 2, 0, 1}),
			Diffs:    utils.GenerateListStringInt64([]int64{1, 2, 3, 4}),
		},
	}

	var request models.Request_GameData_GameTop

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	cardIds := []string{}
	for _, player := range request.Players {
		if player.CardId != "" {
			cardIds = append(cardIds, player.CardId)
		}
	}

	exdatas, err := m.getExChallengeDatas(cardIds)
	if err != nil {
		return nil, err
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

			exChallengeProgress, err := m.db.GetExChallengeProgress(context.TODO(), db.GetExChallengeProgressParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
			})
			if err != nil {
				return nil, err
			}

			exChallengeProgressMap := make([]int64, 10)
			exChallengeIsNewMap := []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
			for _, v := range exChallengeProgress {
				exChallengeProgressMap[v.Exid] = v.Serial

				if v.Seen == 1 {
					exChallengeIsNewMap[v.Exid] = 0
				}
			}

			isNewExChallenge := 0
			if len(exdatas) > 0 {
				isNewExChallenge = exChallengeIsNewMap[exdatas[len(exdatas)-1].ExId]
			}

			players[playerIdx] = models.Response_GameData_GameTop_Player{
				Number:     player.Number,
				Recovery:   int(playerData.Recovery),
				SkillSeqs:  string(skillSeqs),
				MusicSeqs:  string(musicSeqs),
				SkillPercs: skillPercs,

				Ex: models.Response_GameData_GameTop_Player_Ex{
					// If this is set to 1 then the player is forced into the EX challenge screen every game start.
					// There should ideally be a table that tracks what EX challenges have been presented already,
					// and and the EX challenge progress also needs to be tracked somewhere (checks in gameend?)
					New:   isNewExChallenge,
					Value: utils.GenerateListStringInt64(exChallengeProgressMap),
				},

				Ir: models.Response_GameData_GameTop_Player_IR{
					New: 0, // Same for this
				},

				Course: ircourses,
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

		ExData: exdatas,
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

	// Hacky shit, I don't want to store the ex challenge settings in a table so just get the before and after ex challenge data lists to get the proper vacancy count for the response
	cardIds := []string{}
	for _, player := range request.Players {
		if player.CardId != "" {
			cardIds = append(cardIds, player.CardId)
		}
	}

	exdatasBefore, err := m.getExChallengeDatas(cardIds)
	if err != nil {
		return nil, err
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

		if player.Ex != nil {
			err := m.db.UpsertExChallenge(context.TODO(), db.UpsertExChallengeParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Exid:     int64(player.Ex.ExId),
				Seen:     int64(player.Ex.Seen),
				Clear:    int64(player.Ex.Clear),
			})
			if err != nil {
				return nil, err
			}

			if player.Ex.Clear == 1 {
				err := m.db.AwardSerialForExChallenge(context.TODO(), db.AwardSerialForExChallengeParams{
					GameType: m.gameType,
					Cardid:   player.CardId,
					Exid:     int64(player.Ex.ExId),
				})
				if err != nil {
					return nil, err
				}
			}
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

	err = m.db.UpdateShopPoints(context.TODO(), db.UpdateShopPointsParams{
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
		if player.CardId != "" {
			cardIds = append(cardIds, player.CardId)
		}

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

		exSerial := 0

		if player.CardId != "" {
			exSerial_, err := m.db.GetExChallengeSerial(context.TODO(), db.GetExChallengeSerialParams{
				GameType: m.gameType,
				Cardid:   player.CardId,
				Exid:     int64(player.Ex.ExId),
			})
			if err != nil && err != sql.ErrNoRows {
				return nil, err
			}

			exSerial = int(exSerial_)
		}

		playerResponses = append(playerResponses, models.Response_GameData_GameEnd_Player{
			Number:       player.Number,
			Skill:        int(playerSkill),
			SkillOrder:   int(playerRank),
			SkillOrderNr: int(playerCount),
			Stages:       playerStageResponses,
			Ex: models.Response_GameData_GameEnd_ExData_Ex{
				Serial: exSerial,
			},
		})
	}

	exdatasAfter, err := m.getExChallengeDatas(cardIds)
	if err != nil {
		return nil, err
	}

	currentVacancyCount := 0
	for _, v := range exdatasAfter {
		if v.ExId == exdatasBefore[len(exdatasBefore)-1].ExId {
			currentVacancyCount = v.Vacant
			break
		}
	}

	return &models.Response_GameData_GameEnd{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		System: models.Response_System{
			Status: 0,
		},

		Players: playerResponses,

		ExData: models.Response_GameData_GameEnd_ExData{
			Vacant: currentVacancyCount,
		},
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
