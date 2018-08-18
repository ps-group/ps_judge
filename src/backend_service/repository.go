package main

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// BackendRepository - models builder service database
type BackendRepository struct {
	db *sql.DB
}

// NewBackendRepository - creates repository with given database connection
func NewBackendRepository(db *sql.DB) *BackendRepository {
	var r BackendRepository
	r.db = db
	return &r
}

func (r *BackendRepository) query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return nil, errors.Wrap(err, "sql query '"+query+"' failed")
	}
	return rows, nil
}

func (r *BackendRepository) prepare(query string) (*sql.Stmt, error) {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "sql prepare '"+query+"' failed")
	}
	return stmt, nil
}

// UserModel - models user info in database
type UserModel struct {
	ID           int64
	ContestID    int64
	Username     string
	PasswordHash string
	Roles        []string
}

// SolutionModel - models solution in database
type SolutionModel struct {
	ID           int64
	UserID       int64
	AssignmentID int64
	Score        int64
}

// AssignmentInfoModel - models brief info about assignment in database
type AssignmentInfoModel struct {
	ID        int64
	ContestID int64
	UUID      string
	Title     string
}

// AssignmentFullModel - models full info about assignment in database
type AssignmentFullModel struct {
	ID          int64
	ContestID   int64
	UUID        string
	Title       string
	Description string
}

func (r *BackendRepository) getUserInfo(id int64) (*UserModel, error) {
	rows, err := r.query("SELECT `username`, `password`, `active_contest_id`, `roles` FROM user WHERE `id`=?", id)
	if err != nil {
		return nil, err
	}
	// If no such user, return nil.
	if !rows.Next() {
		return nil, nil
	}
	var user UserModel
	user.ID = id
	var roles []byte
	err = rows.Scan(&user.Username, &user.PasswordHash, &user.ContestID, &roles)
	user.Roles = []string{string(roles)}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *BackendRepository) getUserInfoByUsername(username string) (*UserModel, error) {
	rows, err := r.query("SELECT `id`, `password`, `active_contest_id`, `roles` FROM user WHERE `username`=?", username)
	if err != nil {
		return nil, err
	}
	// If no such user, return nil.
	if !rows.Next() {
		return nil, nil
	}
	var user UserModel
	user.Username = username
	var roles []byte
	err = rows.Scan(&user.ID, &user.PasswordHash, &user.ContestID, &roles)
	user.Roles = strings.Split(string(roles), ",")
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *BackendRepository) getUserSolutions(userID int64) ([]SolutionModel, error) {
	var results []SolutionModel
	rows, err := r.query("SELECT `id`, `assignment_id`, `score` FROM solution WHERE user_id=?", userID)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result SolutionModel
		result.UserID = userID
		err = rows.Scan(&result.ID, &result.AssignmentID, &result.Score)
		if err != nil {
			return results, errors.Wrap(err, "failed to scan SQL rows")
		}
		results = append(results, result)
	}
	return results, nil
}

