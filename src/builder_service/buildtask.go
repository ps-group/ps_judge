package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type buildTask struct {
	language language
	source   string
	key      string
	cases    []TestCase
	reports  chan BuildReport
}

func (t *buildTask) Run(workerID int) error {
	workdir := fmt.Sprintf("builder_%d", workerID)
	logrus.WithField("uuid", t.key).Info("running build")
	result := buildSolution(t.source, t.language, t.cases, workdir)
	report := t.createBuildReport(result)
	t.reports <- report
	return nil
}

type buildTaskGenerator struct {
	connector DatabaseConnector
	reports   chan BuildReport
}

func (t *buildTask) createBuildReport(result BuildResult) BuildReport {
	var report BuildReport
	report.Key = t.key
	if result.internalError != nil {
		report.Exception = result.internalError.Error()
		report.Status = StatusException
	} else if result.buildError != nil {
		report.BuildLog = result.buildError.Error()
		report.Status = StatusFailed
	} else {
		report.Status = StatusSucceed
		report.TestsTotal = int64(len(result.testCaseErrors))
		report.TestsPassed = 0
		for i, err := range result.testCaseErrors {
			if err == nil {
				report.TestsPassed++
			} else {
				report.TestsLog += fmt.Sprintf("--- FAILURE IN TEST %d ---\n%s\n", i, err.Error())
			}
		}
	}
	return report
}

func newBuildTaskGenerator(connector DatabaseConnector, reports chan BuildReport) *buildTaskGenerator {
	var generator buildTaskGenerator
	generator.connector = connector
	generator.reports = reports

	return &generator
}

func (g *buildTaskGenerator) Next() (bool, Task) {
	db, err := g.connector.Connect()
	if err != nil {
		logrus.WithField("error", err).Error("database connect failed")
		return false, nil
	}
	defer db.Close()
	repo := NewBuilderRepository(db)
	build, err := repo.PullPendingBuild()
	if err != nil {
		logrus.WithField("error", err).Error("read task from database failed")
		return false, nil
	}
	if build == nil {
		return false, nil
	}
	cases, err := repo.GetTestCases(build.AssignmentID)
	if err != nil {
		logrus.WithField("error", err).Error("cannot read test cases")
		return false, nil
	}
	var task buildTask
	task.language = build.Language
	task.source = build.Source
	task.key = build.Key
	task.cases = cases
	task.reports = g.reports
	return true, &task
}
