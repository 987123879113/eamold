package models

import (
	"encoding/xml"
)

type Response struct {
	XMLName xml.Name `xml:"response"`
	Status  int      `xml:"status,attr"`
	Fault   int      `xml:"fault,attr,omitempty"`

	Body []any
}
