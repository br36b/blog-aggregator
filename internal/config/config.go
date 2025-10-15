package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("Unable to get configuration file path: %v", err)
	}

	file, err := os.Open(configFilePath)
	if err != nil {
		return Config{}, fmt.Errorf("Unable to open configuration file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	var userConfig Config
	if err := decoder.Decode(&userConfig); err != nil {
		return Config{}, fmt.Errorf("Unable to decode configuration file: %v", err)
	}

	return userConfig, nil
}

func (cfg Config) SetUser(user string) error {
	cfg.CurrentUserName = user

	err := write(cfg)
	if err != nil {
		return fmt.Errorf("Unable to write to configuration file: %v", err)
	}

	return nil
}

func getConfigFilePath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get home directory for user: %v", err)
	}

	configPath := userHome + "/" + configFileName
	return configPath, nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("Unable to get configuration file path")
	}

	file, err := os.Create(configFilePath)
	if err != nil {
		return fmt.Errorf("Unable to create/overwrite configuration file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&cfg); err != nil {
		return fmt.Errorf("Unable to write to the configuration file")
	}

	return nil
}
