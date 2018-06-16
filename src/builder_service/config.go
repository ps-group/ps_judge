package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

const (
	configName = "builder_service.json"
)

// Config - video server instance configuration
type Config struct {
	MySQLUser     string `json:"mysql_user"`
	MySQLPassword string `json:"mysql_password"`
	MySQLHost     string `json:"mysql_host"`
	MySQLDB       string `json:"mysql_db"`
	ServerURL     string `json:"server_url"`
}

// ParseConfig loads instance configuration from pre-defined path (relative to executable)
func ParseConfig() (*Config, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}

	configPath := path.Join(path.Dir(executable), configName)
	content, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func NewMySQLConnector(config *Config) DatabaseConnector {
	var connector MySQLConnector
	connector.User = config.MySQLUser
	connector.Password = config.MySQLPassword
	connector.Host = config.MySQLHost
	connector.DatabaseName = config.MySQLDB
	return &connector
}
