package main

import (
	"github.com/pkg/errors"
)

// BuildInfoResponse - contains full build information.
type BuildInfoResponse struct {
	UUID    string `json:"uuid"`
	Status  Status `json:"status"`
	Score   int64  `json:"score"`
	Details string `json:"details"`
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
	Output         string `json:"output"`
	Expected       string `json:"expected"`
}

// RegisterResponse - contains UUID of registered object.
type RegisterResponse struct {
	UUID string `json:"uuid"`
}

func getBuildInfo(c APIContext) error {
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
	row, err := repo.GetBuildInfo(key)
	if err != nil {
		return err
	}

	score := int64(0)
	row.Score.Scan(&score)
	details := ""
	row.Report.Scan(&details)
	res := &BuildInfoResponse{
		UUID:    key,
		Status:  row.Status,
		Score:   score,
		Details: details,
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
		Output:       params.Output,
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
