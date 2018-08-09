package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/streadway/amqp"
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
	AmqpSocket    string `json:"amqp_socket"`
	ServerURL     string `json:"builder_url"`
	LogFileName   string `json:"log_file_name"`
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

// DatabaseConnector - creates SQL database connection
type DatabaseConnector interface {
	Connect() (*sql.DB, error)
}

type mySQLConnector struct {
	User         string
	Password     string
	Host         string
	DatabaseName string
}

func (c *mySQLConnector) Connect() (*sql.DB, error) {
	params := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", c.User, c.Password, c.Host, c.DatabaseName)
	db, err := sql.Open("mysql", params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect database")
	}
	return db, nil
}

// NewMySQLConnector - creates MySQL database connector
func NewMySQLConnector(config *Config) DatabaseConnector {
	var connector mySQLConnector
	connector.User = config.MySQLUser
	connector.Password = config.MySQLPassword
	connector.Host = config.MySQLHost
	connector.DatabaseName = config.MySQLDB
	return &connector
}

// MessageRouterFactory - creates new message routers
type MessageRouterFactory interface {
	NewMessageRouter() *MessageRouter
}

type messageRouterFactoryImpl struct {
	socket string
}

// NewMessageRouterFactory - creates factory that creates routers on given socket
func NewMessageRouterFactory(socket string) MessageRouterFactory {
	f := new(messageRouterFactoryImpl)
	f.socket = socket
	return f
}

// NewMessageRouter - creates message router
func (f *messageRouterFactoryImpl) NewMessageRouter() *MessageRouter {
	var router MessageRouter
	router.conn, router.lastError = amqp.Dial(f.socket)
	if router.lastError == nil {
		defer func() {
			if router.channel == nil {
				router.conn.Close()
				router.conn = nil
			}
		}()
		router.channel, router.lastError = router.conn.Channel()
	}
	return &router
}
