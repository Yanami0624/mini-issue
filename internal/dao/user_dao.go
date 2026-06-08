package dao

import (
	"database/sql"
	"fmt"
	"mini-issue/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserDAO struct {
	db *sqlx.DB
}

func NewUserDAO(db *sqlx.DB) *UserDAO {
	udao := new(UserDAO)
	udao.db = db
	return udao
}

func (udao *UserDAO) CreateUser(username, hashedPassword string) error {
	timestamp := time.Now()
	_, err := udao.db.Exec("insert into `user` (username, password, created_at) values (?, ?, ?)", username, hashedPassword, timestamp)
	if err != nil {
		fmt.Println("failed: CreateUser()", err)
	}
	return err
}

func (udao *UserDAO) GetByUsername(username string) (*model.User, error) {
	var user model.User
	query := `
		select id, username, password, created_at
		from ` + "`user`" + `
		where username = ?
	`

	err := udao.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println("failed: GetByUsername()", err)
		return nil, err
	}

	return &user, nil
}

func (udao *UserDAO) GetByUserID(userid int64) (*model.User, error) {
	var user model.User
	query := `
		select id, username, password, created_at
		from user
		where id = ?
	`
	err := udao.db.Get(&user, query, userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		fmt.Println("failed: GetByUsername()", err)
		return nil, err
	}

	return &user, nil
}
