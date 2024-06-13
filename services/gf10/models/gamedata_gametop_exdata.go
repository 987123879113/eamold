package models

import (
	"eamold/utils"
)

type ExRank int

const (
	ExRankE ExRank = iota
	ExRankD
	ExRankC
	ExRankB
	ExRankA
	ExRankS
	ExRankSS
)

type ExDifficulty int

const (
	ExDifficultyPractice ExDifficulty = iota
	ExDifficultyBasic
	ExDifficultyAdvanced
	ExDifficultyExtreme
)

type ExStage int

const (
	ExStage1st ExStage = iota
	ExStage2nd
	ExStage3rd
	ExStage4th
	ExStage5th
	ExStageExtra
	ExStageEncore
	ExStagePremiumEncore
)

type ExJudgement int

const (
	_ ExJudgement = iota
	ExJudgementPerfect
	ExJudgementGreat
	ExJudgementGood
	ExJudgementPoor
	ExJudgementMiss
)

type ExCommand int

const (
	_ ExCommand = iota
	ExCommandSpeed
	ExCommandHiddenSudden
	ExCommandDark
	ExCommandReverse
	ExCommandPosition
)

type ExCommandSetting int

const (
	ExCommandSpeed_Off ExCommandSetting = iota

	ExCommandHiddenSudden_Off ExCommandSetting = iota
	ExCommandHiddenSudden_Hidden
	ExCommandHiddenSudden_Sudden
	ExCommandHiddenSudden_HiddenSudden

	ExCommandDark_Off ExCommandSetting = iota
	ExCommandDark_On

	ExCommandReverse_Off ExCommandSetting = iota
	ExCommandReverse_On

	ExCommandPosition_TypeA ExCommandSetting = iota
	ExCommandPosition_TypeB
	ExCommandPosition_TypeC
	ExCommandPosition_Off
)

func NewExDataChallengeMinLevelAndMinRank(level int, rank ExRank) string {
	// Clear a song that is level <level> or higher with an <rank>
	// 難度値<level>以上の曲をランク<rank>以上でクリア
	return utils.GenerateListStringInt64([]int64{
		0,
		int64(level),
		int64(rank),
	})
}

func NewExDataChallengeMinLevelAndMinCombo(level int, combo int) string {
	// Clear a song that is level X or higher with a combo of Y or higher
	// <level>以上の曲を<combo>コンボ以上でクリア
	return utils.GenerateListStringInt64([]int64{
		1,
		int64(level),
		int64(combo),
	})
}

func NewExDataChallengeOnStageAndDifficultyWithMinPercent(stage ExStage, difficulty ExDifficulty, perc int) string {
	// On stage <stage>, clear a song on <difficulty> difficulty with <percentage>% or higher perfects
	// <stage>ステージにて、<difficulty>の曲をPerfect <percentage>%以上でクリア
	return utils.GenerateListStringInt64([]int64{
		2,
		int64(stage),
		int64(difficulty),
		int64(perc),
	})
}

func NewExDataChallengeMinLevelAndMaxJudgementCount(level int, judgementMode ExJudgement, count int) string {
	// Clear a song that is level <level> or higher with less than <count> <judgementMode>
	// 難度値<level>以上の曲で<judgementMode>の判定数が<count>以下でクリア
	return utils.GenerateListStringInt64([]int64{
		3,
		int64(level),
		int64(judgementMode),
		int64(count),
	})
}

func NewExDataChallengeMinLevelAndMinJudgementCount(level int, judgementMode ExJudgement, count int) string {
	// Clear a song that is level X or higher with Y or more <grade measurement>
	// 難度値<level>以上の曲で<judgementMode>の判定数が<count>以上でクリア
	return utils.GenerateListStringInt64([]int64{
		4,
		int64(level),
		int64(judgementMode),
		int64(count),
	})
}

func NewExDataChallengeMinLevelAndFixedSetting(level int, command ExCommand, setting ExCommandSetting) string {
	// Clear a song that is level <level> or higher with <command> command set to <setting>
	// 難度値<level>以上の曲で<command>コマンドを<setting>にしてクリア
	return utils.GenerateListStringInt64([]int64{
		5,
		int64(level),
		int64(command),
		int64(setting),
	})
}

func NewExDataChallengeMinLevelAndFullCombo(rank ExRank) string {
	// Clear a song with a <rank> rank or better while getting a full combo
	// ランク<rank>以上でFull Comboクリア
	return utils.GenerateListStringInt64([]int64{
		6,
		int64(rank),
	})
}

func NewExDataChallengeOnStageAndDifficultyWithFullCombo(stage ExStage, difficulty ExDifficulty) string {
	// Clear a song with a <rank> rank or better while getting a full combo
	// <stage>ステージにて、<difficulty>の曲をFull Comboでクリア
	return utils.GenerateListStringInt64([]int64{
		7,
		int64(stage),
		int64(difficulty),
	})
}
