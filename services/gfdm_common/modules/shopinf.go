package modules

import (
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
	"eamold/utils"
)

type moduleShopinf struct {
	name          string
	serverAddress string
}

func NewModuleShopinf(serverAddress string) *moduleShopinf {
	return &moduleShopinf{
		name:          "shopinf",
		serverAddress: serverAddress,
	}
}

func (m *moduleShopinf) Name() string {
	return m.name
}

func (m *moduleShopinf) Url() *string {
	return nil
}

func (m *moduleShopinf) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "shopinf":
		{
			switch elm.Method {
			case "getsservip":
				return m.getsservip(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *moduleShopinf) getsservip(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_Shopinf_GetSServIp

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	return &models.Response_Shopinf_GetSServIp{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		ShopServerIp: models.Response_Shopinf_GetSServIp_ShopServerIp{
			IpAddress: m.serverAddress,
		},
	}, nil
}
