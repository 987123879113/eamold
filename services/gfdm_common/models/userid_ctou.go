package models

import "encoding/xml"

type Request_UserId_Ctou struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Model *string `xml:"model,attr"`
	Card  string  `xml:"card,attr"`
}

type Response_UserId_Ctou struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	User   string `xml:"user,attr"`
	Active int    `xml:"active,attr"`
}
