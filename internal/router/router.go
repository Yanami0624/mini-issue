package router

import (
	"github.com/gin-gonic/gin"

	"mini-issue/internal/controller"
	"mini-issue/internal/middleware"
)

func NewRouter(ctl *controller.UserController) *gin.Engine {
	r := gin.Default()

	r.POST("/register", ctl.Register)
	r.POST("/login", ctl.Login)

	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/me", ctl.Me)
	}

	return r
}
