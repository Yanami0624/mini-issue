package service

import (
	"errors"
	"mini-issue/internal/dao"
	"mini-issue/internal/model"
	"mini-issue/internal/model"
	"mini-issue/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type IssueService struct {
	dao *dao.IssueDAO
}

func NewIssueService(dao *dao.IssueDAO) *IssueService {
	return &IssueService{dao}
}

func (service * IssueService) 
