package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Link2DB() *sqlx.DB {
	db, err := sqlx.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		fmt.Println("Open mysql failed", err)
		return nil
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Ping mysql failed", err)
		db.Close()
		return nil
	}
	db.SetMaxOpenConns(40)
	db.SetMaxIdleConns(20)
	return db
}
