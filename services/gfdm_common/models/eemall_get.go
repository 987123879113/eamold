package models

import "encoding/xml"

type Request_Eemall_Get struct {
	CardId *string `xml:"cardid,attr"`
}

type Response_Eemall_Get struct {
	XMLName xml.Name

	Status   int `xml:"status,attr"`
	NowPoint int `xml:"now_point,attr"`
	AddPoint int `xml:"add_point,attr"`

	Items []Response_Eemall_Get_Item `xml:"item"`
}

type Response_Eemall_Get_Item struct {
	Num int `xml:"num,attr"`
}
