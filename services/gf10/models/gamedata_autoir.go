package models

import "encoding/xml"

type Request_GameData_AutoIR struct {
	XMLName         xml.Name
	Method          string `xml:"method,attr"`
	MachineSerialId string `xml:"sid,attr"`

	Env struct {
		Round    int    `xml:"round,attr"`
		Category int    `xml:"category,attr"`
		ShopName string `xml:"shopname,attr"`
		Pref     int    `xml:"pref,attr"`
	} `xml:"env"`

	Players []struct {
		Number   int    `xml:"no,attr"`
		CardId   string `xml:"cardid,attr"`
		Class    int    `xml:"class,attr"`
		Name     string `xml:"name,attr"`
		Pass     string `xml:"pass,attr"`
		Score    int    `xml:"score,attr"`
		IrPass   string `xml:"irpass,attr"`
		IrRegId  string `xml:"irregid,attr"`
		Combo    int    `xml:"combo,attr"`
		Commands int    `xml:"commands,attr"`
	} `xml:"player"`
}

type Response_GameData_AutoIR struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Player []Response_GameData_AutoIR_Player `xml:"player"`
}

type Response_GameData_AutoIR_Player struct {
	Number int `xml:"no,attr"`
	Result int `xml:"result,attr"`
}
