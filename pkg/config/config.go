package config

import (
	"encoding/json"
	"os"

	"github.com/ionutinit/riddles-api/pkg/logger"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Database struct {
		Host            string `json:"host"`
		Port            int    `json:"port"`
		User            string `json:"user"`
		Password        string `json:"password"`
		Dbname          string `json:"dbname"`
		Sslmode         string `json:"sslmode"`
		MaxOpenConns    int    `json:"maxOpenConns"`
		MaxIdleConns    int    `json:"maxIdleConns"`
		MaxConnLifetime int    `json:"maxConnLifetime"`
	} `json:"database"`
	ServerPort string   `json:"serverPort"`
	BaseURL    string   `json:"baseUrl"`
	AllowedIPs []string `json:"allowedIPs"`
}

var AppConfig Config

func LoadConfig(configPath string) {
	configFile, err := os.Open(configPath)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":      err,
			"configPath": configPath,
		}).Fatal("Error opening config file")
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{
			"error":      err,
			"configPath": configPath,
		}).Fatal("Error decoding config file")
	}
}
