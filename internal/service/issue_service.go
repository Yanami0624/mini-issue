package service

import (
	"errors"
	"mini-issue/internal/dao"
	"mini-issue/internal/model"
)

type IssueService struct {
	dao *dao.IssueDAO
}

func NewIssueService(dao *dao.IssueDAO) *IssueService {
	return &IssueService{dao}
}

func (s *IssueService) CreateIssue(userid int64, req model.CreateIssueRequest) error {
	switch {
	case len(req.Title) == 0:
		return errors.New("empty title")
	case !isValidStatus(req.Status):
		return errors.New("invalid status")
	case isValidPriority(req.Priority):
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

	return nil
}

func (s *IssueService) UpdateIssue(userid, issueid int64, req model.UpdateIssueRequest) (*model.Issue, error) {
	issue, err := s.dao.GetByIssueID(issueid)
	switch {
	case err != nil:
		return nil, errors.New("invalid issueid")
	case issue.UserID != userid:
		return nil, errors.New("incorrect userid")
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

	return issue, nil
}

func (s *IssueService) DeleteIssue(userid, issueid int64) error {
	issue, err := s.dao.GetByIssueID(issueid)
	switch {
	case err != nil:
		return errors.New("invalid issueid")
	case issue.UserID != userid:
		return errors.New("incorrect userid")
	}
	err = s.dao.DeleteIssue(issueid)
	if err != nil {
		return err
	}
	return nil
}

func (s *IssueService) GetIssueByID(issueid int64) (*model.Issue, error) {
	issue, err := s.dao.GetByIssueID(issueid)
	if err != nil {
		return nil, errors.New("invalid issueid")
	}
	return issue, nil
}

func (s *IssueService) ListIssues(page, pagesize int) (*model.IssueListResponse, error) {
	page = max(page, 1)
	pagesize = max(pagesize, 10)
	pagesize = min(pagesize, 100)

	offset := (page - 1) * pagesize
	list, err := s.dao.ListIssues(pagesize, offset)
	if err != nil {
		return nil, err
	}

	total, _ := s.dao.CountIssues()
	return &model.IssueListResponse{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pagesize,
	}, nil
}

func isValidStatus(status string) bool {
	return status == "OPEN" || status == "IN_PROGRESS" || status == "DONE"
}
func isValidPriority(priority int) bool {
	return priority >= 1 && priority <= 3
}
