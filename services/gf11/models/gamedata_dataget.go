package models

import "encoding/xml"

type Request_GameData_DataGet struct {
	XMLName         xml.Name
	Method          string `xml:"method,attr"`
	MachineSerialId string `xml:"sid,attr"`

	Favorite struct {
		Count int `xml:"nr,attr"`
	} `xml:"favorite"`

	ShopRank struct {
		Count     int `xml:"nr,attr"`
		PrefCount int `xml:"pref_nr,attr"`
		Pref      int `xml:"pref,attr"`
	} `xml:"shoprank"`
}

type Response_GameData_DataGet struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System   Response_System                    `xml:"system"`
	Favorite Response_GameData_DataGet_Favorite `xml:"favorite"`
	ShopRank Response_GameData_DataGet_ShopRank `xml:"shoprank"`
	IR       Response_GameData_DataGet_IR       `xml:"ir"`
	LandData Response_GameData_DataGet_LandData `xml:"landdata"`
}

type Response_GameData_DataGet_Favorite struct {
	Count  int    `xml:"nr,attr"`
	NetIDs string `xml:"netids,attr"`
	From   int    `xml:"from,attr"`
	To     int    `xml:"to,attr"`
}

type Response_GameData_DataGet_IR struct {
	All uint `xml:"all,attr"`
	Com uint `xml:"com,attr"`
}

type Response_GameData_DataGet_ShopRank_Shop struct {
	Count    int    `xml:"nr,attr"`
	Names    string `xml:"names,attr"`
	Points   string `xml:"points,attr"`
	Prefs    string `xml:"pref,attr"`
	MyOrder  uint   `xml:"myorder,attr"`
	MyPoints uint   `xml:"mypoint,attr"`
}

type Response_GameData_DataGet_ShopRank_PrefShop struct {
	Count   int    `xml:"nr,attr"`
	Names   string `xml:"names,attr"`
	Points  string `xml:"points,attr"`
	MyOrder uint   `xml:"myorder,attr"`
}

type Response_GameData_DataGet_ShopRank struct {
	From     int                                         `xml:"from,attr"`
	To       int                                         `xml:"to,attr"`
	Shop     Response_GameData_DataGet_ShopRank_Shop     `xml:"shop"`
	PrefShop Response_GameData_DataGet_ShopRank_PrefShop `xml:"prefshop"`
}

type Response_GameData_DataGet_LandData_Land struct {
	Team   int `xml:"team,attr"` // 0-4
	Area   int `xml:"area,attr"`
	Hidden int `xml:"hidden,attr"`
}

type Response_GameData_DataGet_LandData struct {
	Round uint                                      `xml:"round,attr"`
	Land  []Response_GameData_DataGet_LandData_Land `xml:"land"`
}
