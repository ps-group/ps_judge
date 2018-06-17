package main

import (
	"github.com/pkg/errors"
)

// BuildInfo - contains full build information.
type BuildInfo struct {
	UUID    string `json:"uuid"`
	Status  Status `json:"status"`
	Score   int    `json:"score"`
	Details string `json:"details"`
}

// RegisterBuild - contains information required to register new build
// Language - either "c++" or "pascal"
type RegisterBuild struct {
	UUID           string   `json:"uuid"`
	AssignmentUUID string   `json:"assignment_uuid"`
	Language       language `json:"language"`
	Source         string   `json:"source"`
	WebHookURL     string   `json:"web_hook_url"`
}

// RegisterTestCase - contains information required to register tes case
type RegisterTestCase struct {
	UUID           string `json:"uuid"`
	AssignmentUUID string `json:"assignment_uuid"`
	Input          string
	Output         string
	Expected       string
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

	info := &BuildInfo{
		UUID:    key,
		Status:  row.Status,
		Score:   row.Score,
		Details: row.Report,
	}
	return c.WriteJSON(info)
}

func createBuild(c APIContext) error {
	var params RegisterBuild
	err := c.ReadJSON(params)
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

	return err
}

func createTestCase(c APIContext) error {
	var params RegisterTestCase
	err := c.ReadJSON(params)
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

	return err
}
