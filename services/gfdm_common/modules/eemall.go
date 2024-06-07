package modules

import (
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
	"eamold/utils"
)

type ModuleEemall struct {
	name string
}

func NewModuleEemall() *ModuleEemall {
	return &ModuleEemall{
		name: "eemall",
	}
}

func (m *ModuleEemall) Name() string {
	return m.name
}

func (m *ModuleEemall) Url() *string {
	return nil
}

func (m *ModuleEemall) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "eemall":
		{
			switch elm.Method {
			case "get":
				return m.eemall_get(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleEemall) eemall_get(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_Eemall_Get

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	// TODO: This info should be pulled from a database ideally
	// Points are be displayed in GF10DM9 on card in

	items := []models.Response_Eemall_Get_Item{}
	for i := range 32 {
		items = append(items, models.Response_Eemall_Get_Item{
			Num: 0x100 + i,
		})
	}

	return &models.Response_Eemall_Get{
		XMLName: xml.Name{Local: elm.Module},

		Status:   0,
		NowPoint: 1000,
		AddPoint: 100,
		Items:    items,
	}, nil
}
