package main

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type BuildRepository interface {
	RegisterBuild(params RegisterBuildParams) error
	PullPendingBuild() (*PendingBuildResult, error)
	FinishBuild(params FinishBuildParams) error
	GetBuildInfo(key string) (*BuildInfoResult, error)
	GetAssignmentId(key string) (int, error)
}

type RegisterBuildParams struct {
	AssignmentId int
	Key          string
	Language     language
	Source       string
	WebHookURL   string
}

type RegisterTestCaseParams struct {
	AssignmentId int
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
	Score      int
	Report     string
	WebHookURL string
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

func (r *BuildRepositoryImpl) RegisterBuild(params RegisterBuildParams) error {
	q := `INSERT INTO build SET assignment_id=?, key=?, status='pending', language=?, source=?,  web_hook_url=?`
	rows, err := r.db.Query(q, params.Key, params.AssignmentId, params.Language, params.Source, params.WebHookURL)
	if err != nil {
		rows.Close()
	}
	return err
}

func (r *BuildRepositoryImpl) RegisterTestCase(params RegisterTestCaseParams) error {
	q := `INSERT INTO testcase SET assignment_id=?, key=?, input=?, output=?, expected=?`
	rows, err := r.db.Query(q, params.AssignmentId, params.Key, params.Input, params.Output, params.Expected)
	if err != nil {
		rows.Close()
	}
	return err
}

func (r *BuildRepositoryImpl) PullPendingBuild() (*PendingBuildResult, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(`SELECT key, language, source FROM build WHERE status='pending'`)
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
	_, err = tx.Query(`UPDATE build SET status='building' WHERE key=?`, build.Key)
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
	q := `UPDATE build SET status=?, score=?, report=? WHERE key=?`
	status := "failed"
	if params.Succeed {
		status = "succeed"
	}
	_, err := r.db.Query(q, status, params.Score, params.Report, params.Key)

	return err
}

func (r *BuildRepositoryImpl) GetBuildInfo(key string) (*BuildInfoResult, error) {
	q := `SELECT succeed, score, report, web_hook_url FROM video WHERE key = ?`
	rows, err := r.db.Query(q, key)

	if !rows.Next() {
		return nil, errors.New("video with key '" + key + "' not found")
	}

	var item BuildInfoResult
	err = rows.Scan(&item.Status, &item.Score, &item.Report, &item.WebHookURL)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *BuildRepositoryImpl) GetAssignmentId(key string) (int, error) {
	q := `SELECT id FROM build WHERE key = ?`
	rows, err := r.db.Query(q, key)

	if !rows.Next() {
		return 0, errors.New("build with key '" + key + "' not found")
	}

	var id int
	err = rows.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