func (r *BackendRepository) getUserAssignmentSolution(userID int64, assignmentID int64) (*SolutionModel, error) {
	rows, err := r.query("SELECT `id`, `score` FROM solution WHERE user_id=? AND assignment_id=? LIMIT 1", userID, assignmentID)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, nil
	}

	var result SolutionModel
	result.UserID = userID
	result.AssignmentID = assignmentID
	err = rows.Scan(&result.ID, &result.Score)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) getSolution(id int64) (*SolutionModel, error) {
	rows, err := r.query("SELECT `score`, `user_id`, `assignment_id` FROM solution WHERE id=?", id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("no solution with given ID")
	}

	var result SolutionModel
	result.ID = id
	err = rows.Scan(&result.Score, &result.UserID, &result.AssignmentID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) createSolution(userID int64, assignmentID int64) (*SolutionModel, error) {
	stmt, err := r.prepare("INSERT INTO solution (user_id, assignment_id, score) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(userID, assignmentID, 0)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &SolutionModel{
		ID:           id,
		UserID:       userID,
		AssignmentID: assignmentID,
	}, nil
}

func (r *BackendRepository) updateSolutionScore(solutionID int64, score int64) error {
	_, err := r.query("UPDATE solution SET score=? WHERE id=?", score, solutionID)
	return err
}

func (r *BackendRepository) getAssignment(assignmentID int64) (*AssignmentInfoModel, error) {
	rows, err := r.query("SELECT `contest_id`, `uuid`, `title` FROM assignment WHERE id=?", assignmentID)
	if err != nil {
		return nil, err
	}

	// No commit - it's OK.
	if !rows.Next() {
		return nil, nil
	}

	var result AssignmentInfoModel
	result.ID = assignmentID
	err = rows.Scan(&result.ContestID, &result.UUID, &result.Title)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) getAssignmentFull(assignmentID int64) (*AssignmentFullModel, error) {
	rows, err := r.query("SELECT `contest_id`, `uuid`, `title`, `article` FROM assignment WHERE id=?", assignmentID)
	if err != nil {
		return nil, err
	}

	// No commit - it's OK.
	if !rows.Next() {
		return nil, nil
	}

	var result AssignmentFullModel
	result.ID = assignmentID
	err = rows.Scan(&result.ContestID, &result.UUID, &result.Title, &result.Description)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) getContestAssignments(contestID int64) ([]AssignmentInfoModel, error) {
	var results []AssignmentInfoModel
	rows, err := r.query("SELECT `id`, `uuid`, `title` FROM assignment WHERE contest_id=?", contestID)
	if err != nil {
		return results, err
	}

	for rows.Next() {
		var result AssignmentInfoModel
		result.ContestID = contestID
		err = rows.Scan(&result.ID, &result.UUID, &result.Title)
		if err != nil {
			return results, errors.Wrap(err, "failed to scan SQL rows")
		}
		results = append(results, result)
	}
	return results, nil
}

// CommitModel - represents commit in database
type CommitModel struct {
	ID          int64
	SolutionID  int64
	BuildStatus string
	BuildScore  int64
}

func (r *BackendRepository) createCommit(solutionID int64, uuid string) error {
	_, err := r.query("INSERT INTO commit (solution_id, uuid) VALUES (?, ?)", solutionID, uuid)
	return err
}

func (r *BackendRepository) getLastCommit(solutionID int64) (*CommitModel, error) {
	rows, err := r.query("SELECT `id`, `build_status`, `build_score` FROM commit WHERE solution_id=? ORDER BY `id` DESC LIMIT 1", solutionID)
	if err != nil {
		return nil, err
	}

	// No commit - it's OK.
	if !rows.Next() {
		return nil, nil
	}

	var score sql.NullInt64
	var result CommitModel
	result.SolutionID = solutionID
	err = rows.Scan(&result.ID, &result.BuildStatus, &score)
	result.BuildScore = score.Int64
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) getCommitInfoByUUID(uuid string) (*CommitModel, error) {
	rows, err := r.query("SELECT `id`, `build_status`, `build_score`, `solution_id` FROM commit WHERE uuid=?", uuid)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errors.New("no commit with uuid=" + uuid)
	}

	var score sql.NullInt64
	var result CommitModel
	err = rows.Scan(&result.ID, &result.BuildStatus, &score, &result.SolutionID)
	result.BuildScore = score.Int64
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan SQL rows")
	}

	return &result, nil
}

func (r *BackendRepository) updateCommit(model *CommitModel) error {
	_, err := r.query("UPDATE commit SET build_status=?, build_score=? WHERE id=?", model.BuildStatus, model.BuildScore, model.ID)
	return err
}

// ContestModel - models contest in database
type ContestModel struct {
	ID        int64
	Title     string
	StartTime time.Time
	EndTime   time.Time
}

// Creates contest and sets ID if succeed
func (r *BackendRepository) createContest(model *ContestModel) error {
	stmt, err := r.prepare("INSERT INTO contest (title, start_time, end_time) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(model.Title, model.StartTime, model.EndTime)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	model.ID = id
	return nil
}

// Creates user and sets ID if succeed
func (r *BackendRepository) createUser(model *UserModel) error {
	stmt, err := r.prepare("INSERT INTO user (username, password, roles, active_contest_id) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	roles := []byte(strings.Join(model.Roles, ","))
	res, err := stmt.Exec(model.Username, model.PasswordHash, roles, model.ContestID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	model.ID = id
	return nil
}

func (r *BackendRepository) createAssignment(model *AssignmentFullModel) error {
	stmt, err := r.prepare("INSERT INTO assignment (uuid, contest_id, title, article) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(model.UUID, model.ContestID, model.Title, model.Description)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	model.ID = id
	return nil
}
