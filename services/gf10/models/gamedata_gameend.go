package models

import "encoding/xml"

type Request_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Id struct {
		MachineSerialId string `xml:"sid,attr"`
		Pcbid           string `xml:"pcbid,attr"`
	} `xml:"id"`

	Favorite struct {
		Count  int    `xml:"nr,attr"`
		NetIds string `xml:"netids,attr"`
	} `xml:"favorite"`

	ShopRank struct {
		Point    int    `xml:"point,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"shoprank"`

	Session struct {
		Flag int `xml:"flag,attr"`

		Members []struct {
			Game   int    `xml:"game,attr"`
			CardId string `xml:"cardid,attr"`
			Name   string `xml:"name,attr"`
		} `xml:"member"`
	} `xml:"session"`

	StageData struct {
		Mode  int `xml:"mode,attr"`
		Count int `xml:"nr,attr"`

		Stages []struct {
			Number int `xml:"no,attr"`

			// When something == 2
			CourseID *int `xml:"crsid,attr"`

			// Otherwise...
			NetID    int `xml:"netid,attr"`
			Selected int `xml:"selected,attr"`
			Extra    int `xml:"extra,attr"`
		} `xml:"stage"`
	} `xml:"stagedata"`

	Players []struct {
		Number int    `xml:"no,attr"`
		CardId string `xml:"cardid,attr"`
		Name   string `xml:"name,attr"`

		Condition struct {
			Styles   int `xml:"styles,attr"`
			Hidden   int `xml:"hidden,attr"`
			Recovery int `xml:"recovery,attr"`
		} `xml:"condition"`

		Judge struct {
			Perfect int `xml:"perfect,attr"`
			Great   int `xml:"great,attr"`
			Good    int `xml:"good,attr"`
			Poor    int `xml:"poor,attr"`
			Miss    int `xml:"miss,attr"`
		} `xml:"judge"`

		Play struct {
			Time     int `xml:"time,attr"`
			Count    int `xml:"nr,attr"`
			Clear    int `xml:"clear,attr"`
			GameMode int `xml:"gamemode,attr"`
			MNR      int `xml:"mnr,attr"`
			NNR      int `xml:"nnr,attr"`

			Stages []struct {
				NetId   int `xml:"netid,attr"`
				SeqMode int `xml:"seqmode,attr"`
				Flags   int `xml:"flags,attr"`
				Score   int `xml:"score,attr"`
				Clear   int `xml:"clear,attr"`
				Skill   int `xml:"skill,attr"`
				Combo   int `xml:"combo,attr"`
				Encore  int `xml:"encore,attr"`
				Extra   int `xml:"extra,attr"`
				Perc    int `xml:"perc,attr"`
			} `xml:"stage"`
		} `xml:"play"`

		Ex struct {
			Round int `xml:"round,attr"`
			ExId  int `xml:"exid,attr"`
			Seen  int `xml:"seen,attr"`
			Clear int `xml:"clear,attr"`
		} `xml:"ex"`

		IR struct {
			Round int `xml:"round,attr"`
		} `xml:"ir"`
	} `xml:"player"`
}

type Response_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Players  []Response_GameData_GameEnd_Player `xml:"player"`
	ExData   Response_GameData_GameEnd_ExData   `xml:"exdata"`
	ShopRank Response_GameData_GameEnd_ShopRank `xml:"shoprank"`
}

type Response_GameData_GameEnd_Player_Stage struct {
	All   int `xml:"all,attr"`
	Order int `xml:"order,attr"`
	Best  int `xml:"best,attr"`
}

type Response_GameData_GameEnd_Player struct {
	Number       int                                      `xml:"no,attr"`
	Skill        int                                      `xml:"skill,attr"`
	SkillOrder   int                                      `xml:"skill_order,attr"`
	SkillOrderNr int                                      `xml:"skill_order_nr,attr"`
	Stages       []Response_GameData_GameEnd_Player_Stage `xml:"stage"`
}

type Response_GameData_GameEnd_ExData_Ex struct {
	Serial int `xml:"serial,attr"`
}

type Response_GameData_GameEnd_ExData struct {
	Vacant int                                 `xml:"vacant,attr"`
	Ex     Response_GameData_GameEnd_ExData_Ex `xml:"ex"`
}

type Response_GameData_GameEnd_ShopRank struct {
	Hidden int `xml:"hidden,attr"`
}
