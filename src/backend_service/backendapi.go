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

func parseID(req restapi.Request) (int64, error) {
	return strconv.ParseInt(req.Var("id"), 10, 64)
}

// CreateResponse - contains ID of the created entity
type CreateResponse struct {
	ID int64 `json:"id"`
}

// UserInfo contains user info response
type UserInfo struct {
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
	ContestID int64    `json:"contest_id"`
}

// BriefSolutionInfo - contains brief solution info
type BriefSolutionInfo struct {
	AssignmentID    int64  `json:"assignment_id"`
	AssignmentTitle string `json:"assignment_title"`
	CommitID        int64  `json:"commit_id"`
	Score           int64  `json:"score"`
	BuildStatus     string `json:"build_status"`
}

// LoginUserParams - contains user login parameters
type LoginUserParams struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

// LoginUserResponse - contains user login response
type LoginUserResponse struct {
	Succeed bool     `json:"succeed"`
	UserID  int64    `json:"user_id"`
	User    UserInfo `json:"user"`
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
		return &restapi.Ok{
			LoginUserResponse{
				Succeed: true,
				UserID:  info.ID,
				User: UserInfo{
					Username:  info.Username,
					Roles:     info.Roles,
					ContestID: info.ContestID,
				},
			},
		}
	}

	return &restapi.Unauthorized{
		LoginUserResponse{
			Succeed: false,
		},
	}
}

func getUserInfo(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	userID, err := parseID(req)
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	info, err := repository.getUserInfo(userID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{UserInfo{
		Username:  info.Username,
		Roles:     info.Roles,
		ContestID: info.ContestID,
	}}
}

func getUserSolutions(ctx interface{}, req restapi.Request) restapi.Response {
	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	userID, err := parseID(req)
	if err != nil {
		return &restapi.BadRequest{errors.Wrap(err, "invalid id")}
	}

	solutions, err := repository.getUserSolutions(userID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	userInfo, err := repository.getUserInfo(userID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	assignments, err := repository.getContestAssignments(userInfo.ContestID)
	if err != nil {
		return &restapi.InternalError{err}
	}

	assignmentTitles := make(map[int64]string)
	for _, assignment := range assignments {
		assignmentTitles[assignment.ID] = assignment.Title
	}

	var results []BriefSolutionInfo
	for _, solution := range solutions {
		commit, err := repository.getLastCommit(solution.ID)
		if err != nil {
			return &restapi.InternalError{err}
		}
		result := BriefSolutionInfo{
			AssignmentID:    solution.AssignmentID,
			AssignmentTitle: assignmentTitles[solution.AssignmentID],
			Score:           solution.Score,
			CommitID:        commit.ID,
			BuildStatus:     commit.BuildStatus,
		}
		results = append(results, result)
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
	userID, err := parseID(req)
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

// AssignmentInfo - brief info about assignment
type AssignmentInfo struct {
	ID        int64  `json:"id"`
	ContestID int64  `json:"contest_id"`
	UUID      string `json:"uuid"`
	Title     string `json:"title"`
}

func getContestAssignments(ctx interface{}, req restapi.Request) restapi.Response {
	contestID, err := parseID(req)
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

	var infos []AssignmentInfo
	for _, assignment := range assignments {
		info := AssignmentInfo{
			ID:        assignment.ID,
			ContestID: assignment.ContestID,
			UUID:      assignment.UUID,
			Title:     assignment.Title,
		}
		infos = append(infos, info)
	}
	return &restapi.Ok{infos}
}

// FullAssignmentInfo - result of assignment info request
type FullAssignmentInfo struct {
	ID          int64  `json:"id"`
	ContestID   int64  `json:"contest_id"`
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func getAssignmentInfo(ctx interface{}, req restapi.Request) restapi.Response {
	assignmentID, err := parseID(req)
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

	result := FullAssignmentInfo{
		ID:          assignmentID,
		ContestID:   assignment.ContestID,
		UUID:        assignment.UUID,
		Title:       assignment.Title,
		Description: assignment.Description,
	}
	return &restapi.Ok{result}
}

// CreateContestParams - parameters for the new contest
type CreateContestParams struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func createContest(ctx interface{}, req restapi.Request) restapi.Response {
	var params CreateContestParams
	err := req.ReadJSON(&params)
	if err != nil {
		return &restapi.BadRequest{err}
	}

	if params.StartTime.After(params.EndTime) {
		return &restapi.BadRequest{errors.New("contest start time cannot be bigger than end time")}
	}

	c := ctx.(*apiContext)
	defer c.Close()
	repository, err := c.ConnectDB()
	if err != nil {
		return &restapi.InternalError{err}
	}

	model := ContestModel{
		Title:     params.Title,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
	}
	err = repository.createContest(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}
	return &restapi.Ok{&CreateResponse{
		ID: model.ID,
	}}
}

// CreateUserParams - parameters for the new user
type CreateUserParams struct {
	Username     string   `json:"username"`
	PasswordHash string   `json:"password_hash"`
	Roles        []string `json:"roles"`
	ContestID    int64    `json:"contest_id"`
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
		ContestID:    params.ContestID,
	}
	err = repository.createUser(&model)
	if err != nil {
		return &restapi.InternalError{err}
	}
	return &restapi.Ok{&CreateResponse{
		ID: model.ID,
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
	return &restapi.Ok{&CreateResponse{
		ID: model.ID,
	}}
}

// CreateTestCaseParams - parameters of the new assignment test case
type CreateTestCaseParams struct {
	UUID         string `json:"uuid"`
	AssignmentID string `json:"assignment_uuid"`
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

	_, err = c.builderService.RegisterTestCase(params.UUID, params.AssignmentID, params.Input, params.Expected)
	if err != nil {
		return &restapi.InternalError{err}
	}

	return &restapi.Ok{nil}
}
