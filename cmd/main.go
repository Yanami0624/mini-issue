package main

import (
	"log"

	"mini-issue/internal/controller"
	"mini-issue/internal/dao"
	"mini-issue/internal/router"
	"mini-issue/internal/service"
	"mini-issue/pkg/db"
)

func main() {
	mysqlDB := db.NewMySQL()
	defer mysqlDB.Close()

	userDAO := dao.NewUserDAO(mysqlDB)
	userService := service.NewUserService(userDAO)
	userController := controller.NewUserController(userService)

	r := router.NewRouter(userController)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
