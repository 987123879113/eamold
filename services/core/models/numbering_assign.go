package models

import "encoding/xml"

type Request_Numbering_Assign struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Label  string  `xml:"label,attr"` // cardv1
	Model  *string `xml:"model,attr"`
	Format *string `xml:"format,attr"` // card16m10
}

type Response_Numbering_Assign struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Number string `xml:"number,attr"`
}
