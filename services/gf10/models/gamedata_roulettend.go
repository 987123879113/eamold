package models

import "encoding/xml"

type Request_GameData_RoulettEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Players []Request_GameData_RoulettEnd_Player `xml:"player"`
}

type Request_GameData_RoulettEnd_Player struct {
	CardId string `xml:"cardid,attr"`

	Puzzle struct {
		Number int `xml:"no,attr"`
		Flags  int `xml:"flags,attr"`
		Hidden int `xml:"hidden,attr"`
	} `xml:"puzzle"`
}

type Response_GameData_RoulettEnd struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`
}
