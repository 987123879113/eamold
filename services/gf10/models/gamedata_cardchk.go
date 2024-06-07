package models

import (
	"encoding/xml"

	"eamold/services/gfdm_common/constants"
)

type Request_Gamedata_CardChk struct {
	Card struct {
		Id string `xml:"id,attr"`
	} `xml:"card"`
}

type Response_Gamedata_CardChk struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Card *Response_Gamedata_CardChk_Card `xml:"card,omitempty"`
}

type Response_Gamedata_CardChk_Card struct {
	Status constants.CardStatus `xml:"status,attr"` // if status = 2 then the rest of the fields are read
	Pass   string               `xml:"pass,attr,omitempty"`
	Skill  int64                `xml:"skill,attr,omitempty"`
	Color  int64                `xml:"color,attr,omitempty"`
}
