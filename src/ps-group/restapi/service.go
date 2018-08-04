package restapi

import (
	"context"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func openFileLogger(filename string) *os.File {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	logrus.SetOutput(file)
	return file
}

// Service - service interface
type Service interface {
	Start()
	Shutdown()
}

type serviceImpl struct {
	router  *mux.Router
	server  *http.Server
	logFile *os.File
	started bool
}

// NewService - creates new service with given configuration
func NewService(config ServiceConfig) Service {
	var s serviceImpl
	s.router = newRouter(config.RouterConfig, config.Context)
	s.server = &http.Server{Addr: config.ServerURL, Handler: s.router}

	s.logFile = openFileLogger(config.LogFileName)
	logrus.WithFields(logrus.Fields{"url": config.ServerURL}).Info("starting server")

	return &s
}

func (s *serviceImpl) Start() {
	go func() {
		logrus.Fatal(s.server.ListenAndServe())
	}()
	s.started = true
}

func (s *serviceImpl) Shutdown() {
	if s.logFile != nil {
		s.logFile.Close()
		s.logFile = nil
	}
	if s.started {
		s.server.Shutdown(context.Background())
		s.started = false
	}
}
