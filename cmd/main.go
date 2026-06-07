package main

import (
	"log"
	. "mini-issue/internal/dao"
	. "mini-issue/pkg/db"
)

func main() {
	db := Link2DB()
	if db == nil {
		log.Fatal("connect mysql failed")
	}
	defer db.Close()

	user := NewUserDAO(db)
	if err := user.CreateUser("Alice", "1234"); err != nil {
		log.Fatal(err)
	}
}
