package modules

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"

	internal_models "eamold/internal/models"
	"eamold/services/gfdm_common/models"
)

type MessageDataProvider interface {
	GetMessages(ctx context.Context) ([]string, error)
}

type ModuleMessage struct {
	name string
	db   MessageDataProvider
}

func NewModuleMessage(db MessageDataProvider) *ModuleMessage {
	return &ModuleMessage{
		name: "message",
		db:   db,
	}
}

func (m *ModuleMessage) Name() string {
	return m.name
}

func (m *ModuleMessage) Url() *string {
	return nil
}

func (m *ModuleMessage) Dispatch(elm internal_models.MethodXmlElement) (any, error) {
	switch elm.Module {
	case "message":
		{
			switch elm.Method {
			case "get":
				return m.get(elm)
			}
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *ModuleMessage) get(elm internal_models.MethodXmlElement) (any, error) {
	allMessages, err := m.db.GetMessages(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("message.get: %v", err)
	}

	for i := range allMessages {
		allMessages[i] = strings.ReplaceAll(allMessages[i], "?", "？")
		allMessages[i] = strings.ReplaceAll(allMessages[i], "!", "！")
		allMessages[i] = strings.ReplaceAll(allMessages[i], "\r\n", "?")
		allMessages[i] = strings.ReplaceAll(allMessages[i], "\n", "?")
		allMessages[i] = strings.ReplaceAll(allMessages[i], "\\n", "?")
		allMessages[i] = strings.TrimRight(allMessages[i], "?")
	}

	return &models.Response_Message_Get{
		XMLName: xml.Name{Local: elm.Module},
		Method:  elm.Method,

		String: strings.Join(allMessages, "!"),
	}, nil
}
