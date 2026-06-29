package model

import "time"

type NotificationTask struct {
	Type      string    `json:"type"`
	UserID    int64     `json:"user_id"`
	IssueID   int64     `json:"issue_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	Notification_IssueCreated = "ISSUE_CREATED"
)

