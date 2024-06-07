package models

import "encoding/xml"

type Request_UserId_AddUser struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	// All of these are optional
	Model *string `xml:"model,attr"`
	User  *string `xml:"user,attr"`
	Pass  *string `xml:"pass,attr"`
	Card  *string `xml:"card,attr"`
	Flags *int    `xml:"aflag,attr"`
}

type Response_UserId_AddUser struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Alter string `xml:"alter,attr"`
}
