package service

import (
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"mini-issue/internal/dao"
	"mini-issue/internal/model"
	jwtpkg "mini-issue/pkg/jwt"
)

func TestUserServiceRegisterRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		req     model.RegisterRequest
		wantErr string
	}{
		{
			name:    "empty username",
			req:     model.RegisterRequest{Username: "", Password: "123456"},
			wantErr: "username can not be empty",
		},
		{
			name:    "short password",
			req:     model.RegisterRequest{Username: "alice", Password: "123"},
			wantErr: "password should have 6 characters at least",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us, mock, closeDB := newMockUserService(t)
			defer closeDB()

			expectGetByUsernameNoRows(mock, tt.req.Username)

			err := us.Register(tt.req)
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("Register() error = %v, want %q", err, tt.wantErr)
			}
			assertServiceExpectationsMet(t, mock)
		})
	}
}

func TestUserServiceRegisterRejectsExistingUser(t *testing.T) {
	us, mock, closeDB := newMockUserService(t)
	defer closeDB()

	rows := userRows().AddRow(int64(1), "alice", "hashed-password", time.Now())
	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs("alice").
		WillReturnRows(rows)

	err := us.Register(model.RegisterRequest{Username: "alice", Password: "123456"})
	if err == nil || err.Error() != "user already exist" {
		t.Fatalf("Register() error = %v, want user already exist", err)
	}
	assertServiceExpectationsMet(t, mock)
}

func TestUserServiceRegisterCreatesUserWithHashedPassword(t *testing.T) {
	us, mock, closeDB := newMockUserService(t)
	defer closeDB()

	expectGetByUsernameNoRows(mock, "alice")
	mock.ExpectExec("insert into `user` \\(username, password, created_at\\) values \\(\\?, \\?, \\?\\)").
		WithArgs("alice", bcryptHashOf("123456"), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := us.Register(model.RegisterRequest{Username: "alice", Password: "123456"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	assertServiceExpectationsMet(t, mock)
}

func TestUserServiceLoginReturnsToken(t *testing.T) {
	us, mock, closeDB := newMockUserService(t)
	defer closeDB()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	rows := userRows().AddRow(int64(7), "alice", string(hashedPassword), time.Now())
	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs("alice").
		WillReturnRows(rows)

	resp, err := us.Login(model.LoginRequest{Username: "alice", Password: "123456"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if resp == nil || resp.Token == "" {
		t.Fatalf("Login() response = %+v, want token", resp)
	}

	claims, err := jwtpkg.ParseToken(resp.Token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 7 || claims.Username != "alice" {
		t.Fatalf("claims = %+v, want user 7 alice", claims)
	}
	assertServiceExpectationsMet(t, mock)
}

func TestUserServiceLoginRejectsIncorrectPassword(t *testing.T) {
	us, mock, closeDB := newMockUserService(t)
	defer closeDB()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("right-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	rows := userRows().AddRow(int64(7), "alice", string(hashedPassword), time.Now())
	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs("alice").
		WillReturnRows(rows)

	resp, err := us.Login(model.LoginRequest{Username: "alice", Password: "wrong-password"})
	if err == nil || err.Error() != "incorrect password" {
		t.Fatalf("Login() error = %v, want incorrect password", err)
	}
	if resp != nil {
		t.Fatalf("Login() response = %+v, want nil", resp)
	}
	assertServiceExpectationsMet(t, mock)
}

func TestUserServiceGetMeReturnsUser(t *testing.T) {
	us, mock, closeDB := newMockUserService(t)
	defer closeDB()

	rows := userRows().AddRow(int64(8), "bob", "hashed-password", time.Now())
	mock.ExpectQuery("select id, username, password, created_at\\s+from user\\s+where id = \\?").
		WithArgs(int64(8)).
		WillReturnRows(rows)

	user, err := us.GetMe(8)
	if err != nil {
		t.Fatalf("GetMe() error = %v", err)
	}
	if user == nil || user.ID != 8 || user.Username != "bob" {
		t.Fatalf("GetMe() user = %+v, want bob with id 8", user)
	}
	assertServiceExpectationsMet(t, mock)
}

type bcryptHashOf string

func (m bcryptHashOf) Match(value driver.Value) bool {
	hash, ok := value.(string)
	if !ok {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(string(m))) == nil
}

func newMockUserService(t *testing.T) (*UserService, sqlmock.Sqlmock, func()) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return NewUserService(dao.NewUserDAO(sqlxDB)), mock, func() {
		_ = sqlxDB.Close()
	}
}

func expectGetByUsernameNoRows(mock sqlmock.Sqlmock, username string) {
	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)
}

func usernameQueryPattern() string {
	return "select id, username, password, created_at\\s+from `user`\\s+where username = \\?"
}

func userRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "password", "created_at"})
}

func assertServiceExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
