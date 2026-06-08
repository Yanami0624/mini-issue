package dao

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestUserDAOCreateUser(t *testing.T) {
	udao, mock, closeDB := newMockUserDAO(t)
	defer closeDB()

	mock.ExpectExec("insert into `user` \\(username, password, created_at\\) values \\(\\?, \\?, \\?\\)").
		WithArgs("alice", "hashed-password", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := udao.CreateUser("alice", "hashed-password"); err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
	assertExpectationsMet(t, mock)
}

func TestUserDAOGetByUsernameReturnsUser(t *testing.T) {
	udao, mock, closeDB := newMockUserDAO(t)
	defer closeDB()

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at"}).
		AddRow(int64(1), "alice", "hashed-password", createdAt)

	mock.ExpectQuery("select id, username, password, created_at\\s+from `user`\\s+where username = \\?").
		WithArgs("alice").
		WillReturnRows(rows)

	user, err := udao.GetByUsername("alice")
	if err != nil {
		t.Fatalf("GetByUsername() error = %v", err)
	}
	if user == nil {
		t.Fatal("GetByUsername() user = nil, want user")
	}
	if user.ID != 1 || user.Username != "alice" || user.Password != "hashed-password" {
		t.Fatalf("GetByUsername() user = %+v, want alice with id 1", user)
	}
	assertExpectationsMet(t, mock)
}

func TestUserDAOGetByUsernameReturnsNilWhenNotFound(t *testing.T) {
	udao, mock, closeDB := newMockUserDAO(t)
	defer closeDB()

	mock.ExpectQuery("select id, username, password, created_at\\s+from `user`\\s+where username = \\?").
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	user, err := udao.GetByUsername("missing")
	if err != nil {
		t.Fatalf("GetByUsername() error = %v", err)
	}
	if user != nil {
		t.Fatalf("GetByUsername() user = %+v, want nil", user)
	}
	assertExpectationsMet(t, mock)
}

func TestUserDAOGetByUserIDReturnsUser(t *testing.T) {
	udao, mock, closeDB := newMockUserDAO(t)
	defer closeDB()

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "username", "password", "created_at"}).
		AddRow(int64(2), "bob", "hashed-password", createdAt)

	mock.ExpectQuery("select id, username, password, created_at\\s+from user\\s+where id = \\?").
		WithArgs(int64(2)).
		WillReturnRows(rows)

	user, err := udao.GetByUserID(2)
	if err != nil {
		t.Fatalf("GetByUserID() error = %v", err)
	}
	if user == nil {
		t.Fatal("GetByUserID() user = nil, want user")
	}
	if user.ID != 2 || user.Username != "bob" {
		t.Fatalf("GetByUserID() user = %+v, want bob with id 2", user)
	}
	assertExpectationsMet(t, mock)
}

func newMockUserDAO(t *testing.T) (*UserDAO, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return NewUserDAO(sqlxDB), mock, func() {
		_ = sqlxDB.Close()
	}
}

func assertExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
