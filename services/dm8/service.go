package dm8

import (
	"database/sql"

	"eamold/internal/models"
	"eamold/internal/services_manager"
	"eamold/services/core"
	core_modules "eamold/services/core/modules"
	"eamold/services/gf9/db"
	"eamold/services/gf9/modules"
	"eamold/services/gf9/providers"
	gfdm_constants "eamold/services/gfdm_common/constants"
	gfdm_common "eamold/services/gfdm_common/modules"
)

const SERVICE_NAME = "C38:*:*:*"
const GAME_TYPE = int(gfdm_constants.GameTypeDrums)

type service struct {
	manager services_manager.ServicesManager

	moduleMap services_manager.ModuleMap
}

func (s *service) AcceptsRequest(model models.ModelString) bool {
	return model.Model() == models.ModelString(SERVICE_NAME).Model()
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
			"binary": gfdm_common.NewModuleBinary(),

			"local": modules.NewModuleLocal(db, GAME_TYPE),

			"keepalive": core_modules.NewModuleConstant("keepalive", "pa=127.0.0.1&ga=127.0.0.1&ping=ping://127.0.0.1&ntp=ntp://162.159.200.123"),
		},
	}
}
