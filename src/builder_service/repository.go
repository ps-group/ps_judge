package main

import (
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type BuildRepository interface {
	RegisterBuild(params RegisterBuildParams) error
	PullPendingBuild() (*PendingBuildResult, error)
	AddBuildReport(params BuildReport) error
	GetBuildReport(key string) (*BuildReport, error)
	GetAssignmentID(key string) (int, error)
}

type RegisterBuildParams struct {
	AssignmentID int64
	Key          string
	Language     language
	Source       string
}

type RegisterTestCaseParams struct {
	AssignmentID int64
	Key          string
	Input        string
	Expected     string
}

type BuildReport struct {
	Key         string
	Exception   string
	BuildLog    string
	TestsPassed int64
	TestsTotal  int64
	Status      Status
}

type PendingBuildResult struct {
	Key      string
	Source   string
	Language language
}

type BuildRepositoryImpl struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *BuildRepositoryImpl {
	var r BuildRepositoryImpl
	r.db = db
	return &r
}

func (r *BuildRepositoryImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := r.db.Query(query, args...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return nil, errors.Wrap(err, "sql query '"+query+"' failed")
	}
	return rows, nil
}

func (r *BuildRepositoryImpl) Prepare(query string) (*sql.Stmt, error) {
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "sql prepare '"+query+"' failed")
	}
	return stmt, nil
}

func (r *BuildRepositoryImpl) RegisterBuild(params RegisterBuildParams) error {
	q := "INSERT INTO build (`assignment_id`, `key`, `status`, `language`, `source`) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Query(q, params.AssignmentID, params.Key, "pending", params.Language, params.Source)
	return err
}

func (r *BuildRepositoryImpl) RegisterTestCase(params RegisterTestCaseParams) error {
	q := "INSERT INTO testcase (`assignment_id`, `key`, `input`, `expected`) VALUES (?, ?, ?, ?)"
	_, err := r.Query(q, params.AssignmentID, params.Key, params.Input, params.Expected)
	return err
}

func (r *BuildRepositoryImpl) PullPendingBuild() (*PendingBuildResult, error) {
	rows, err := r.Query("SELECT `key`, `language`, `source` FROM build WHERE `status` = 'pending' LIMIT 1")
	if err != nil {
		return nil, err
	}
	// If no pending build, return nil.
	if !rows.Next() {
		return nil, nil
	}
	var build PendingBuildResult
	err = rows.Scan(&build.Key, &build.Language, &build.Source)
	if err != nil {
		return nil, err
	}
	_, err = r.Query("UPDATE build SET status='building' WHERE `key`=?", build.Key)
	if err != nil {
		return nil, err
	}

	return &build, nil
}

func (r *BuildRepositoryImpl) AddBuildReport(params BuildReport) error {
	buildID, err := r.getBuildID(params.Key)

	_, err = r.Query(
		"INSERT INTO report (`build_id`, `tests_passed`, `tests_total`, `exception`, `build_log`) VALUES (?, ?, ?, ?, ?)",
		buildID, params.TestsPassed, params.TestsTotal, params.Exception, params.BuildLog)
	if err != nil {
		return errors.Wrap(err, "SQL INSERT query failed")
	}
	_, err = r.Query("UPDATE build SET status=? WHERE `id`=?", params.Status, buildID)
	if err != nil {
		return errors.Wrap(err, "SQL UPDATE query failed")
	}

	return err
}

func (r *BuildRepositoryImpl) GetBuildStatus(key string) (Status, error) {
	rows, err := r.Query("SELECT status FROM build WHERE `key`=", key)
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

func (r *BuildRepositoryImpl) GetBuildReport(key string) (*BuildReport, error) {
	rows, err := r.Query("SELECT id, status FROM build WHERE `key`=?", key)
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

	rows, err = r.Query("SELECT status, tests_passed, tests_total, exception, build_log FROM report WHERE `build_id`=?", buildID)
	if err != nil {
		return nil, errors.Wrap(err, "SQL SELECT query failed")
	}
	if !rows.Next() {
		return nil, errors.New("build with key '" + key + "' not found")
	}

	var report BuildReport
	report.Key = key
	report.Status = status
	err = rows.Scan(&report.Status, &report.TestsPassed, &report.TestsTotal, &report.Exception, &report.BuildLog)
	if err != nil {
		return nil, errors.Wrap(err, "scan SQL result failed")
	}

	return &report, nil
}

func (r *BuildRepositoryImpl) GetAssignmentID(key string) (int64, error) {
	rows, err := r.Query("SELECT id FROM assignment WHERE `key`=?", key)

	if !rows.Next() {
		stmt, err := r.Prepare("INSERT INTO assignment (`key`) VALUES (?)")
		if err != nil {
			return 0, err
		}
		res, err := stmt.Exec(key)
		if err != nil {
			return 0, err
		}
		id, _ := res.LastInsertId()
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

func (r *BuildRepositoryImpl) getBuildID(key string) (int64, error) {
	rows, err := r.Query("SELECT id FROM build WHERE `key`=?", key)
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
