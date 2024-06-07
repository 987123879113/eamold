package models

import "encoding/xml"

type Request_Shopinf_GetSServIp struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	SourceIP string `xml:"srcip,attr"`
}

type Response_Shopinf_GetSServIp struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	ShopServerIp Response_Shopinf_GetSServIp_ShopServerIp `xml:"sservip"`
}

type Response_Shopinf_GetSServIp_ShopServerIp struct {
	IpAddress string `xml:"ipaddr,attr"`
}
