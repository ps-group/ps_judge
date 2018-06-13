package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

const (
	configName = "builder_service.json"
)

// Config - video server instance configuration
type Config struct {
	Workdir       string `json:"workdir"`
	MySQLUser     string `json:"mysql_user"`
	MySQLPassword string `json:"mysql_password"`
	MySQLHost     string `json:"mysql_host"`
	MySQLDB       string `json:"mysql_db"`
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

	_, err = os.Stat(config.Workdir)
	if err != nil && !os.IsExist(err) {
		return nil, errors.New("workdir does not exist: " + config.Workdir)
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
