package dm7puv

import (
	"database/sql"

	"eamold/internal/models"
	"eamold/internal/services_manager"
	"eamold/services/core"
	"eamold/services/gf8puv/db"
	"eamold/services/gf8puv/modules"
	"eamold/services/gf8puv/providers"
	gfdm_constants "eamold/services/gfdm_common/constants"
	gfdm_common "eamold/services/gfdm_common/modules"
)

const SERVICE_NAME = "C07:*:B:*"
const GAME_TYPE = int(gfdm_constants.GameTypeDrums)

type service struct {
	manager services_manager.ServicesManager

	moduleMap services_manager.ModuleMap
}

func (s *service) AcceptsRequest(model models.ModelString) bool {
	serviceModel := models.ModelString(SERVICE_NAME)
	return model.Model() == serviceModel.Model() && model.Spec() == serviceModel.Spec()
}

func (s *service) GetModuleMap() services_manager.ModuleMap {
	return s.moduleMap
}

func New(manager services_manager.ServicesManager, conn *sql.DB) services_manager.Service {
	db := db.New(conn)

	return &service{
		manager: manager,
		moduleMap: services_manager.ModuleMap{
			"services":   manager.ResolveModule(core.SERVICE_NAME, "services"),
			"posevent":   manager.ResolveModule(core.SERVICE_NAME, "posevent"),
			"pcbtracker": manager.ResolveModule(core.SERVICE_NAME, "pcbtracker"),
			"numbering":  manager.ResolveModule(core.SERVICE_NAME, "numbering"),

			"userid": gfdm_common.NewModuleUserId(
				providers.NewUserIdDataProvider(db, GAME_TYPE),
			),
			"message": gfdm_common.NewModuleMessage(
				providers.NewMessageDataProvider(db, GAME_TYPE),
			),
			"demomusic": gfdm_common.NewModuleDemoMusic(
				providers.NewDemoMusicDataProvider(db, GAME_TYPE),
			),

			"local": modules.NewModuleLocal(db, GAME_TYPE),
		},
	}
}
