package main

import (
	"database/sql"
)

type apiContext struct {
	dbConnector    DatabaseConnector
	db             *sql.DB
	builderService BuilderService
}

func newAPIContext(dbConnector DatabaseConnector, builderService BuilderService) *apiContext {
	c := new(apiContext)
	c.dbConnector = dbConnector
	c.builderService = builderService
	return c
}

func (c *apiContext) BuilderAPI() BuilderService {
	return c.builderService
}

func (c *apiContext) ConnectDB() (*BackendRepository, error) {
	db, err := c.dbConnector.Connect()
	if err != nil {
		return nil, err
	}
	c.db = db
	return NewBackendRepository(db), nil
}

func (c *apiContext) Close() {
	if c.db != nil {
		c.db.Close()
		c.db = nil
	}
}