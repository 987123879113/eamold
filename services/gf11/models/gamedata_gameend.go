package models

import "encoding/xml"

type Request_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Id struct {
		MachineSerialId string `xml:"sid,attr"`
		Pcbid           string `xml:"pcbid,attr"`
	} `xml:"id"`

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

	ShopRank struct {
		Point    int    `xml:"point,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"shoprank"`

	Session struct {
		Flag int `xml:"flag,attr"`
	} `xml:"session"`

	IR struct {
		IrAll int `xml:"irall,attr"`
		IrCom int `xml:"ircom,attr"`
	} `xml:"ir"`

	LandData struct {
		Round int `xml:"round,attr"`
	} `xml:"landdata"`

	Players []struct {
		Number int    `xml:"no,attr"`
		GdId   *int   `xml:"gdid,attr"`
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
			Clear int `xml:"clear,attr"`
			Time  int `xml:"time,attr"`
			MNR   int `xml:"mnr,attr"`
			NNR   int `xml:"nnr,attr"`
			Miss  int `xml:"miss,attr"`

			Stages []struct {
				Number  int `xml:"no,attr"`
				SeqMode int `xml:"seqmode,attr"`
				Flags   int `xml:"flags,attr"`
				Perc    int `xml:"perc,attr"`
				Score   int `xml:"score,attr"`
				Clear   int `xml:"clear,attr"`
				Skill   int `xml:"skill,attr"`
				Combo   int `xml:"combo,attr"`
				IrAll   int `xml:"irall,attr"`
				IrCom   int `xml:"ircom,attr"`
			} `xml:"stage"`
		} `xml:"play"`

		Land struct {
			Team   int `xml:"team,attr"`
			Point  int `xml:"point,attr"`
			SPoint int `xml:"spoint,attr"`
		} `xml:"land"`
	} `xml:"player"`
}

type Response_GameData_GameEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System   Response_System                    `xml:"system"`
	Players  []Response_GameData_GameEnd_Player `xml:"player"`
	LandData Response_GameData_GameEnd_LandData `xml:"landdata"`
}

type Response_GameData_GameEnd_Player_Stage struct {
	All   int `xml:"all,attr"`
	Order int `xml:"order,attr"`
	Best  int `xml:"best,attr"`
}

type Response_GameData_GameEnd_Player struct {
	Number     int                                      `xml:"no,attr"`
	Skill      int                                      `xml:"skill,attr"`
	SkillOrder int                                      `xml:"skill_order,attr"`
	SkillAll   int                                      `xml:"skill_all,attr"`
	Stages     []Response_GameData_GameEnd_Player_Stage `xml:"stage"`
}

type Response_GameData_GameEnd_LandData_Land struct {
	// If LandData.Session == 1
	OrgAreas  int `xml:"org_areas,attr"`
	OrgHidden int `xml:"org_hidden,attr"`
	Areas     int `xml:"areas,attr"`
	Hidden    int `xml:"hidden,attr"`

	// If LandData.Session == 0
	OrgPoints int `xml:"org_points,attr"`
	OrgNext   int `xml:"org_next,attr"`
	Point     int `xml:"point,attr"`
	Next      int `xml:"next,attr"`
	Order     int `xml:"order,attr"`
	All       int `xml:"all,attr"`
	Mag       int `xml:"mag,attr"`
}

type Response_GameData_GameEnd_LandData struct {
	From    int                                     `xml:"from,attr"`
	To      int                                     `xml:"to,attr"`
	Session int                                     `xml:"session,attr"`
	Land    Response_GameData_GameEnd_LandData_Land `xml:"land"`
}
