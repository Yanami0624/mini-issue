package router

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"mini-issue/internal/controller"
	"mini-issue/internal/dao"
	"mini-issue/internal/service"
	jwtpkg "mini-issue/pkg/jwt"
)

func TestRouterRegisterSuccess(t *testing.T) {
	router, mock, closeDB := newTestRouter(t)
	defer closeDB()

	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs("alice").
		WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("insert into `user` \\(username, password, created_at\\) values \\(\\?, \\?, \\?\\)").
		WithArgs("alice", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/register", jsonBody(`{"username":"alice","password":"123456"}`))
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	assertRouterResponseCode(t, recorder.Body.Bytes(), 0)
	assertRouterExpectationsMet(t, mock)
}

func TestRouterLoginSuccess(t *testing.T) {
	router, mock, closeDB := newTestRouter(t)
	defer closeDB()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	rows := userRows().AddRow(int64(3), "alice", string(hashedPassword), time.Now())
	mock.ExpectQuery(usernameQueryPattern()).
		WithArgs("alice").
		WillReturnRows(rows)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", jsonBody(`{"username":"alice","password":"123456"}`))
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var body struct {
		Code int `json:"code"`
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body.Code != 0 {
		t.Fatalf("code = %d, want 0", body.Code)
	}
	if body.Data.Token == "" {
		t.Fatal("token is empty")
	}
	claims, err := jwtpkg.ParseToken(body.Data.Token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 3 || claims.Username != "alice" {
		t.Fatalf("claims = %+v, want user 3 alice", claims)
	}
	assertRouterExpectationsMet(t, mock)
}

func TestRouterMeRequiresAuth(t *testing.T) {
	router, _, closeDB := newTestRouter(t)
	defer closeDB()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
}

func TestRouterMeSuccess(t *testing.T) {
	router, mock, closeDB := newTestRouter(t)
	defer closeDB()

	tokenString, err := jwtpkg.GenerateToken(4, "bob")
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}
	rows := userRows().AddRow(int64(4), "bob", "hashed-password", time.Now())
	mock.ExpectQuery("select id, username, password, created_at\\s+from user\\s+where id = \\?").
		WithArgs(int64(4)).
		WillReturnRows(rows)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var body struct {
		Code int `json:"code"`
		Data struct {
			ID       int64  `json:"id"`
			Username string `json:"username"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body.Code != 0 {
		t.Fatalf("code = %d, want 0", body.Code)
	}
	if body.Data.ID != 4 || body.Data.Username != "bob" {
		t.Fatalf("data = %+v, want user 4 bob", body.Data)
	}
	assertRouterExpectationsMet(t, mock)
}

func newTestRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, func()) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	userDAO := dao.NewUserDAO(sqlxDB)
	userService := service.NewUserService(userDAO)
	userController := controller.NewUserController(userService)

	issueDAO := dao.NewIssueDAO(sqlxDB)
	issueService := service.NewIssueService(issueDAO)
	issueController := controller.NewIssueController(issueService)

	return NewRouter(userController, issueController), mock, func() {
		_ = sqlxDB.Close()
	}
}

func jsonBody(body string) *bytes.Reader {
	return bytes.NewReader([]byte(body))
}

func assertRouterResponseCode(t *testing.T, raw []byte, want int) {
	t.Helper()

	var body struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if body.Code != want {
		t.Fatalf("code = %d, want %d", body.Code, want)
	}
}

func usernameQueryPattern() string {
	return "select id, username, password, created_at\\s+from `user`\\s+where username = \\?"
}

func userRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "username", "password", "created_at"})
}

func assertRouterExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}
