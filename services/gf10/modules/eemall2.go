package modules

import (
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gf10/db"
	"eamold/services/gf10/models"
	"eamold/utils"
)

type ModuleEemall2 struct {
	name     string
	db       *db.Queries
	gameType int64
}

func NewModuleEemall2(db *db.Queries, gameType int) *ModuleEemall2 {
	return &ModuleEemall2{
		name:     "eemall2",
		db:       db,
		gameType: int64(gameType),
	}
}

func (m *ModuleEemall2) Name() string {
	return m.name
}

func (m *ModuleEemall2) Url() *string {
	return nil
}

func (m *ModuleEemall2) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "gd_rpg":
		{
			switch elm.Method {
			case "put":
				return m.gd_rpg_put(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleEemall2) gd_rpg_put(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_GdRpg_Put

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	return &models.Response_GdRpg_Put{
		XMLName: xml.Name{Local: elm.Module},

		Status: 0,
		Fault:  0,
	}, nil
}
