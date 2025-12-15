package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"ketukApps/internal/models"
	"ketukApps/internal/services"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupUnblockingHandlerTest(t *testing.T) (*UnblockingHandler, sqlmock.Sqlmock, *gin.Engine) {
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

	service := services.NewUnblockingService(gormDB)
	handler := NewUnblockingHandler(service)
	r := gin.Default()

	return handler, mock, r
}

func TestUnblockingHandler_CreateUnblocking(t *testing.T) {
	handler, mock, r := setupUnblockingHandlerTest(t)
	r.POST("/api/unblockings/v1", handler.CreateUnblocking)

	t.Run("Success", func(t *testing.T) {
		body := models.CreateUnblockingRequest{
			Tahun:     2023,
			Semester:  "Ganjil",
			StartDate: time.Now(),
			EndDate:   time.Now().Add(24 * time.Hour),
			UserID:    1,
		}
		jsonBody, _ := json.Marshal(body)

		mock.ExpectBegin()
		// Mock Collision check
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "unblocking"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/api/unblockings/v1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/unblockings/v1", bytes.NewBufferString("bad"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		body := models.CreateUnblockingRequest{Tahun: 2023, UserID: 1}
		jsonBody, _ := json.Marshal(body)

		mock.ExpectBegin()
		// Mock Collision check
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "unblocking"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodPost, "/api/unblockings/v1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUnblockingHandler_GetUnblockingByID(t *testing.T) {
	handler, mock, r := setupUnblockingHandlerTest(t)
	r.GET("/api/unblockings/v1/:id", handler.GetUnblockingByID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking" WHERE "unblocking"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking"`)).
			WithArgs(1, 1).
			WillReturnError(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUnblockingHandler_GetAllUnblockings(t *testing.T) {
	handler, mock, r := setupUnblockingHandlerTest(t)
	r.GET("/api/unblockings/v1", handler.GetAllUnblockings)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestUnblockingHandler_GetUnblockingsByUserID(t *testing.T) {
	handler, mock, r := setupUnblockingHandlerTest(t)
	r.GET("/api/unblockings/v1/user/:user_id", handler.GetUnblockingsByUserID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking" WHERE user_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1/user/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/unblockings/v1/user/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
