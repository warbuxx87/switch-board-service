package util

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Configuration struct {
	ListenIp   string
	ListenPort int

	SentryUrl string

	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string

	JWTSecret string
	// 0 = Info, 1 = Debug
	LogLevel int
}

type ActionResult struct {
}

func LoadConfigurationFile() (*Configuration, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return nil, err
	}
	log.Printf("Loading configuration from %s", configPath)
	configuration := Configuration{}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &configuration); err != nil {
		return nil, err
	}
	return &configuration, nil
}

func resolveConfigPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "configs", "config.json"), nil
}

func ShowBanner() {

	fmt.Printf(`
    ___    ____                  _             _____   ____   _____
   /   |  / / /_____  ____ ___  (_)           / ___/  / __ ) / ___/
  / /| | / / //_/ _ \/ __  __ \/ /  ______    \__ \  / __  | \__ \
 / ___ |/ /  < /  __/ / / / / / /  /_____/   ___/ / / /_/ / ___/ /
/_/  |_/_/_/|_|\___/_/ /_/ /_/_/            /____(_)_____(_)____/

		`)

	fmt.Printf("- Alkemi Switch Board Service - V 0.0.1")
	fmt.Printf("\n")
	fmt.Printf("\n")
}
