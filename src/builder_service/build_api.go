package main

import (
	"database/sql"
	"ps-group/restapi"
)

type apiContext struct {
	dbConnector DatabaseConnector
}

func (c *apiContext) ConnectDB() (*sql.DB, error) {
	return c.dbConnector.Connect()
}

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
	TestsLog    string `json:"tests_log"`
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

func getBuildReport(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)

	key := req.Var("uuid")
	if len(key) == 0 {
		return restapi.NewInternalError("missed 'uuid' request parameter")
	}

	db, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}
	defer db.Close()
	repo := NewRepository(db)
	report, err := repo.GetBuildReport(key)
	if err != nil {
		return &restapi.InternalError{err}
	}

	res := &BuildReportResponse{
		UUID:        key,
		Status:      report.Status,
		Exception:   report.Exception,
		BuildLog:    report.BuildLog,
		TestsLog:    report.TestsLog,
		TestsPassed: report.TestsPassed,
		TestsTotal:  report.TestsTotal,
	}
	return &restapi.Ok{&res}
}

func getBuildStatus(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	key := req.Var("uuid")
	if len(key) == 0 {
		return restapi.NewInternalError("missed 'uuid' request parameter")
	}

	db, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}
	defer db.Close()
	repo := NewRepository(db)
	status, err := repo.GetBuildStatus(key)
	if err != nil {
		return &restapi.InternalError{err}
	}

	res := &BuildStatusResponse{
		UUID:   key,
		Status: status,
	}
	return &restapi.Ok{&res}
}

func createBuild(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)

	var params RegisterBuildRequest
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.InternalError{err}
	}

	db, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}
	defer db.Close()

	repo := NewRepository(db)
	assignmentID, err := repo.GetAssignmentID(params.AssignmentUUID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	err = repo.RegisterBuild(RegisterBuildParams{
		AssignmentID: assignmentID,
		Key:          params.UUID,
		Language:     params.Language,
		Source:       params.Source,
	})
	if err != nil {
		return &restapi.InternalError{err}
	}

	res := RegisterResponse{
		UUID: params.UUID,
	}
	return &restapi.Ok{&res}
}

func createTestCase(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)

	var params RegisterTestCaseRequest
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.InternalError{err}
	}

	db, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}
	defer db.Close()

	repo := NewRepository(db)
	assignmentID, err := repo.GetAssignmentID(params.AssignmentUUID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	err = repo.RegisterTestCase(RegisterTestCaseParams{
		AssignmentID: assignmentID,
		Key:          params.UUID,
		Input:        params.Input,
		Expected:     params.Expected,
	})
	if err != nil {
		return &restapi.InternalError{err}
	}

	res := RegisterResponse{
		UUID: params.UUID,
	}
	return &restapi.Ok{&res}
}
