package router

import (
	"github.com/gin-gonic/gin"

	"mini-issue/internal/controller"
	"mini-issue/internal/middleware"
)

func NewRouter(
	uc *controller.UserController,
	ic *controller.IssueController,
) *gin.Engine {
	r := gin.Default()

	r.POST("/register", uc.Register)
	r.POST("/login", uc.Login)

	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", uc.Me)
	}

	auth.POST("/issues", ic.CreateIssue)
	auth.GET("/issues", ic.ListIssues)
	auth.GET("/issues/:id", ic.GetIssue)
	auth.PATCH("/issues/:id", ic.UpdateIssue)
	auth.DELETE("/issues/:id", ic.DeleteIssue)

	return r
}

