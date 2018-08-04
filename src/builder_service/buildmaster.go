package main

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// BuildMaster - reads build tasks from database and passes them to the workers.
type BuildMaster struct {
	reports              chan BuildReport
	stopWorkers          chan struct{}
	stopListening        chan struct{}
	workersWaitGroup     *sync.WaitGroup
	generator            TaskGenerator
	dbConnector          DatabaseConnector
	messageRouterFactory MessageRouterFactory
}

// NewBuildMaster - creates build master with given database
func NewBuildMaster(dbConnector DatabaseConnector, messageRouterFactory MessageRouterFactory) *BuildMaster {
	var master BuildMaster
	master.reports = make(chan BuildReport)
	master.stopWorkers = make(chan struct{})
	master.stopListening = make(chan struct{})
	master.generator = newBuildTaskGenerator(dbConnector, master.reports)
	master.dbConnector = dbConnector
	master.messageRouterFactory = messageRouterFactory

	return &master
}

// RunWorkerPool - runs workers that accept tasks
func (master *BuildMaster) RunWorkerPool() {
	master.workersWaitGroup = RunWorkerPool(master.generator, master.stopWorkers)
	go master.listenBuildReports()
}

// Shutdown - stops all workers and closes channels
func (master *BuildMaster) Shutdown() {
	master.stopWorkers <- struct{}{}
	master.stopListening <- struct{}{}
	<-master.stopListening
	master.workersWaitGroup.Wait()
	close(master.reports)
	close(master.stopWorkers)
	close(master.stopListening)
}

func (master *BuildMaster) listenBuildReports() {
	for {
		select {
		case report := <-master.reports:
			err := master.processBuildReport(report)
			if err != nil {
				logrus.Errorf("cannot process build report: %v", err)
			}
		case <-master.stopListening:
			master.stopListening <- struct{}{}
			return
		}
	}
}

func (master *BuildMaster) processBuildReport(report BuildReport) error {
	db, err := master.dbConnector.Connect()
	if err != nil {
		return errors.Wrap(err, "database connect failed")
	}
	defer db.Close()

	repo := NewRepository(db)
	err = repo.AddBuildReport(report)
	if err != nil {
		return errors.Wrap(err, "cannot add build report")
	}
	err = master.fireBuildFinished(report.Key, report.Status == StatusSucceed)
	if err != nil {
		return errors.Wrap(err, "cannot post build finished")
	}

	return nil
}

func (master *BuildMaster) fireBuildFinished(key string, succeed bool) error {
	router := master.messageRouterFactory.NewMessageRouter()
	router.DeclareExchange(ExchangeBuildFinished)
	router.PublishJSON(ExchangeBuildFinished, BuildFinishedEvent{
		Key:     key,
		Succeed: succeed,
	})
	return router.Error()
}
