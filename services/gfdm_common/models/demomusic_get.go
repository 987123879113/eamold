package models

import "encoding/xml"

type Response_DemoMusic_Get struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	MusicIDs string `xml:"musicids,attr"`
}
