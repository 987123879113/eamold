package services_manager

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/beevik/etree"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"eamold/internal/config"
	"eamold/internal/models"
	"eamold/utils/arc4"
	"eamold/utils/lzss"
)

// All modules must implement this
type Module interface {
	Name() string
	Dispatch(elm models.MethodXmlElement) (any, error)
	Url() *string
}

type ModuleMap map[string]Module

type ServicesManager interface {
	Handler(c echo.Context) error

	ResolveModule(serviceName string, moduleName string) Module
	RegisterService(name string, service Service)

	// Special functions for services module
	GetRegisteredServiceModules(model models.ModelString) map[string]Module
	Url(name string, module Module) string
}

type Service interface {
	AcceptsRequest(model models.ModelString) bool
	GetModuleMap() ModuleMap
}

// Actual services manager code
type servicesManager struct {
	serviceMap          map[string]Service
	moduleMapsByService map[string]ModuleMap
	serviceMapKeyLookup map[string]string
	config              config.Config
}

func (m *servicesManager) ResolveModule(serviceName string, moduleName string) Module {
	if v, ok := m.moduleMapsByService[serviceName]; ok {
		if v2, ok2 := v[moduleName]; ok2 {
			return v2
		}
	}

	return nil
}

func (m *servicesManager) RegisterService(name string, service Service) {
	m.moduleMapsByService[name] = service.GetModuleMap()
	m.serviceMap[name] = service
}

func (m *servicesManager) GetRegisteredServiceModules(model models.ModelString) map[string]Module {
	if k, ok := m.serviceMapKeyLookup[string(model)]; ok {
		return m.moduleMapsByService[k]
	}

	if modules, ok := m.getAcceptingService(model); ok {
		return modules
	}

	return nil
}

func (m *servicesManager) UrlPath(name string, module Module) string {
	urlOverride := module.Url()

	if urlOverride != nil {
		return *urlOverride
	}

	return strings.Replace("/"+strings.Join([]string{name, module.Name()}, "/"), "//", "/", -1)
}

func (m *servicesManager) Url(name string, module Module) string {
	urlOverride := module.Url()

	if urlOverride != nil {
		return *urlOverride
	}

	return m.config.MakeUrl("/")
}

type MethodLookup map[string]func(models.RequestCall) (any, error)

func (m *servicesManager) getAcceptingService(model models.ModelString) (ModuleMap, bool) {
	if k, ok := m.serviceMapKeyLookup[string(model)]; ok {
		return m.moduleMapsByService[k], ok
	}

	for serviceKey, service := range m.serviceMap {
		if service.AcceptsRequest(model) {
			// Cache service that accepts given model
			m.serviceMapKeyLookup[string(model)] = serviceKey
			return m.moduleMapsByService[serviceKey], true
		}
	}

	return nil, false
}

func (m *servicesManager) Dispatch(elm models.MethodXmlElement) (any, error) {
	if modules, ok := m.getAcceptingService(elm.Model); ok {
		if callback, ok2 := modules[elm.Module]; ok2 {
			return callback.Dispatch(elm)
		}

		// Try all dispatch functions for the case where the module name provided in the XML doesn't match the service name
		for _, callback := range modules {
			if callback == nil {
				continue
			}

			r, err := callback.Dispatch(elm)
			if err == nil {
				// found match
				return r, nil
			}

			log.Printf("err: %v\n", err)
		}
	}

	return nil, fmt.Errorf("unknown call %s %s %s", elm.Model, elm.Module, elm.Method)
}

