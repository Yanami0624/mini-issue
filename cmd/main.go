package main

import (
	"log"

	cache "mini-issue/internal/cache"
	"mini-issue/internal/controller"
	"mini-issue/internal/dao"
	"mini-issue/internal/router"
	"mini-issue/internal/service"
	"mini-issue/internal/worker"
	rediscache "mini-issue/pkg/cache"
	"mini-issue/pkg/db"
)

func main() {
	mysqlDB := db.NewMySQL()
	defer mysqlDB.Close()

	redisClient, _ := rediscache.NewRedisClient()
	defer redisClient.Close()

	userDAO := dao.NewUserDAO(mysqlDB)
	userService := service.NewUserService(userDAO)
	userController := controller.NewUserController(userService)

	issueDAO := dao.NewIssueDAO(mysqlDB)
	issueCache := cache.NewIssueCache(redisClient)
	notificationProducer := worker.NewNotificationProducer(redisClient)
	issueService := service.NewIssueService(
		issueDAO,
		issueCache,
		notificationProducer,
	)
	issueController := controller.NewIssueController(issueService)

	r := router.NewRouter(userController, issueController)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
