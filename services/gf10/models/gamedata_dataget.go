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
}

type Response_GameData_DataGet_Favorite struct {
	Count  int    `xml:"nr,attr"`
	NetIDs string `xml:"netids,attr"`
	Round  int    `xml:"round,attr"`
}

type Response_GameData_DataGet_IR struct {
	Round int `xml:"round,attr"`
}

type Response_GameData_DataGet_ShopRank_Shop struct {
	Count    int    `xml:"nr,attr"`
	Names    string `xml:"name,attr"`
	Points   string `xml:"point,attr"`
	Prefs    string `xml:"pref,attr"`
	MyOrder  uint   `xml:"myorder,attr"`
	MyPoints uint   `xml:"mypoint,attr"`
}

type Response_GameData_DataGet_ShopRank_PrefShop struct {
	Count   int    `xml:"nr,attr"`
	Names   string `xml:"name,attr"`
	Points  string `xml:"point,attr"`
	MyOrder uint   `xml:"myorder,attr"`
}

type Response_GameData_DataGet_ShopRank struct {
	Round    int                                         `xml:"round,attr"`
	MyHidden int                                         `xml:"myhidden,attr"`
	Shop     Response_GameData_DataGet_ShopRank_Shop     `xml:"shop"`
	PrefShop Response_GameData_DataGet_ShopRank_PrefShop `xml:"prefshop"`
}
