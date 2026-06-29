package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mini-issue/internal/model"
	"mini-issue/internal/service"
	"mini-issue/pkg/response"
)

type IssueController struct {
	s *service.IssueService
}

func NewIssueController(is *service.IssueService) *IssueController {
	return &IssueController{is}
}

func (c *IssueController) CreateIssue(ctx *gin.Context) {
	userid, err := GetUserID(ctx)
	if err != nil {
		response.Fail(ctx, 401, "unauthorized")
		return
	}

	var req model.CreateIssueRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(ctx, 400, "invalid request body")
		return
	}

	err = c.s.CreateIssue(userid, req)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, nil)
}

func (c *IssueController) GetIssue(ctx *gin.Context) {
	issueid, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Fail(ctx, 400, "invalid issue id")
		return
	}
	issue, err := c.s.GetIssueByID(issueid)
	if err != nil {
		if err.Error() == "issue not found" {
			response.Fail(ctx, 404, err.Error())
			return
		}
		response.Fail(ctx, 400, err.Error())
		return
	}
	response.Success(ctx, issue)
}

func (c *IssueController) ListIssues(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))

	result, err := c.s.ListIssues(page, pageSize)
	if err != nil {
		response.Fail(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, result)
}

func (c *IssueController) UpdateIssue(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		response.Fail(ctx, 401, "unauthorized")
		return
	}

	issueID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Fail(ctx, 400, "invalid issue id")
		return
	}

	var req model.UpdateIssueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, 400, "invalid request body")
		return
	}

	issue, err := c.s.UpdateIssue(userID, issueID, req)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Fail(ctx, 403, err.Error())
			return
		}
		if err.Error() == "issue not found" {
			response.Fail(ctx, 404, err.Error())
			return
		}
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, issue)
}

func (c *IssueController) DeleteIssue(ctx *gin.Context) {
	userID, err := GetUserID(ctx)
	if err != nil {
		response.Fail(ctx, 401, "unauthorized")
		return
	}

	issueID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Fail(ctx, 400, "invalid issue id")
		return
	}

	err = c.s.DeleteIssue(userID, issueID)
	if err != nil {
		if err.Error() == "forbidden" {
			response.Fail(ctx, 403, err.Error())
			return
		}
		if err.Error() == "issue not found" {
			response.Fail(ctx, 404, err.Error())
			return
		}
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, nil)
}
