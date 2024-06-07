package models

import "encoding/xml"

type PcbidStatus int

const (
	PcbidStatusUnknown PcbidStatus = iota
	PcbidStatusValid
	PcbidStatusBlacklisted
)

type Request_PcbTracker_Alive struct {
	Model *string `xml:"model,attr"` // TODO: When is this actually set?
}

type Response_PcbTracker_Alive struct {
	XMLName xml.Name
	Method  string `xml:"method,attr"`

	Expire int `xml:"expire,attr"`
}
