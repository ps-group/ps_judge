package main

import (
	"database/sql"
	"fmt"
)

type DatabaseConnector interface {
	Connect() (*sql.DB, error)
}

type MySQLConnector struct {
	User         string
	Password     string
	Host         string
	DatabaseName string
}

func (c *MySQLConnector) Connect() (*sql.DB, error) {
	params := fmt.Sprintf("%s:%s@%s/%s", c.User, c.Password, c.Host, c.DatabaseName)
	return sql.Open("mysql", params)
}
