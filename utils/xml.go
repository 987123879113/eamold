package utils

import (
	"bytes"
	"encoding/xml"

	"github.com/beevik/etree"
)

func UnserializeEtreeElement(e *etree.Element, out any) error {
	if e == nil {
		return nil
	}

	var b bytes.Buffer
	e.WriteTo(&b, &etree.WriteSettings{})

	if err := xml.Unmarshal(b.Bytes(), out); err != nil {
		return err
	}

	return nil
}
