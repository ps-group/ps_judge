package main

import (
	"database/sql"
	"ps-group/restapi"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type apiContext struct {
	dbConnector    DatabaseConnector
	db             *sql.DB
	builderService BuilderService
}

type valuesMap map[string]interface{}
type valuesMapList []valuesMap

func newAPIContext(dbConnector DatabaseConnector, builderService BuilderService) *apiContext {
	c := new(apiContext)
	c.dbConnector = dbConnector
	c.builderService = builderService
	return c
}

func (c *apiContext) BuilderAPI() BuilderService {
	return c.builderService
}

func (c *apiContext) ConnectDB() (*BackendRepository, error) {
	db, err := c.dbConnector.Connect()
	if err != nil {
		return nil, err
	}
	c.db = db
	return NewBackendRepository(db), nil
}

func (c *apiContext) Close() {
	if c.db != nil {
		c.db.Close()
		c.db = nil
	}
}

func parseID(req restapi.Request, name string) (int64, error) {
	return strconv.ParseInt(req.Var(name), 10, 64)
}

// LoginUserParams - contains user login parameters
type LoginUserParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

func loginUser(ctx interface{}, req restapi.Request) restapi.Response {
	var params LoginUserParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	info, err := repository.getUserInfoByUsername(params.Username)
	if err != nil {
		return &restapi.InternalError{err}
	}

	if info != nil && info.PasswordHash == params.PasswordHash {
		return &restapi.Ok{valuesMap{
			"succeed": true,
			"user_id": info.ID,
			"user": valuesMap{
				"username": info.Username,
				"roles":    info.Roles,
			},
		}}
	}

	return &restapi.Unauthorized{valuesMap{
		"succeed": false,
	}}
}

func getUserInfo(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	userID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	info, err := repository.getUserInfo(userID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{valuesMap{
		"username": info.Username,
		"roles":    info.Roles,
	}}
}

func getUserContestList(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	userID, err := parseID(req, "user_id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid user_id")}
	}

	contests, err := repository.getUserContestList(userID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	var results valuesMapList
	for _, contest := range contests {
		results = append(results, valuesMap{
			"id":    contest.ID,
			"title": contest.Title,
		})
	}
	return &restapi.Ok{results}
}

func getUserContestSolutions(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	userID, err := parseID(req, "user_id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid user_id")}
	}

	contestID, err := parseID(req, "contest_id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid contest_id")}
	}

	solutions, err := repository.getUserContestSolutions(userID, contestID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	var results valuesMapList
	for _, solution := range solutions {
		commit, err := repository.getLastCommit(solution.ID)
		if err != nil {
			return &restapi.InternalError{err}
		}
		results = append(results, valuesMap{
			"assignment_id":    solution.AssignmentID,
			"assignment_title": solution.AssignmentTitle,
			"score":            solution.Score,
			"commit_id":        commit.ID,
			"build_status":     commit.BuildStatus,
		})
	}
	return &restapi.Ok{results}
}

