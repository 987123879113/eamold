package models

import "encoding/xml"

type Request_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Players []struct {
		Number   int    `xml:"no,attr"`
		CardId   string `xml:"cardid,attr"`
		Recovery int    `xml:"recovery,attr"`
	} `xml:"player"`
}

type Response_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Players []Response_GameData_GameTop_Player `xml:"player"`
}

type Response_GameData_GameTop_Player struct {
	Number   int    `xml:"no,attr"`
	Status   int    `xml:"status,attr"`
	Recovery int    `xml:"recovery,attr"`
	SkillMid string `xml:"skillmid,attr"` // up to 30 music IDs
}
