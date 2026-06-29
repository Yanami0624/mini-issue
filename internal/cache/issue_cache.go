package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"mini-issue/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type IssueCache struct {
	cli *redis.Client
}

func NewIssueCache(cli *redis.Client) *IssueCache {
	return &IssueCache{cli}
}

func IssueDetailKey(issueID int64) string {
	return fmt.Sprintf("issue:detail:%d", issueID)
}

func IssueListKey(page, pagesize int) string {
	return fmt.Sprintf("issue:list:page:%d:size:%d", page, pagesize)
}

func (c *IssueCache) GetIssueDetail(ctx context.Context, IssueID int64) (*model.Issue, error) {
	key := IssueDetailKey(IssueID)

	val, err := c.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var issue model.Issue
	json.Unmarshal([]byte(val), &issue)
	return &issue, nil
}

func (c *IssueCache) SetIssueDetail(ctx context.Context, issue *model.Issue, ttl time.Duration) error {
	key := IssueDetailKey(issue.ID)
	data, _ := json.Marshal(issue)
	return c.cli.Set(ctx, key, data, ttl).Err()
}

func (c *IssueCache) DeleteIssueDetail(ctx context.Context, issueID int64) error {
	key := IssueDetailKey(issueID)
	return c.cli.Del(ctx, key).Err()
}

func (c *IssueCache) GetIssueList(ctx context.Context, page, pagesize int) (*model.IssueListResponse, error) {
	key := IssueListKey(page, pagesize)

	val, err := c.cli.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result model.IssueListResponse
	json.Unmarshal([]byte(val), &result)
	return &result, nil
}

func (c *IssueCache) SetIssueList(ctx context.Context, page, pagesize int, result *model.IssueListResponse, ttl time.Duration) error {
	key := IssueListKey(page, pagesize)
	data, _ := json.Marshal(result)
	return c.cli.Set(ctx, key, data, ttl).Err()
}

func (c *IssueCache) DeleteIssueList(ctx context.Context, page, pageSize int) error {
	key := IssueListKey(page, pageSize)
	return c.cli.Del(ctx, key).Err()
}