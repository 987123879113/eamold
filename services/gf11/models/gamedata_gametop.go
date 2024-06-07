package models

import "encoding/xml"

type Request_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	// If this element exists in the request then it was a session mode play
	Session *struct{} `xml:"session"`

	Players []struct {
		Number   int  `xml:"no,attr"`
		GdId     *int `xml:"gdid,attr"`
		Recovery int  `xml:"recovery,attr"`
	} `xml:"player"`
}

type Response_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System   Response_System                    `xml:"system"`
	IRData   Response_GameData_GameTop_IR       `xml:"irdata"`
	LandData Response_GameData_GameTop_LandData `xml:"landdata"`
	Players  []Response_GameData_GameTop_Player `xml:"player"`
}

type Response_GameData_GameTop_Player_SkillPerc struct {
	SeqMode int    `xml:"seqmode,attr"`
	Values  string `xml:",innerxml"`
}

type Response_GameData_GameTop_Player struct {
	Number   int `xml:"no,attr"`
	Recovery int `xml:"recovery,attr"`

	SkillSeqs  string                                       `xml:"skillseqs"`
	MusicSeqs  string                                       `xml:"musicseqs"`
	SkillPercs []Response_GameData_GameTop_Player_SkillPerc `xml:"skillpercs"`
}

type Response_GameData_GameTop_IR struct {
	All uint `xml:"all,attr"`
	Com uint `xml:"com,attr"`
}

type Response_GameData_GameTop_LandData_Land struct {
	Team   int `xml:"team,attr"`
	Hidden int `xml:"hidden,attr"`
}

type Response_GameData_GameTop_LandData struct {
	Round int                                       `xml:"round,attr"`
	Land  []Response_GameData_GameTop_LandData_Land `xml:"land"`
}
