package models

import "encoding/xml"

type Request_GameData_Rend struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Players []struct {
		GdId int `xml:"gdid,attr"`

		Puzzle struct {
			Number int `xml:"no,attr"`
			Flags  int `xml:"flags,attr"`
			Hidden int `xml:"hidden,attr"`
		} `xml:"puzzle"`
	} `xml:"player"`
}

type Response_GameData_Rend struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`
}
