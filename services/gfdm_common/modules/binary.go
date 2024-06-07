package modules

import (
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
	"eamold/utils"
)

type ModuleBinary struct {
	name string
}

func NewModuleBinary() *ModuleBinary {
	return &ModuleBinary{
		name: "binary",
	}
}

func (m *ModuleBinary) Name() string {
	return m.name
}

func (m *ModuleBinary) Url() *string {
	return nil
}

func (m *ModuleBinary) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "binary":
		{
			switch elm.Method {
			case "get":
				return m.get(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleBinary) get(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_Binary_Get

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	// TODO: What is this?

	return nil, nil
}
