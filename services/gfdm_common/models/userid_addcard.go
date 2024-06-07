package models

import "encoding/xml"

type Request_UserId_AddCard struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Model *string `xml:"model,attr"`

	User string `xml:"user,attr"`
	Pass string `xml:"pass,attr"`
	Card string `xml:"card,attr"`
}

type Response_UserId_AddCard struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`
}
