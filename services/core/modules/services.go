package modules

import (
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/internal/services_manager"
	"eamold/services/core/models"
)

type ModuleServices struct {
	name    string
	manager services_manager.ServicesManager
}

func NewModuleServices(manager services_manager.ServicesManager) *ModuleServices {
	return &ModuleServices{
		name:    "services",
		manager: manager,
	}
}

func (m *ModuleServices) Name() string {
	return m.name
}

func (m *ModuleServices) Url() *string {
	return nil
}

func (m *ModuleServices) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "services":
		{
			switch elm.Method {
			case "get":
				return m.get(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleServices) get(elm internal_models.MethodXmlElement) (any, error) {
	modules := m.manager.GetRegisteredServiceModules(elm.Model)

	if modules == nil {
		return nil, fmt.Errorf("service is not registered")
	}

	serviceItems := make([]models.Response_Services_Get_Item, 0, len(modules))
	for _, entry := range modules {
		serviceItems = append(serviceItems, models.Response_Services_Get_Item{
			Name: entry.Name(),
			Url:  m.manager.Url(elm.Model.Model(), entry),
		})
	}

	return &models.Response_Services_Get{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		Expire: 600,
		Mode:   "operation",
		Items:  serviceItems,
	}, nil
}
