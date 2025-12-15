package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"ketukApps/internal/models"
	"ketukApps/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUserHandlerTest(t *testing.T) (*UserHandler, sqlmock.Sqlmock, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	service := services.NewUserService(gormDB)
	handler := NewUserHandler(service)
	r := gin.Default()

	return handler, mock, r
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	handler, mock, r := setupUserHandlerTest(t)
	r.GET("/api/users/v1", handler.GetAllUsers)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test@example.com"))

		req, _ := http.NewRequest(http.MethodGet, "/api/users/v1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnError(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/users/v1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserHandler_GetUserByID(t *testing.T) {
	handler, mock, r := setupUserHandlerTest(t)
	r.GET("/api/users/v1/:id", handler.GetUserByID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).AddRow(1, "test@example.com"))

		req, _ := http.NewRequest(http.MethodGet, "/api/users/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/users/v1/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/users/v1/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	handler, mock, r := setupUserHandlerTest(t)
	r.POST("/api/users/v1", handler.CreateUser)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreateUserRequest{Email: "new@example.com", Name: "New User"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		// Mock User existence check
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("new@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/api/users/v1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/users/v1", bytes.NewBufferString("bad"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		reqBody := models.CreateUserRequest{Email: "new@example.com", Name: "New User"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs("new@example.com").
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodPost, "/api/users/v1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	handler, mock, r := setupUserHandlerTest(t)
	r.PUT("/api/users/v1/:id", handler.UpdateUser)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.UpdateUserRequest{Name: "Updated"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1, 1). // id, limit
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Old"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Updated"))

		req, _ := http.NewRequest(http.MethodPut, "/api/users/v1/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		reqBody := models.UpdateUserRequest{Name: "Updated"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodPut, "/api/users/v1/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	handler, mock, r := setupUserHandlerTest(t)
	r.DELETE("/api/users/v1/:id", handler.DeleteUser)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodDelete, "/api/users/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodDelete, "/api/users/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
