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
}

type RegisterBuildParams struct {
	Key        string
	Source     string
	WebHookURL string
}

type FinishBuildParams struct {
	Key     string
	Succeed bool
	Score   int
	Report  string
}

type BuildInfoResult struct {
	Succeed    bool
	Score      int
	Report     string
	WebHookURL string
}

type PendingBuildResult struct {
	Key    string
	Source string
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
	q := `INSERT INTO build SET key=?, source=?,  web_hook_url=?, status='pending'`
	rows, err := r.db.Query(q, params.Key, params.Source, params.WebHookURL)
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
	rows, err := tx.Query(`SELECT key, source FROM build WHERE status='pending'`)
	if err != nil {
		return nil, err
	}
	// If no pending build, return nil.
	if !rows.Next() {
		return nil, nil
	}
	var build PendingBuildResult
	err = rows.Scan(&build.Key, &build.Source)
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
	err = rows.Scan(&item.Succeed, &item.Score, &item.Report, &item.WebHookURL)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
