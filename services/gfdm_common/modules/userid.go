package modules

import (
	"context"
	"encoding/xml"
	"fmt"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
	"eamold/utils"
)

// Since all of the SQL is being handled by sqlc, it's easier to just define an interface and
// have all of the games pass in their own data provider using their own SQL commands
type UserIdDataProvider interface {
	GetUserInfoFromCardId(ctx context.Context, cardId string) (UserIdUserInfo, error)
}

type UserIdUserInfo struct {
	Name   string
	Active int
}

type ModuleUserId struct {
	name string
	db   UserIdDataProvider
}

func NewModuleUserId(db UserIdDataProvider) *ModuleUserId {
	return &ModuleUserId{
		name: "userid",
		db:   db,
	}
}

func (m *ModuleUserId) Name() string {
	return m.name
}

func (m *ModuleUserId) Url() *string {
	return nil
}

func (m *ModuleUserId) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "userid":
		{
			switch elm.Method {
			case "ctou":
				return m.ctou(elm)
			case "adduser":
				return m.adduser(elm)
			case "addcard":
				return m.addcard(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleUserId) ctou(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_UserId_Ctou

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	profile, err := m.db.GetUserInfoFromCardId(context.TODO(), request.Card)
	if err != nil {
		return nil, fmt.Errorf("ctou: GetProfileByCardId: %v", err)
	}

	return &models.Response_UserId_Ctou{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		// What does this actually do? The flags don't seem to change anything in-game
		User:   profile.Name,
		Active: profile.Active,
	}, nil
}

func (m *ModuleUserId) adduser(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_UserId_AddUser

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	return &models.Response_UserId_AddUser{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}

func (m *ModuleUserId) addcard(elm internal_models.MethodXmlElement) (any, error) {
	var request models.Request_UserId_AddCard

	if err := utils.UnserializeEtreeElement(elm.Element, &request); err != nil {
		panic(err)
	}

	return &models.Response_UserId_AddCard{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,
	}, nil
}
