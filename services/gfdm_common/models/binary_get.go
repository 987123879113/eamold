package models

import "encoding/xml"

type Request_Binary_Get struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Data []struct {
		Code string `xml:"string,attr"`
	} `xml:"data"`
}
