// This module is provided as a simple way to implement things like ntp or keepalive URLS in the services list
package modules

import (
	"fmt"

	internal_models "eamold/internal/models"
)

type ModuleConstant struct {
	name string
	url  string
}

func NewModuleConstant(name string, url string) *ModuleConstant {
	return &ModuleConstant{
		name: name,
		url:  url,
	}
}

func (m *ModuleConstant) Name() string {
	return m.name
}

func (m *ModuleConstant) Url() *string {
	return &m.url
}

func (m *ModuleConstant) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}
