package modules

import (
	"context"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/core/db"
	"eamold/services/core/models"
	"eamold/utils"
)

type ModulePcbtracker struct {
	name string
	db   *db.Queries
}

func NewModulePcbtracker(db *db.Queries) *ModulePcbtracker {
	return &ModulePcbtracker{
		name: "pcbtracker",
		db:   db,
	}
}

func (m *ModulePcbtracker) Name() string {
	return m.name
}

func (m *ModulePcbtracker) Url() *string {
	return nil
}

func (m *ModulePcbtracker) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "pcbtracker":
		{
			switch elm.Method {
			case "alive":
				return m.alive(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModulePcbtracker) alive(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_PcbTracker_Alive

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	const ENABLE_PCBID_CHECK = false
	if ENABLE_PCBID_CHECK {
		status, err := m.db.GetPcbidStatus(context.TODO(), elm.SourceId)

		if err != nil {
			return nil, fmt.Errorf("invalid pcbid")
		}

		if status == int64(models.PcbidStatusBlacklisted) {
			return nil, fmt.Errorf("blacklisted pcbid")
		}
	}

	expire := 0 // TODO: Set to a non-0 value?

	return &models.Response_PcbTracker_Alive{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
		Expire:  expire,
	}, nil
}
