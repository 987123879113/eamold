package models

import "encoding/xml"

type Request_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System struct {
		Status int `xml:"status,attr"`
	} `xml:"system"`

	Players []struct {
		Number   int    `xml:"no,attr"`
		CardId   string `xml:"cardid,attr"`
		Recovery int    `xml:"recovery,attr"`
	} `xml:"player"`
}

type Response_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Players []Response_GameData_GameTop_Player `xml:"player"`
}

type Response_GameData_GameTop_Player struct {
	Number    int    `xml:"no,attr"`
	Status    int    `xml:"status,attr"`
	Recovery  int    `xml:"recovery,attr"`
	SkillSeqs string `xml:"skillseq,attr"`
	MusicSeqs string `xml:"musicseq,attr"`
}
