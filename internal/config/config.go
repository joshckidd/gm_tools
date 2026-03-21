// internal package used for handling the config file for the cli

package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// the config file should be a json file containing the api url and the current user token
const configFileName = ".gm-tools.json"

type CliConfig struct {
	APIUrl           string `json:"api_url"`
	CurrentUserToken string `json:"current_user_token"`
}

// read the config file
func Read() (CliConfig, error) {
	var config CliConfig

	configFile, err := getConfigFilePath()
	if err != nil {
		return CliConfig{}, err
	}

	configJson, err := os.ReadFile(configFile)
	if err != nil {
		return CliConfig{}, err
	}

	err = json.Unmarshal(configJson, &config)
	if err != nil {
		return CliConfig{}, err
	}

	return config, nil
}

// update the token in the config file
func (config CliConfig) SetToken() error {
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

// get the path of the config file
func getConfigFilePath() (string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configFile := filepath.Join(homePath, configFileName)

	return configFile, nil
}
