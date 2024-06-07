package modules

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/core/db"
	"eamold/services/core/models"
	"eamold/utils"
)

type ModuleNumbering struct {
	name string
	db   *db.Queries
}

func NewModuleNumbering(db *db.Queries) *ModuleNumbering {
	return &ModuleNumbering{
		name: "numbering",
		db:   db,
	}
}

func (m *ModuleNumbering) Name() string {
	return m.name
}

func (m *ModuleNumbering) Url() *string {
	return nil
}

func (m *ModuleNumbering) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "numbering":
		{
			switch elm.Method {
			case "assign":
				return m.assign(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleNumbering) assign(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_Numbering_Assign
	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	if request.Label != "cardv1" || (request.Format != nil && *request.Format != "card16m10") {
		return nil, fmt.Errorf("assign: unknown parameters provided %v", request)
	}

	var number string
	for range 100 {
		numberBytes := make([]byte, 8)
		if _, err := rand.Read(numberBytes); err != nil {
			return nil, fmt.Errorf("assign: number generation: %v", err)
		}

		number = hex.EncodeToString(numberBytes)

		used, err := m.db.GetCardNumberStatus(context.TODO(), db.GetCardNumberStatusParams{
			Label:  request.Label,
			Number: number,
		})

		if err != nil {
			return nil, fmt.Errorf("assign: number availability check: %v", err)
		}

		if used == 0 {
			break
		}

		number = ""
	}

	if len(number) == 0 {
		return nil, fmt.Errorf("assign: could not find available number")
	}

	err := m.db.AddUsedCardNumber(context.TODO(), db.AddUsedCardNumberParams{
		Label:  request.Label,
		Number: number,
	})

	if err != nil {
		return nil, fmt.Errorf("assign: insert used number: %v", err)
	}

	return &models.Response_Numbering_Assign{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		Number: number,
	}, nil
}
