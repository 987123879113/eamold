package models

import "encoding/xml"

type Response_Message_Get struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	String string `xml:"string,attr"`
}
