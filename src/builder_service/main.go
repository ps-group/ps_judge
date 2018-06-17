package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := ParseConfig()
	if err != nil {
		panic(err)
	}

	databaseConnector := NewMySQLConnector(config)
	webhookService := newWebhookService()
	killChan := getKillSignalChan()

	logFile := openFileLogger(config.LogFileName)
	defer logFile.Close()
	server := startServer(databaseConnector, config.ServerURL)
	waitForKillSignal(killChan)
	webhookService.close()
	server.Shutdown(context.Background())
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

func startServer(connector DatabaseConnector, serverURL string) *http.Server {
	logrus.WithFields(logrus.Fields{"url": serverURL}).Info("starting server")
	server := &http.Server{Addr: serverURL, Handler: newRouter(connector)}
	go func() {
		logrus.Fatal(server.ListenAndServe())
	}()

	return server
}
