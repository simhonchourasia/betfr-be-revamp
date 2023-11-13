package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

// Add field here when new config element in json
type Config struct {
	DatabaseURL    string `json:"databaseURL"`
	SecretKey      string `json:"secretKey"`
	Domain         string `json:"domain"`
	Port           string `json:"port"`
	Debug          bool   `json:"debug"`
	MigrationsOnly bool   `json:"migrationsOnly"`
	OriginFE       string `json:"originFE"`
}

var GlobalConfig Config

func SetupConfig() error {
	var cfgErr error = nil
	var configOnce sync.Once
	configOnce.Do(func() {
		log.Println("Reading config file...")
		jsonFile, err := os.Open("config/default.json")
		defer func() {
			cfgErr = jsonFile.Close()
		}()
		if err != nil {
			cfgErr = err
		} else {
			configBytes, err := io.ReadAll(jsonFile)
			if err != nil {
				cfgErr = err
			} else {
				json.Unmarshal(configBytes, &GlobalConfig)
			}
		}
	})
	return cfgErr
}
