package main

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
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
	params := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", c.User, c.Password, c.Host, c.DatabaseName)
	db, err := sql.Open("mysql", params)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect database")
	}
	return db, nil
}
