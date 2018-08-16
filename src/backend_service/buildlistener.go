package main

import (
	"ps-group/judgeevents"

	"github.com/sirupsen/logrus"
)

const (
	BackendBuildEventsQueue = "BackendBuildEvents"
	MaxPercentage           = 100
)

type buildListener struct {
	events    judgeevents.BuilderEvents
	builder   BuilderService
	connector DatabaseConnector
}

func newBuildListener(connector DatabaseConnector, builder BuilderService, socket string) *buildListener {
	listener := new(buildListener)
	listener.events = judgeevents.NewBuilderEvents(socket)
	listener.builder = builder
	listener.connector = connector
	return listener
}

func (listener *buildListener) Start() error {
	listener.events.ConsumeBuildFinished(BackendBuildEventsQueue, func(event judgeevents.BuildFinishedEvent) {
		err := listener.processBuild(event)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"uuid":  event.Key,
				"error": err,
			}).Error("cannot process build")
		} else {
			logrus.WithFields(logrus.Fields{
				"uuid":   event.Key,
				"status": event.Succeed,
			}).Info("build finished")
		}
	})
	return listener.events.Error()
}

func (listener *buildListener) Close() {
	listener.events.Close()
}

func (listener *buildListener) processBuild(event judgeevents.BuildFinishedEvent) error {
	status := "succeed"
	if !event.Succeed {
		status = "failed"
	}

	var newScore int64
	if event.Succeed {
		report, err := listener.builder.GetBuildReport(event.Key)
		if err != nil {
			return err
		}
		if report.TestsTotal > 0 {
			newScore = MaxPercentage * report.TestsPassed / report.TestsTotal
		}
	}

	db, err := listener.connector.Connect()
	if err != nil {
		return err
	}
	repo := NewBackendRepository(db)

	commit, err := listener.updateCommit(repo, event.Key, status, newScore)
	if err != nil {
		return err
	}
	return listener.updateSolution(repo, commit.SolutionID, commit.BuildScore)
}

func (listener *buildListener) updateCommit(repo *BackendRepository, buildUUID string, status string, newScore int64) (*CommitModel, error) {

	commit, err := repo.getCommitInfoByUUID(buildUUID)
	if err != nil {
		return nil, err
	}

	commit.BuildScore = newScore
	commit.BuildStatus = status
	err = repo.updateCommit(commit)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func (listener *buildListener) updateSolution(repo *BackendRepository, solutionID int64, newScore int64) error {
	solution, err := repo.getSolution(solutionID)
	if err != nil {
		return err
	}

	if solution.Score < newScore {
		err := repo.updateSolutionScore(solutionID, newScore)
		if err != nil {
			return err
		}
	}

	return nil
}
