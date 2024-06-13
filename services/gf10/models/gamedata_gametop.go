package models

import "encoding/xml"

type Request_GameData_GameTop struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	// If this element exists in the request then it was a session mode play
	Session *struct{} `xml:"session"`

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

	ExData  []Response_GameData_GameTop_ExData `xml:"exdata"`
	IRData  Response_GameData_GameTop_IRData   `xml:"irdata"`
	Players []Response_GameData_GameTop_Player `xml:"player"`
}

type Response_GameData_GameTop_IRData struct {
	Round uint `xml:"round,attr"`
}

type Response_GameData_GameTop_ExData struct {
	Round      int    `xml:"round,attr"`
	ExId       int    `xml:"exid,attr"`
	Skill      int    `xml:"skill,attr"`
	Parameters string `xml:"parameters,attr"`
	Vacant     int    `xml:"vacant,attr"`
	Open       int    `xml:"open,attr"`
	Close      int    `xml:"close,attr"`
}

type Response_GameData_GameTop_Player_SkillPerc struct {
	SeqMode int    `xml:"seqmode,attr"`
	Values  string `xml:",innerxml"`
}

type Response_GameData_GameTop_Player struct {
	Number     int                                          `xml:"no,attr"`
	Recovery   int                                          `xml:"recovery,attr"`
	SkillSeqs  string                                       `xml:"skillseqs"`
	MusicSeqs  string                                       `xml:"musicseqs"`
	SkillPercs []Response_GameData_GameTop_Player_SkillPerc `xml:"skillpercs"`

	Ex     Response_GameData_GameTop_Player_Ex       `xml:"ex"`
	Ir     Response_GameData_GameTop_Player_IR       `xml:"ir"`
	Course []Response_GameData_GameTop_Player_Course `xml:"course"`
}

type Response_GameData_GameTop_Player_Ex struct {
	New int `xml:"new,attr"`
}

type Response_GameData_GameTop_Player_IR struct {
	New int `xml:"new,attr"`
}

type Response_GameData_GameTop_Player_Course struct {
	Class    int    `xml:"class,attr"`
	MusicIds string `xml:"musicids,attr"`
	Seqs     string `xml:"seqs,attr"`
	Diffs    string `xml:"diffs,attr"`
}
