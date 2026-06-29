package service

import (
	"context"
	"errors"
	"log"
	"mini-issue/internal/cache"
	"mini-issue/internal/dao"
	"mini-issue/internal/model"
	"mini-issue/internal/worker"
	"time"
)

const (
	issueDetailTTL = 5 * time.Minute
	issueListTTL   = 1 * time.Minute
)

type IssueService struct {
	dao *dao.IssueDAO
	cache *cache.IssueCache
	notif *worker.NotificationProducer
}

func NewIssueService(dao *dao.IssueDAO,
	cache *cache.IssueCache,
	notif *worker.NotificationProducer) *IssueService {
	return &IssueService{dao, cache, notif}
}

func (s *IssueService) CreateIssue(userid int64, req model.CreateIssueRequest) error {
	switch {
	case len(req.Title) == 0:
		return errors.New("empty title")
	case !isValidStatus(req.Status):
		return errors.New("invalid status")
	case !isValidPriority(req.Priority):
		return errors.New("invalid priority")
	}

	issue := model.Issue{
		UserID:   userid,
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status,
		Priority: req.Priority,
	}

	if err := s.dao.CreateIssue(issue); err != nil {
		return err
	}

	ctx, cancel := redisContext()
	defer cancel()
	s.cache.DeleteIssueList(ctx, 1, 10)
	
	task := model.NotificationTask {
		Type: model.Notification_IssueCreated,
		UserID: userid,
		IssueID: issue.ID,
		Title: issue.Title,
		CreatedAt: time.Now(),
	}
	s.notif.PushNotificationTask(ctx, task)

	return nil
}

func (s *IssueService) UpdateIssue(userid, issueid int64, req model.UpdateIssueRequest) (*model.Issue, error) {
	issue, err := s.dao.GetByIssueID(issueid)
	switch {
	case err != nil:
		return nil, errors.New("invalid issueid")
	case issue == nil:
		return nil, errors.New("issue not found")
	case issue.UserID != userid:
		return nil, errors.New("forbidden")
	case !isValidStatus(req.Status):
		return nil, errors.New("invalid issue status")
	case !isValidPriority(req.Priority):
		return nil, errors.New("invalid issue priority")
	}

	issue.Title = req.Title
	issue.Content = req.Content
	issue.Priority = req.Priority
	issue.Status = req.Status

	err = s.dao.UpdateIssue(issue)
	if err != nil {
		return nil, err
	}

	ctx, cancel := redisContext()
	defer cancel()
	s.cache.DeleteIssueDetail(ctx, issueid)

	return issue, nil
}

func (s *IssueService) DeleteIssue(userid, issueid int64) error {
	issue, err := s.dao.GetByIssueID(issueid)
	switch {
	case err != nil:
		return errors.New("invalid issueid")
	case issue == nil:
		return errors.New("issue not found")
	case issue.UserID != userid:
		return errors.New("forbidden")
	}
	err = s.dao.DeleteIssue(issueid)
	if err != nil {
		return err
	}

	ctx, cancel := redisContext()
	defer cancel()
	s.cache.DeleteIssueDetail(ctx, issueid)
	
	return nil
}

func (s *IssueService) GetIssueByID(issueid int64) (*model.Issue, error) {
	ctx, cancel := redisContext()
	defer cancel()

	cacheIssue, err := s.cache.GetIssueDetail(ctx, issueid)
	if err != nil {
		log.Printf("get issue detail cache failed: %v", err)
	} else if cacheIssue != nil {
		log.Printf("issue detail cache hit : id=%d", issueid)
		return cacheIssue, nil
	} else {
		log.Printf("issue detail cache miss: id=%d", issueid)
	}
	
	issue, err := s.dao.GetByIssueID(issueid)
	if err != nil {
		return nil, errors.New("invalid issueid")
	}
	if issue == nil {
		return nil, errors.New("issue not found")
	}
	s.cache.SetIssueDetail(ctx, issue, issueDetailTTL)
	return issue, nil
}

func (s *IssueService) ListIssues(page, pagesize int) (*model.IssueListResponse, error) {
	page = max(page, 1)
	pagesize = max(pagesize, 10)
	pagesize = min(pagesize, 100)

	// 不为list建缓存，因为不好管理，当一个issue变化，所有包含该issue的list缓存都会失效

	// ctx, cancel := redisContext()
	// defer cancel()
	// cacheList, err := s.cache.GetIssueList(ctx, page, pagesize)
	// if err != nil {
	// 	log.Printf("get issue list cache failed: %v", err)
	// } else if cacheList != nil {
	// 	log.Printf("issue list cache hit : page=%d", page)
	// } else {
	// 	log.Printf("issue list cache miss: page=%d", page)
	// }

	offset := (page - 1) * pagesize
	list, err := s.dao.ListIssues(pagesize, offset)
	if err != nil {
		return nil, err
	}

	total, err := s.dao.CountIssues()
	if err != nil {
		return nil, err
	}

	resp := model.IssueListResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pagesize,
	}
	// s.cache.SetIssueList(ctx, page, pagesize, &resp, issueListTTL)
	return &resp, nil
}

func isValidStatus(status string) bool {
	return status == "OPEN" || status == "IN_PROGRESS" || status == "DONE"
}
func isValidPriority(priority int) bool {
	return priority >= 1 && priority <= 3
}

func redisContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}