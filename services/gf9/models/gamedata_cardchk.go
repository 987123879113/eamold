package models

import (
	"encoding/xml"

	"eamold/services/gfdm_common/constants"
)

type Request_GameData_CardChk struct {
	Card struct {
		Id string `xml:"id,attr"`
	} `xml:"card"`
}

type Response_GameData_CardChk struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	System Response_System `xml:"system"`

	Card *Response_GameData_CardChk_Card `xml:"card,omitempty"`
}

type Response_GameData_CardChk_Card struct {
	Status constants.CardStatus `xml:"status,attr"` // if status = 2 then the rest of the fields are read
}
