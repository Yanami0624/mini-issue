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

func (dao *IssueDAO) CreateIssue(issue model.Issue) error {
	timestamp := time.Now()
	_, err := dao.db.Exec(
		"insert into issues (user_id, title, content, status, priority, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?)",
		issue.UserID, issue.Title, issue.Content, issue.Status, issue.Priority, timestamp, timestamp)
	if err != nil {
		fmt.Println("failed: CreateIssue()", err)
	}
	return err
}

func (dao *IssueDAO) GetByIssueID(issueid int64) (*model.Issue, error) {
	var issue model.Issue
	query := `
		select id, user_id, title, content, status, priority, created_at, updated_at
		from issues
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
		SET title = ?, content = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := dao.db.Exec(
		query,
		issue.Title,
		issue.Content,
		issue.Status,
		issue.Priority,
		time.Now(),
		issue.ID,
	)

	return err
}

func (dao *IssueDAO) DeleteIssue(issueid int64) error {
	query := `
		delete from issues
		where id = ?
	`

	_, err := dao.db.Exec(query, issueid)
	if err != nil {
		fmt.Println("failed: GetByUsername()", err)
		return err
	}

	return nil
}

func (dao *IssueDAO) ListIssues(limit, offset int) ([]model.Issue, error) {
	var issues []model.Issue
	query := `
		select id, user_id, title, content, status, priority, created_at, updated_at
		from issues
		order by id desc
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
