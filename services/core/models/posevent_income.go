package models

import "encoding/xml"

type Response_PosEvent_Income struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Dummy uint `xml:"dummy,attr"`
}
