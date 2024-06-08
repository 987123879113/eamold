package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"time"

	_ "modernc.org/sqlite"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"eamold/internal/config"
	"eamold/internal/services_manager"

	"eamold/services/core"

	"eamold/services/dm10"
	"eamold/services/dm7"
	"eamold/services/dm7puv"
	"eamold/services/dm8"
	"eamold/services/dm9"

	"eamold/services/gf10"
	"eamold/services/gf11"
	"eamold/services/gf8"
	"eamold/services/gf8puv"
	"eamold/services/gf9"
)

func main() {
	config, err := config.NewConfig("config.yml")
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic("must provide valid configuration file")
	}

	db, err := sql.Open(config.Database.Driver, config.Database.DataSource)
	if err != nil {
		panic(err)
	}

	manager := services_manager.NewServicesManager(*config)

	manager.RegisterService(core.SERVICE_NAME, core.New(manager, db))

	manager.RegisterService(gf8.SERVICE_NAME, gf8.New(manager, db))
	manager.RegisterService(gf8puv.SERVICE_NAME, gf8puv.New(manager, db))
	manager.RegisterService(gf9.SERVICE_NAME, gf9.New(manager, db))
	manager.RegisterService(gf10.SERVICE_NAME, gf10.New(manager, db, config.EemallShopServer.ServerAddress))
	manager.RegisterService(gf11.SERVICE_NAME, gf11.New(manager, db, config.EemallShopServer.ServerAddress))

	manager.RegisterService(dm7.SERVICE_NAME, dm7.New(manager, db))
	manager.RegisterService(dm7puv.SERVICE_NAME, dm7puv.New(manager, db))
	manager.RegisterService(dm8.SERVICE_NAME, dm8.New(manager, db))
	manager.RegisterService(dm9.SERVICE_NAME, dm9.New(manager, db, config.EemallShopServer.ServerAddress))
	manager.RegisterService(dm10.SERVICE_NAME, dm10.New(manager, db, config.EemallShopServer.ServerAddress))

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", manager.Handler)
	e.POST("/services", manager.Handler)

	// Expose any configured static paths
	for _, folder := range config.Static.Folders {
		if folder.FileList != nil {
			jsonFile, err := os.Open(*folder.FileList)
			if err != nil {
				log.Printf("Could not open file list: %v", err)
				continue
			}

			jsonBytes, err := io.ReadAll(jsonFile)
			if err != nil {
				panic(err)
			}

			var staticFileList struct {
				Files []struct {
					OverridePath *string `json:"override_path"`
					FilePath     string  `json:"path"`
				}
			}

			if err := json.Unmarshal(jsonBytes, &staticFileList); err != nil {
				panic(err)
			}

			for _, file := range staticFileList.Files {
				filePath := file.FilePath

				if file.OverridePath != nil {
					filePath = *file.OverridePath
				}

				absPath, err := url.JoinPath(folder.StaticPath, filePath)
				if err != nil {
					panic(err)
				}

				dataPath := path.Join(folder.DataPath, file.FilePath)

				log.Printf("exposing static file: %v -> %v\n", absPath, dataPath)
				e.File(absPath, dataPath)
			}
		} else if config.Static.ExposeFolders {
			log.Printf("exposing static folder: %v -> %v\n", folder.StaticPath, folder.DataPath)
			e.Static(folder.StaticPath, folder.DataPath)
		}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal(fmt.Errorf("shutting down the server: %v", err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
