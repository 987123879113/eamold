package modules

import (
	"context"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
	"eamold/utils"
)

type DemoMusicDataProvider interface {
	GetDemoMusic(ctx context.Context, limit int) ([]int64, error)
}

type ModuleDemoMusic struct {
	name string
	db   DemoMusicDataProvider
}

func NewModuleDemoMusic(db DemoMusicDataProvider) *ModuleDemoMusic {
	return &ModuleDemoMusic{
		name: "demomusic",
		db:   db,
	}
}

func (m *ModuleDemoMusic) Name() string {
	return m.name
}

func (m *ModuleDemoMusic) Url() *string {
	return nil
}

func (m *ModuleDemoMusic) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "demomusic":
		{
			switch elm.Method {
			case "get":
				return m.get(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleDemoMusic) get(elm internal_models.MethodXmlElement) (any, error) {
	demoMusicList, err := m.db.GetDemoMusic(context.TODO(), 5)
	if err != nil {
		return nil, fmt.Errorf("demomusic.get: %v", err)
	}

	return &models.Response_DemoMusic_Get{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		MusicIDs: utils.GenerateListStringInt64(demoMusicList),
	}, nil
}
