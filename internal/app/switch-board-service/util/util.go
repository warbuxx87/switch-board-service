package util

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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

func LoadConfigurationFile(filename string) (*Configuration, error) {
	configPath, err := resolveConfigPath(filename)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(configPath); err != nil {
	if err := generateDefaultConfig(filename); err != nil {
		return nil, err
		}
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


func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
// verifyPassword will return nil if inputPassword is same as hashedPassword
func VerifyPassword(inputPassword, hashedPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		return err
	}
	return nil
}


func resolveConfigPath(filename string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, "configs", filename), nil
}

func generateDefaultConfig(filename string) error {
	configPath, err := resolveConfigPath(filename)
	if err != nil {
		return err
	}
	cfg := defaultConfig()
	data, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configPath, data, 0664)
}

func defaultConfig() *Configuration {
	return &Configuration{
		ListenIp:		"127.0.0.1",
		ListenPort:		7020,
		SentryUrl:		"",
		DatabaseName:		"db1",
		DatabaseUsername:	"ubuntu",
		DatabasePassword:	"pwd",
		LogLevel:		1,
		JWTSecret:		"somesecret",
	}
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
