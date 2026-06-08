package dao

import (
	"database/sql"
	"fmt"
	"mini-issue/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type IssueDAO struct {
	db *sqlx.DB
}

func NewIssueDAO(db *sqlx.DB) *IssueDAO {
	return &IssueDAO{db}
}

func (dao *IssueDAO) CreateIssue(
	userid int64,
	title, content, status string,
	priority int) error {
	timestamp := time.Now()
	_, err := dao.db.Exec(
		"insert into issue (user_id, title, content, status, priority, created_at, updated_at)",
		userid, title, content, status, priority, timestamp, timestamp)
	if err != nil {
		fmt.Println("failed: CreateIssue()", err)
	}
	return err
}

func (dao *IssueDAO) GetByIssueID(issueid int64) (*model.Issue, error) {
	var issue model.Issue
	query := `
		select *
		from issue
		where id = ?
	`

	err := dao.db.Get(&issue, query, issueid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println("failed: GetByIssueID()", err)
		return nil, err
	}

	return &issue, nil
}

func (dao *IssueDAO) UpdateIssue(issue *model.Issue) error {
	query := `
		UPDATE issues
		SET title = ?, content = ?, status = ?, priority = ?
		WHERE id = ?
	`

	_, err := dao.db.Exec(
		query,
		issue.Title,
		issue.Content,
		issue.Status,
		issue.Priority,
		issue.ID,
	)

	return err
}

func (dao *IssueDAO) DeleteIssue(issueid int64, update_comment string) error {
	query := `
		delete from issue
		where id = ?
	`

	_, err := dao.db.Exec(query, update_comment, issueid)
	if err != nil {
		fmt.Println("failed: GetByUsername()", err)
		return err
	}

	return nil
}

func (dao *IssueDAO) ListIssues(limit, offset int) ([]model.Issue, error) {
	var issues []model.Issue
	query := `
		select * 
		from issues
		oedered by id
		limit ? offset ?
	`

	err := dao.db.Select(&issues, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return issues, nil
}

func (dao *IssueDAO) CountIssues() (int64, error) {
	var total int64
	query := `
		select count(*)
		from issues
	`

	err := dao.db.Get(&total, query)
	return total, err
}
