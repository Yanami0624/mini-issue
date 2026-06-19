package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"mini-issue/internal/model"
	"mini-issue/internal/service"
	"mini-issue/pkg/response"
)

type UserController struct {
	s *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{us}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req model.RegisterRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(ctx, 400, "invalid request body")
		return
	}

	err = c.s.Register(req)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, nil)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req model.LoginRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		response.Fail(ctx, 400, "invalid request body")
		return
	}

	resp, err := c.s.Login(req)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, resp)
}

func (c *UserController) Me(ctx *gin.Context) {
	value, exists := ctx.Get("user_id")
	if !exists {
		response.Fail(ctx, 401, "unauthorized")
		return
	}

	userID, ok := value.(int64)
	if !ok {
		response.Fail(ctx, 500, "invalid user id in context")
		return
	}

	user, err := c.s.GetMe(userID)
	if err != nil {
		response.Fail(ctx, 400, err.Error())
		return
	}

	response.Success(ctx, user)
}

func GetUserID(c *gin.Context) (int64, error) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, strconv.ErrSyntax
	}

	userID, ok := value.(int64)
	if !ok {
		return 0, strconv.ErrSyntax
	}

	return userID, nil
}
