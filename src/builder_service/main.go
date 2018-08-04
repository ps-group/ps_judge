package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"ps-group/restapi"
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		panic(err)
	}

	databaseConnector := NewMySQLConnector(config)
	messageRouterFactory := NewMessageRouterFactory(config.AmqpSocket)
	context := &apiContext{databaseConnector}

	master := NewBuildMaster(databaseConnector, messageRouterFactory)
	killChan := getKillSignalChan()
	service := restapi.NewService(restapi.ServiceConfig{
		RouterConfig: g_routes,
		ServerURL:    config.ServerURL,
		LogFileName:  config.LogFileName,
		Context:      context,
	})
	defer service.Shutdown()

	// Start services
	service.Start()
	master.RunWorkerPool()

	// Wait for SIGTERM
	waitForKillSignal(killChan)

	master.Shutdown()
}

func getKillSignalChan() chan os.Signal {
	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, os.Kill, os.Interrupt, syscall.SIGTERM)
	return killChan
}

func waitForKillSignal(killChan <-chan os.Signal) {
	killSignal := <-killChan
	switch killSignal {
	case os.Interrupt:
		logrus.Info("got SIGINT, shutting down...")
	case syscall.SIGTERM:
		logrus.Info("got SIGTERM, shutting down...")
	}
}
