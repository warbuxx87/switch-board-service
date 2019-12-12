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

	DatabaseHost string
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string

	JWTSecret string
	// 0 = Info, 1 = Debug
	LogLevel int

	HealthCheckFrequency int
}

type ActionResult struct {
}

func LoadConfigurationFile(filename string) (*Configuration, error) {

	configPath, configPathError := resolveConfigPath(filename)

	if configPathError != nil {
		return nil, configPathError
	}

	_, configPathCheckError := os.Stat(configPath)

	if configPathCheckError != nil {
		if err := generateDefaultConfig(filename); err != nil {
			return nil, err
		}
	}

	log.Printf("Loading configuration from %s", configPath)
	configuration := Configuration{}

	data, dataReadError := ioutil.ReadFile(configPath)

	if dataReadError != nil {
		return nil, dataReadError
	}

	jsonUnMarshalError := json.Unmarshal(data, &configuration);
	if jsonUnMarshalError != nil {
		return nil, jsonUnMarshalError
	}

	return &configuration, nil
}

func HashPassword(password string) (string, error) {
	generatedHash, generateHashError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if generateHashError != nil {
		return "", generateHashError
	}

	return string(generatedHash), nil
}

// verifyPassword will return nil if inputPassword is same as hashedPassword
func VerifyPassword(inputPassword, hashedPassword string) error {
	compareHashError := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword));
	if compareHashError != nil {
		return compareHashError
	}

	return nil
}

func resolveConfigPath(filename string) (string, error) {
	currentWorkingDirectory, currentWorkingDirectoryError := os.Getwd()
	if currentWorkingDirectoryError != nil {
		return "", currentWorkingDirectoryError
	}
	return filepath.Join(currentWorkingDirectory, "configs", filename), nil
}

func generateDefaultConfig(filename string) error {

	configPath, configPathError := resolveConfigPath(filename)

	if configPathError != nil {
		return configPathError
	}

	cfg := defaultConfig()
	data, jsonMarshalError := json.MarshalIndent(cfg, "", "\t")

	if jsonMarshalError != nil {
		return jsonMarshalError
	}

	return ioutil.WriteFile(configPath, data, 0664)
}

func defaultConfig() *Configuration {
	return &Configuration{
		ListenIp:             "127.0.0.1",
		ListenPort:           7020,
		SentryUrl:            "",
		DatabaseName:         "db1",
		DatabaseUsername:     "ubuntu",
		DatabasePassword:     "pwd",
		LogLevel:             1,
		JWTSecret:            "somesecret",
		HealthCheckFrequency: 30,
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
