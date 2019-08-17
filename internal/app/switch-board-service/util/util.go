package util

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
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

	configuration := Configuration{}
	err := gonfig.GetConf("config.json", &configuration)

	if err != nil {
		logrus.WithFields(logrus.Fields{}).Error(err)
		return nil, errors.New("Failed to load configuration file")
	}

	return &configuration, nil
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
