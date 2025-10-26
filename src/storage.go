package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getConfigUsernames() (files []string, err error) {
	fileInfo, err := os.ReadDir("/data")
	if err != nil {
		return
	}

	for _, file := range fileInfo {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			files = append(files, strings.TrimSuffix(file.Name(), ".json"))
		}
	}
	return
}

func readConfig(username string, config *Config) error {
	configName := filepath.Base(username) + ".json"
	configPath := path.Join("/data", configName)

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		*config = Config{Rules: []FilterRule{}}
		return nil
	}
	byteValue, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return err
	}
	return nil
}

func writeConfig(username string, config *Config) error {
	if filepath.Base(username) == "." {
		return fmt.Errorf("empty username not allowed")
	}
	configPath := path.Join("/data", filepath.Base(username)+".json")
	configJson, err := json.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, configJson, 0644)
	if err != nil {
		return err
	}
	return nil
}
