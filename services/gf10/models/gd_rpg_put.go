package models

import "encoding/xml"

type Request_GdRpg_Put struct {
	Session string `xml:"session,attr"`

	Player []struct {
		CardId string `xml:"card_id,attr"`
		PModel string `xml:"pmodel,attr"`

		Skill []struct {
			Point int `xml:"point,attr"`
		} `xml:"skill"`
	} `xml:"player"`
}

type Response_GdRpg_Put struct {
	XMLName xml.Name

	Status int `xml:"status,attr"`
	Fault  int `xml:"fault,attr"`
}
