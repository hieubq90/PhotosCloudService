package app_config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type AppConfiguration struct {
	// Base settings
	ListenHost      string   `yaml:"listten_host"`
	ListenPort      int      `yaml:"listten_port"`
	RuntimeMaxProcs int      `yaml:"maxprocs"`
	RunMode         string   `yaml:"run_mode"`
	HandleDemo      bool     `yaml:"demo"`
	AllowOrigins    []string `yaml:"allow_origins"`
	AllowMethods    []string `yaml:"allow_methods"`
	AllowFileTypes  []string `yaml:"allow_file_types"`
	BodyLimitSize   int      `yaml:"body_limit_size"`
	SaveLocation    string   `yaml:"save_location"`
	DownloadDomain  string   `yaml:"download_domain"`
}

var AppConfig *AppConfiguration

func (c *AppConfiguration) IsAllowedFileType(fileType string) bool {
	if len(c.AllowFileTypes) > 0 {
		for _, allowedType := range c.AllowFileTypes {
			if strings.Compare(allowedType, fileType) == 0 {
				return true
			}
		}
	}
	return false
}

func InitFromYAML() bool {
	filename := os.Args[0] + ".yaml"
	fmt.Println("[PhotosCloudService] Start loading configurations from", filename)
	yamlAbsPath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	// read the raw contents of the file
	data, err := ioutil.ReadFile(yamlAbsPath)
	if err != nil {
		fmt.Println("[PhotosCloudService] Read config file error.")
		panic(err)
	}

	c := AppConfiguration{}
	// put the file's contents as yaml to the default configuration(c)
	if err := yaml.Unmarshal(data, &c); err != nil {
		panic(err)
	}

	AppConfig = &c

	fmt.Println("[PhotosCloudService] Load configurations successful!")
	return true
}
