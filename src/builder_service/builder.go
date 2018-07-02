package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	MaxBuildScore = 100
)

type BuildReport struct {
	log    string
	score  int
	status Status
}

type buildTask struct {
	language language
	source   string
	key      string
	cases    []testCase
	output   chan BuildReport
}

func createBuildReport(result BuildResult) BuildReport {
	var report BuildReport
	if result.internalError != nil {
		report.log = result.internalError.Error()
		report.status = StatusException
	} else if result.buildError != nil {
		report.log = result.buildError.Error()
		report.status = StatusFailed
	} else {
		succeedCount := 0
		for i, err := range result.testCaseErrors {
			if err != nil {
				report.log += fmt.Sprintf("#%d case failed\n", i+1)
			} else {
				succeedCount++
			}
		}
		report.status = StatusSucceed
		report.score = (succeedCount * MaxBuildScore) / len(result.testCaseErrors)
	}
	return report
}

func (t *buildTask) Run(workerID int) error {
	workdir := fmt.Sprintf("builder_%d", workerID)
	logrus.WithField("uuid", t.key).Info("running build")
	result := buildSolution(t.source, t.language, t.cases, workdir)
	report := createBuildReport(result)
	t.output <- report
	return nil
}

type buildTaskGenerator struct {
	connector DatabaseConnector
}

func newBuildTaskGenerator(connector DatabaseConnector) *buildTaskGenerator {
	var generator buildTaskGenerator
	generator.connector = connector
	return &generator
}

func (g *buildTaskGenerator) Next() (bool, Task) {
	db, err := g.connector.Connect()
	if err != nil {
		logrus.WithField("error", err).Error("database connect failed")
		return false, nil
	}
	defer db.Close()
	repo := NewRepository(db)
	build, err := repo.PullPendingBuild()
	if err != nil {
		logrus.WithField("error", err).Error("read task from database failed")
		return false, nil
	}
	if build == nil {
		return false, nil
	}
	var task buildTask
	task.language = build.Language
	task.source = build.Source
	task.key = build.Key
	task.output = make(chan BuildReport)
	return true, &task
}
