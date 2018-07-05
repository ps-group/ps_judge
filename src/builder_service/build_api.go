package main

import (
	"github.com/pkg/errors"
)

// BuildStatusResponse - contains build status.
type BuildStatusResponse struct {
	UUID   string `json:"uuid"`
	Status Status `json:"status"`
}

// BuildReportResponse - contains full build information.
type BuildReportResponse struct {
	UUID        string `json:"uuid"`
	Status      Status `json:"status"`
	Exception   string `json:"exception"`
	BuildLog    string `json:"build_log"`
	TestsPassed int64  `json:"tests_passed"`
	TestsTotal  int64  `json:"tests_total"`
}

// RegisterBuildRequest - contains information required to register new build
// Language - either "c++" or "pascal"
type RegisterBuildRequest struct {
	UUID           string   `json:"uuid"`
	AssignmentUUID string   `json:"assignment_uuid"`
	Language       language `json:"language"`
	Source         string   `json:"source"`
	WebHookURL     string   `json:"web_hook_url"`
}

// RegisterTestCaseRequest - contains information required to register tes case
type RegisterTestCaseRequest struct {
	UUID           string `json:"uuid"`
	AssignmentUUID string `json:"assignment_uuid"`
	Input          string `json:"input"`
	Expected       string `json:"expected"`
}

// RegisterResponse - contains UUID of registered object.
type RegisterResponse struct {
	UUID string `json:"uuid"`
}

func getBuildReport(c APIContext) error {
	key := c.Vars()["uuid"]
	if len(key) == 0 {
		return errors.New("missed 'uuid' request parameter")
	}

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	repo := NewRepository(db)
	report, err := repo.GetBuildReport(key)
	if err != nil {
		return err
	}

	res := &BuildReportResponse{
		UUID:        key,
		Status:      report.Status,
		Exception:   report.Exception,
		BuildLog:    report.BuildLog,
		TestsPassed: report.TestsPassed,
		TestsTotal:  report.TestsTotal,
	}
	return c.WriteJSON(res)
}

func getBuildStatus(c APIContext) error {
	key := c.Vars()["uuid"]
	if len(key) == 0 {
		return errors.New("missed 'uuid' request parameter")
	}

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()
	repo := NewRepository(db)
	status, err := repo.GetBuildStatus(key)
	if err != nil {
		return err
	}

	res := &BuildStatusResponse{
		UUID:   key,
		Status: status,
	}
	return c.WriteJSON(res)
}

func createBuild(c APIContext) error {
	var params RegisterBuildRequest
	err := c.ReadJSON(&params)
	if err != nil {
		return err
	}

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := NewRepository(db)
	assignmentID, err := repo.GetAssignmentID(params.AssignmentUUID)
	if err != nil {
		return err
	}

	err = repo.RegisterBuild(RegisterBuildParams{
		AssignmentID: assignmentID,
		Key:          params.UUID,
		Language:     params.Language,
		Source:       params.Source,
		WebHookURL:   params.WebHookURL,
	})
	if err != nil {
		return err
	}

	res := RegisterResponse{
		UUID: params.UUID,
	}
	return c.WriteJSON(res)
}

func createTestCase(c APIContext) error {
	var params RegisterTestCaseRequest
	err := c.ReadJSON(&params)
	if err != nil {
		return err
	}

	db, err := c.ConnectDB()
	if err != nil {
		return err
	}
	defer db.Close()

	repo := NewRepository(db)
	assignmentID, err := repo.GetAssignmentID(params.AssignmentUUID)
	if err != nil {
		return err
	}

	err = repo.RegisterTestCase(RegisterTestCaseParams{
		AssignmentID: assignmentID,
		Key:          params.UUID,
		Input:        params.Input,
		Expected:     params.Expected,
	})
	if err != nil {
		return err
	}

	res := RegisterResponse{
		UUID: params.UUID,
	}
	return c.WriteJSON(res)
}