func (m *servicesManager) Handler(c echo.Context) error {
	req := c.Request()
	body, _ := io.ReadAll(req.Body)

	// Decrypt body if needed
	if v, ok := req.Header["X-Eamuse-Info"]; ok && len(v) == 1 && len(v[0]) == 15 && v[0][1] == '-' && v[0][10] == '-' {
		// The first byte of the value seems to correspond with encryption method?
		// But I've only ever seen 1 here
		key := v[0][2:10] + v[0][11:15]
		arc4.XORKeyStream(key, body)
	}

	// Decompress body if needed
	if v, ok := req.Header["X-Compress"]; ok {
		if len(v) > 0 && v[0] == "lz77" {
			decompressed, err := lzss.Decompress(body)
			if err != nil {
				panic(err)
			}

			body = decompressed
		} else {
			panic(fmt.Sprintf("unknown compression: %s", v[0]))
		}
	}

	messageEncodings := []encoding.Encoding{
		japanese.EUCJP,
		japanese.ShiftJIS,
		unicode.UTF8,
	}

	var validEncoding encoding.Encoding
	for _, messageEncoding := range messageEncodings {
		bodyString, _, err := transform.String(messageEncoding.NewDecoder(), string(body))
		if err == nil {
			_, _, err := transform.String(messageEncoding.NewEncoder(), bodyString)
			if err == nil {
				validEncoding = messageEncoding
				break
			}
		}
	}

	if validEncoding == nil {
		panic("couldn't find valid encoding for message")
	}

	bodyString, _, err := transform.String(validEncoding.NewDecoder(), string(body))
	if err != nil {
		panic(err)
	}

	log.Print(bodyString)
	log.Printf("validEncoding: %v\n", validEncoding)

	doc := etree.NewDocument()
	if err := doc.ReadFromString(bodyString); err != nil {
		panic(err)
	}

	curTime := time.Now()
	outputPath := ""

	responses := []any{}
	if call := doc.SelectElement("call"); call != nil {
		srcid := call.SelectAttrValue("srcid", "")
		model := call.SelectAttrValue("model", "")

		if srcid == "" {
			log.Error(fmt.Errorf("invalid srcid: %s", srcid))
		}

		if model == "" {
			log.Error(fmt.Errorf("invalid model: %s", model))
		}

		modelmodels := models.ModelString(model)

		if m.config.Server.Logging {
			outputPath = fmt.Sprintf("output/%s_%s_%s", modelmodels.Model(), modelmodels.Dest(), modelmodels.Spec())
			os.MkdirAll(outputPath, os.ModePerm)
			os.WriteFile(fmt.Sprintf("%s/%d_request.xml", outputPath, curTime.Unix()), []byte(bodyString), os.ModePerm)
		}

		for _, methodCall := range call.ChildElements() {
			elm := models.MethodXmlElement{
				Model:    modelmodels,
				SourceId: srcid,
				Module:   methodCall.Tag,
				Method:   methodCall.SelectAttrValue("method", ""),
				Element:  methodCall,
			}

			resp, err := m.Dispatch(elm)

			if err != nil {
				log.Error(string(body))
				log.Error(err)

				// TODO: What's the proper way to respond with an error?
				return c.XML(http.StatusInternalServerError, models.Response{
					Fault:  -1,
					Status: -1,
				})
			} else {
				responses = append(responses, resp)
			}
		}
	}

	// Generate response struct based on the responses of the individual module calls
	respmodels := models.Response{
		Body: responses,
	}

	resp, err := xml.Marshal(respmodels)
	if err != nil {
		log.Error(string(body))
		log.Error(err)

		return c.XML(http.StatusInternalServerError, models.Response{
			Fault:  -1,
			Status: -1,
		})
	}

	// Compress response if needed
	/*
		if v, ok := req.Header["X-Compress"]; ok {

			if len(v) > 0 && v[0] == "lz77" {
				compressed, err := lzss.Compress(resp)
				if err != nil {
					panic(err)
				}

				resp = compressed

				c.Response().Header().Set("X-Compress", v[0])
			} else {
				panic(fmt.Sprintf("unknown compression: %s", v[0]))
			}
		}
	*/

	log.Printf("resp: %s\n", string(resp))
	log.Printf("messageEncoding: %v\n", validEncoding)
	respString, _, err := transform.String(validEncoding.NewEncoder(), string(resp))
	if err != nil {
		panic(err)
	}

	resp = []byte(respString)

	if m.config.Server.Logging {
		os.WriteFile(fmt.Sprintf("%s/%d_response.xml", outputPath, curTime.Unix()), resp, os.ModePerm)
	}

	// Encrypt response if needed
	if v, ok := req.Header["X-Eamuse-Info"]; ok && len(v) == 1 && len(v[0]) == 15 && v[0][1] == '-' && v[0][10] == '-' {
		// The first byte of the value seems to correspond with encryption method?
		// But I've only ever seen 1 here
		key := v[0][2:10] + v[0][11:15]
		arc4.XORKeyStream(key, resp)

		c.Response().Header().Set("X-Eamuse-Info", v[0])
	}

	return c.Blob(http.StatusOK, "application/octet-stream", resp)
}

func NewServicesManager(config config.Config) ServicesManager {
	return &servicesManager{
		serviceMap:          map[string]Service{},
		moduleMapsByService: map[string]ModuleMap{},
		serviceMapKeyLookup: map[string]string{},
		config:              config,
	}
}
