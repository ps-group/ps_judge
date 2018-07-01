package main

import (
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type BuildRepository interface {
	RegisterBuild(params RegisterBuildParams) error
	PullPendingBuild() (*PendingBuildResult, error)
	FinishBuild(params FinishBuildParams) error
	GetBuildInfo(key string) (*BuildInfoResult, error)
	GetAssignmentID(key string) (int, error)
}

type RegisterBuildParams struct {
	AssignmentID int64
	Key          string
	Language     language
	Source       string
	WebHookURL   string
}

type RegisterTestCaseParams struct {
	AssignmentID int64
	Key          string
	Input        string
	Output       string
	Expected     string
}

type FinishBuildParams struct {
	Key     string
	Succeed bool
	Score   int
	Report  string
}

type BuildInfoResult struct {
	Status     Status
	Score      sql.NullInt64
	Report     sql.NullString
	WebHookURL sql.NullString
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
	q := "INSERT INTO build (`assignment_id`, `key`, `status`, `language`, `source`, `web_hook_url`) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := r.Query(q, params.AssignmentID, params.Key, "pending", params.Language, params.Source, params.WebHookURL)
	return err
}

func (r *BuildRepositoryImpl) RegisterTestCase(params RegisterTestCaseParams) error {
	q := "INSERT INTO testcase (`assignment_id`, `key`, `input`, `output`, `expected`) VALUES (?, ?, ?, ?, ?)"
	_, err := r.Query(q, params.AssignmentID, params.Key, params.Input, params.Output, params.Expected)
	return err
}

func (r *BuildRepositoryImpl) PullPendingBuild() (*PendingBuildResult, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT key, language, source FROM build WHERE `status`='pending'")
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
	_, err = tx.Query("UPDATE build SET status='building' WHERE `key`=?", build.Key)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &build, nil
}

func (r *BuildRepositoryImpl) FinishBuild(params FinishBuildParams) error {
	q := "UPDATE build SET status=?, score=?, report=? WHERE `key`=?"
	status := "failed"
	if params.Succeed {
		status = "succeed"
	}
	_, err := r.Query(q, status, params.Score, params.Report, params.Key)

	return err
}

func (r *BuildRepositoryImpl) GetBuildInfo(key string) (*BuildInfoResult, error) {
	q := "SELECT status, score, report, web_hook_url FROM build WHERE `key`=?"
	rows, err := r.Query(q, key)
	if err != nil {
		return nil, errors.Wrap(err, "SQL query failed")
	}

	if !rows.Next() {
		return nil, errors.New("build with key '" + key + "' not found")
	}

	var item BuildInfoResult
	err = rows.Scan(&item.Status, &item.Score, &item.Report, &item.WebHookURL)
	if err != nil {
		return nil, errors.Wrap(err, "scan SQL result failed")
	}

	return &item, nil
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
