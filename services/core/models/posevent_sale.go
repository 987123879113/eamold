package models

import "encoding/xml"

type Response_PosEvent_Sale struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Approve uint `xml:"approve,attr"` // printed as a hex value
	Margin  uint `xml:"margin,attr"`  // (expected) minimum of 1000?
	Dummy   uint `xml:"dummy,attr"`
}
