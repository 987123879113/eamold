package models

import "encoding/xml"

type Request_GameData_CardRegist struct {
	Card *struct {
		Id       string `xml:"id,attr"`
		IrId     string `xml:"irid,attr"`
		Name     string `xml:"name,attr"`
		Pass     string `xml:"pass,attr"`
		Type     uint64 `xml:"type,attr"`
		Update   uint64 `xml:"update,attr"`
		PuzzleNo uint64 `xml:"pazzleno,attr"` // attribute name itself is typo'd
		Recovery uint64 `xml:"recovery,attr"`
	} `xml:"card"`
}

type Response_GameData_CardRegist struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Card Response_GameData_CardRegist_Card `xml:"card,omitempty"`
}

type Response_GameData_CardRegist_Card struct {
	Status int `xml:"status,attr"`
	GdId   int `xml:"gdid,attr"`
}
