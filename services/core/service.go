package core

import (
	"database/sql"

	"eamold/internal/models"
	"eamold/internal/services_manager"
	"eamold/services/core/db"
	"eamold/services/core/modules"
)

const SERVICE_NAME = "core"

type service struct {
	manager services_manager.ServicesManager

	moduleMap services_manager.ModuleMap
}

func (s *service) AcceptsRequest(model models.ModelString) bool {
	return false
}

func (s *service) GetModuleMap() services_manager.ModuleMap {
	return s.moduleMap
}

func New(manager services_manager.ServicesManager, conn *sql.DB) services_manager.Service {
	db := db.New(conn)

	return &service{
		manager: manager,
		moduleMap: services_manager.ModuleMap{
			"services":   modules.NewModuleServices(manager),
			"posevent":   modules.NewModulePosevent(),
			"pcbtracker": modules.NewModulePcbtracker(db),
			"numbering":  modules.NewModuleNumbering(db),
		},
	}
}