func getContestResults(ctx interface{}, req restapi.Request) restapi.Response {
	contestID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	results, err := repository.getContestResults(contestID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{results}
}

// CommitSolutionParams - parameters to commit solution
type CommitSolutionParams struct {
	UUID         string `json:"uuid"`
	AssignmentID int64  `json:"assignment_id"`
	Language     string `json:"language"`
	Source       string `json:"source"`
}

func commitSolution(ctx interface{}, req restapi.Request) restapi.Response {
	userID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	var params CommitSolutionParams
	err = req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid JSON")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	solution, err := repository.getUserAssignmentSolution(userID, params.AssignmentID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	if solution == nil {
		solution, err = repository.createSolution(userID, params.AssignmentID)
		if err != nil {
			return &restapi.InternalError{err}
		}
	}

	assignment, err := repository.getAssignment(params.AssignmentID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	err = repository.createCommit(solution.ID, params.UUID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	response, err := c.BuilderAPI().RegisterNewBuild(params.UUID, assignment.UUID, params.Language, params.Source)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{response}
}

func getCommitReport(ctx interface{}, req restapi.Request) restapi.Response {
	commitID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	commitUUID, err := repository.getCommitUUID(commitID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	response, err := c.BuilderAPI().GetBuildReport(commitUUID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{response}
}

func getContestAssignments(ctx interface{}, req restapi.Request) restapi.Response {
	contestID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	assignments, err := repository.getContestAssignments(contestID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	var infos valuesMapList
	for _, assignment := range assignments {
		infos = append(infos, valuesMap{
			"id":         assignment.ID,
			"contest_id": assignment.ContestID,
			"uuid":       assignment.UUID,
			"title":      assignment.Title,
		})
	}
	return &restapi.Ok{infos}
}

func getAssignmentInfo(ctx interface{}, req restapi.Request) restapi.Response {
	assignmentID, err := parseID(req, "id")
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	assignment, err := repository.getAssignmentFull(assignmentID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	result := valuesMap{
		"id":          assignmentID,
		"contest_id":  assignment.ContestID,
		"uuid":        assignment.UUID,
		"title":       assignment.Title,
		"description": assignment.Description,
	}
	return &restapi.Ok{result}
}

// CreateContestParams - parameters for the new contest
type CreateContestParams struct {
	Title      string `json:"title"`
	MaxReviews uint   `json:"max_reviews"`
}

func createContest(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateContestParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	model := ContestModel{
		Title:      params.Title,
		MaxReviews: params.MaxReviews,
	}
	err = repository.createContest(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}
	return &restapi.Ok{&valuesMap{
		"id": model.ID,
	}}
}

// CreateUserParams - parameters for the new user
type CreateUserParams struct {
	Username     string   `json:"username"`
	PasswordHash string   `json:"password_hash"`
	Roles        []string `json:"roles"`
}

func createUser(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateUserParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	model := UserModel{
		Username:     params.Username,
		PasswordHash: params.PasswordHash,
		Roles:        params.Roles,
	}
	err = repository.createUser(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}
	return &restapi.Ok{&valuesMap{
		"id": model.ID,
	}}
}

// CreateAssignmentParams - parameters for the new contest assignment
type CreateAssignmentParams struct {
	UUID        string `json:"uuid"`
	ContestID   int64  `json:"contest_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func createAssignment(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateAssignmentParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	model := AssignmentFullModel{
		UUID:        params.UUID,
		ContestID:   params.ContestID,
		Title:       params.Title,
		Description: params.Description,
	}
	err = repository.createAssignment(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}
	return &restapi.Ok{&valuesMap{
		"id": model.ID,
	}}
}

// CreateTestCaseParams - parameters of the new assignment test case
type CreateTestCaseParams struct {
	UUID         string `json:"uuid"`
	AssignmentID int64  `json:"assignment_id"`
	Input        string `json:"input"`
	Expected     string `json:"expected"`
}

func createTestCase(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateTestCaseParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	c := ctx.(*apiContext)
	defer c.Close()

	repo, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	assignment, err := repo.getAssignment(params.AssignmentID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	_, err = c.builderService.RegisterTestCase(params.UUID, assignment.UUID, params.Input, params.Expected)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{nil}
}

// CreateAppointmentParams - parameters for the new contest assignment
type CreateAppointmentParams struct {
	GroupID   int64 `json:"group_id"`
	ContestID int64 `json:"contest_id"`
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
}

func assignGroupToContest(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateAppointmentParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	startTime := time.Unix(int64(params.StartTime), 0)
	endTime := time.Unix(int64(params.EndTime), 0)
	if startTime.After(endTime) {
		return &restapi.BadRequest{errors.New("contest start time cannot be bigger than end time")}
	}

	c := ctx.(*apiContext)
	defer c.Close()

	repo, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	model := AppointmentModel{
		GroupID:   params.GroupID,
		ContestID: params.ContestID,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
	}
	repo.createAppointment(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{&valuesMap{
		"id": model.ID,
	}}
}
