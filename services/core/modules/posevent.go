package modules

import (
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/core/models"
)

type ModulePosevent struct {
	name string
}

func NewModulePosevent() *ModulePosevent {
	return &ModulePosevent{
		name: "posevent",
	}
}

func (m *ModulePosevent) Name() string {
	return m.name
}

func (m *ModulePosevent) Url() *string {
	return nil
}

func (m *ModulePosevent) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "posevent":
		{
			switch elm.Method {
			case "income":
				return m.income(elm)
			case "sale":
				return m.sale(elm)
			case "sales":
				return m.sales(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModulePosevent) income(elm internal_models.MethodXmlElement) (any, error) {
	return &models.Response_PosEvent_Income{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}

func (m *ModulePosevent) sale(elm internal_models.MethodXmlElement) (any, error) {
	return &models.Response_PosEvent_Sale{
		XMLName: xml.Name{Local: elm.Module},
		Method:  "sale",
	}, nil
}

func (m *ModulePosevent) sales(elm internal_models.MethodXmlElement) (any, error) {
	return &models.Response_PosEvent_Sales{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}
