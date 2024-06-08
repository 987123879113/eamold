package config

import (
	"fmt"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`

		ServicesUrl string `yaml:"services_url"`

		Logging bool `yaml:"logging"`
	} `yaml:"server"`

	Database struct {
		Driver     string `yaml:"driver"`
		DataSource string `yaml:"datasource"`
	} `yaml:"database"`

	EemallShopServer struct {
		ServerAddress string `yaml:"server_addr"`
	} `yaml:"eemall_shopserver"`

	Static struct {
		ExposeFolders bool `yaml:"expose_folders"`

		Folders []struct {
			StaticPath string  `yaml:"path"`
			DataPath   string  `yaml:"datapath"`
			FileList   *string `yaml:"filelist"`
		} `yaml:"folder"`
	} `yaml:"static"`
}

func (c *Config) MakeUrl(path string) string {
	u, err := url.JoinPath(fmt.Sprintf("http://%s", c.Server.ServicesUrl), path)

	if err != nil {
		return path
	}

	return u
}

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var c Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
