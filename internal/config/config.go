package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gm-tools.json"

type Config struct {
	APIUrl           string `json:"api_url"`
	CurrentUserToken string `json:"current_user_token"`
}

func Read() (Config, error) {
	var config Config

	configFile, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	configJson, err := os.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(configJson, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config Config) SetToken() error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(configFile, configJson, 0600)
	return err
}

func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(homePath, configFileName)

	return configFile, nil
}
