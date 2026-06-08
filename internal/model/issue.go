package model

import "time"

type Issue struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Status    string    `db:"status" json:"status"`
	Priority  int       `db:"priority" json:"priority"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type CreateIssueRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Priority int    `json:"priority"`
}

type UpdateIssueRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Priority int    `json:"priority"`
}

type IssueListResponse struct {
	List     []Issue `json:"list"`
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
}