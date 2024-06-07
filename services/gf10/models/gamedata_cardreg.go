package models

import "encoding/xml"

type Request_GameData_CardReg struct {
	Card *struct {
		Id       string `xml:"id,attr"`
		Name     string `xml:"name,attr"`
		Pass     string `xml:"pass,attr"`
		Type     uint64 `xml:"type,attr"`
		Update   uint64 `xml:"update,attr"`
		Recovery uint64 `xml:"recovery,attr"`
	} `xml:"card"`
}

type Response_GameData_CardReg struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`
}
