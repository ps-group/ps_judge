package main

import (
	"ps-group/restapi"
)

// APIPrefix - prefix for each builder API method, contains API version
const (
	APIPrefix     = "/api/v1/"
	DefaultScheme = "http://"
)

// BuilderService - accessor to the builder service REST API
type BuilderService interface {
	RegisterNewBuild(buildUUID string, assignmentUUID string, language string, source string) (*RegisterResponse, error)
	RegisterTestCase(testUUID string, assignmentUUID string, input string, expected string) (*RegisterResponse, error)
	GetBuildReport(buildUUID string) (*BuildReportResponse, error)
}

type builderServiceImpl struct {
	client *restapi.Client
}

// RegisterResponse - contains UUID of registered object.
type RegisterResponse struct {
	UUID string `json:"uuid"`
}

// BuildReportResponse - contains detailed report about finished build
type BuildReportResponse struct {
	UUID        string `json:"uuid"`
	Status      string `json:"status"`
	Exception   string `json:"exception"`
	BuildLog    string `json:"build_log"`
	TestsLog    string `json:"tests_log"`
	TestsPassed int64  `json:"tests_passed"`
	TestsTotal  int64  `json:"tests_total"`
}

// NewBuilderService - creates new builder service accessor
func NewBuilderService(builderURL string) BuilderService {
	bs := new(builderServiceImpl)
	bs.client = restapi.NewClient(DefaultScheme + builderURL + APIPrefix)
	return bs
}

// RegisterNewBuild - registers new solution build
func (bs *builderServiceImpl) RegisterNewBuild(buildUUID string, assignmentUUID string, language string, source string) (*RegisterResponse, error) {
	params := map[string]string{
		"uuid":            buildUUID,
		"assignment_uuid": assignmentUUID,
		"language":        language,
		"source":          source,
	}
	var result RegisterResponse
	err := bs.client.Post("build/new", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// RegisterTestCase - registers new test case for assignment solutions.
func (bs *builderServiceImpl) RegisterTestCase(testUUID string, assignmentUUID string, input string, expected string) (*RegisterResponse, error) {
	params := map[string]string{
		"uuid":            testUUID,
		"assignment_uuid": assignmentUUID,
		"input":           input,
		"expected":        expected,
	}
	var result RegisterResponse
	err := bs.client.Post("testcase/new", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetBuildReport - queries report for the finished build
func (bs *builderServiceImpl) GetBuildReport(buildUUID string) (*BuildReportResponse, error) {
	var result BuildReportResponse
	err := bs.client.Get("build/report/"+buildUUID, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
