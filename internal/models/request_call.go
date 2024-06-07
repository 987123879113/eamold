package models

import (
	"encoding/xml"
	"strings"

	"github.com/beevik/etree"
)

type RequestCall struct {
	XMLName  xml.Name
	SourceId string      `xml:"srcid,attr"`
	Model    ModelString `xml:"model,attr"`
	Items    []byte      `xml:",innerxml"`

	XMLRaw string
}

type RequestCallItem struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`

	Content string `xml:",innerxml"`
}

func (m RequestCallItem) Method() string {
	method := ""

	for _, attr := range m.Attrs {
		if strings.ToLower(attr.Name.Local) == "method" {
			method = attr.Value
			break
		}
	}

	return method
}

type MethodXmlElement struct {
	Model    ModelString
	SourceId string
	Module   string
	Method   string

	Element *etree.Element
}
