package models

import "encoding/xml"

type Response_Services_Get struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Expire int    `xml:"expire,attr,omitempty"`
	Mode   string `xml:"mode,attr"`
	Items  []Response_Services_Get_Item
}

type Response_Services_Get_Item struct {
	XMLName xml.Name `xml:"item"`
	Name    string   `xml:"name,attr"`
	Url     string   `xml:"url,attr"`
}
