package main

import (
	"database/sql"

	"github.com/pkg/errors"
)

// BuildRepository - represents builder database model
type BuildRepository interface {
	RegisterBuild(params RegisterBuildParams) error
	PullPendingBuild() (*PendingBuildResult, error)
	AddBuildReport(params BuildReport) error
	GetBuildReport(key string) (*BuildReport, error)
	GetAssignmentID(key string) (int, error)
}

// RegisterBuildParams - parameters for DB request
type RegisterBuildParams struct {
	AssignmentID int64
	Key          string
	Language     language
	Source       string
}

// RegisterTestCaseParams - parameters for DB request
type RegisterTestCaseParams struct {
	AssignmentID int64
	Key          string
	Input        string
	Expected     string
}

// BuildReport - parameters for DB request
type BuildReport struct {
	Key         string
	Exception   string
	BuildLog    string
	TestsLog    string
	TestsPassed int64
	TestsTotal  int64
	Status      Status
}

// PendingBuildResult - parameters for DB request
type PendingBuildResult struct {
	AssignmentID int
	Key          string
	Source       string
	Language     language
}

// BuilderRepository - models builder service database
type BuilderRepository struct {
	db *sql.DB
}

// NewBuilderRepository - creates repository with given database connection
func NewBuilderRepository(db *sql.DB) *BuilderRepository {
	var r BuilderRepository
	r.db = db
	return &r
}

func (r *BuilderRepository) query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return nil, errors.Wrap(err, "sql query '"+query+"' failed")
	}
	return rows, nil
}

func (r *BuilderRepository) prepare(query string) (*sql.Stmt, error) {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "sql prepare '"+query+"' failed")
	}
	return stmt, nil
}

// RegisterBuild - registers new build task
func (r *BuilderRepository) RegisterBuild(params RegisterBuildParams) error {
	q := "INSERT INTO build (`assignment_id`, `key`, `status`, `language`, `source`) VALUES (?, ?, ?, ?, ?)"
	_, err := r.query(q, params.AssignmentID, params.Key, "pending", params.Language, params.Source)
	return err
}

// RegisterTestCase - registers new test case for the assignment solutions
func (r *BuilderRepository) RegisterTestCase(params RegisterTestCaseParams) error {
	q := "INSERT INTO testcase (`assignment_id`, `key`, `input`, `expected`) VALUES (?, ?, ?, ?)"
	_, err := r.query(q, params.AssignmentID, params.Key, params.Input, params.Expected)
	return err
}

// PullPendingBuild - pulls one pending build from database and turns it into 'building' status
func (r *BuilderRepository) PullPendingBuild() (*PendingBuildResult, error) {
	rows, err := r.query("SELECT `assignment_id`, `key`, `language`, `source` FROM build WHERE `status` = 'pending' LIMIT 1")
	if err != nil {
		return nil, err
	}
	// If no pending build, return nil.
	if !rows.Next() {
		return nil, nil
	}
	var build PendingBuildResult
	err = rows.Scan(&build.AssignmentID, &build.Key, &build.Language, &build.Source)
	if err != nil {
		return nil, err
	}
	_, err = r.query("UPDATE build SET status='building' WHERE `key`=?", build.Key)
	if err != nil {
		return nil, err
	}

	return &build, nil
}

// AddBuildReport - adds finished build report
func (r *BuilderRepository) AddBuildReport(params BuildReport) error {
	buildID, err := r.GetBuildID(params.Key)

	_, err = r.query(
		"INSERT INTO report (`build_id`, `tests_passed`, `tests_total`, `exception`, `build_log`, `tests_log`) VALUES (?, ?, ?, ?, ?, ?)",
		buildID, params.TestsPassed, params.TestsTotal, params.Exception, params.BuildLog, params.TestsLog)
	if err != nil {
		return errors.Wrap(err, "SQL INSERT query failed")
	}
	_, err = r.query("UPDATE build SET status=? WHERE `id`=?", params.Status, buildID)
	if err != nil {
		return errors.Wrap(err, "SQL UPDATE query failed")
	}

	return err
}

// GetTestCases - returns list of test cases for the assignment solutions
func (r *BuilderRepository) GetTestCases(assignmentID int) ([]TestCase, error) {
	var cases []TestCase
	rows, err := r.query("SELECT input, expected FROM testcase WHERE `assignment_id`=?", assignmentID)
	if err != nil {
		return cases, errors.Wrap(err, "SQL SELECT query failed")
	}
	for rows.Next() {
		var result TestCase
		err = rows.Scan(&result.Input, &result.Expected)
		if err != nil {
			return cases, errors.Wrap(err, "scan SQL result failed")
		}
		cases = append(cases, result)
	}
	return cases, nil
}

// GetBuildStatus - returns build status string
func (r *BuilderRepository) GetBuildStatus(key string) (Status, error) {
	rows, err := r.query("SELECT status FROM build WHERE `key`=?", key)
	if err != nil {
		return "", errors.Wrap(err, "SQL SELECT query failed")
	}
	if !rows.Next() {
		return "", errors.New("build with key '" + key + "' not found")
	}
	var status Status
	err = rows.Scan(&status)
	if err != nil {
		return "", errors.Wrap(err, "scan SQL result failed")
	}
	return status, nil
}

// GetBuildReport - returns detailed report for finished build
func (r *BuilderRepository) GetBuildReport(key string) (*BuildReport, error) {
	rows, err := r.query("SELECT id, status FROM build WHERE `key`=?", key)
	if err != nil {
		return nil, errors.Wrap(err, "SQL SELECT query failed")
	}
	if !rows.Next() {
		return nil, errors.New("build with key '" + key + "' not found")
	}
	var buildID int64
	var status Status
	err = rows.Scan(&buildID, &status)
	if err != nil {
		return nil, errors.Wrap(err, "scan SQL result failed")
	}

	rows, err = r.query("SELECT tests_passed, tests_total, exception, build_log, tests_log FROM report WHERE `build_id`=?", buildID)
	if err != nil {
		return nil, errors.Wrap(err, "SQL SELECT query failed")
	}
	if !rows.Next() {
		return nil, errors.New("report for build with key '" + key + "' not found")
	}

	var report BuildReport
	report.Key = key
	report.Status = status
	err = rows.Scan(&report.TestsPassed, &report.TestsTotal, &report.Exception, &report.BuildLog, &report.TestsLog)
	if err != nil {
		return nil, errors.Wrap(err, "scan SQL result failed")
	}

	return &report, nil
}

// GetAssignmentID - returns assignment ID for given cross-service unique key
func (r *BuilderRepository) GetAssignmentID(key string) (int64, error) {
	rows, err := r.query("SELECT id FROM assignment WHERE `key`=?", key)

	if !rows.Next() {
		stmt, err := r.prepare("INSERT INTO assignment (`key`) VALUES (?)")
		if err != nil {
			return 0, err
		}
		res, err := stmt.Exec(key)
		if err != nil {
			return 0, err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		return id, nil
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetBuildID - returns build ID for given cross-service unique key
func (r *BuilderRepository) GetBuildID(key string) (int64, error) {
	rows, err := r.query("SELECT id FROM build WHERE `key`=?", key)
	if err != nil {
		return 0, errors.Wrap(err, "SQL SELECT query failed")
	}
	if !rows.Next() {
		return 0, errors.New("build with key '" + key + "' not found")
	}
	var buildID int64
	err = rows.Scan(&buildID)
	if err != nil {
		return 0, errors.Wrap(err, "scan SQL result failed")
	}
	return buildID, nil
}
