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

func setupTicketHandlerTest(t *testing.T) (*TicketHandler, sqlmock.Sqlmock, *gin.Engine) {
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

	service := services.NewTicketService(gormDB)
	handler := NewTicketHandler(service)
	r := gin.Default()

	return handler, mock, r
}

func TestTicketHandler_GetAllTickets(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.GET("/api/tickets/v1", handler.GetAllTickets)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Ticket 1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(0).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnError(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTicketHandler_GetTicketByID(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.GET("/api/tickets/v1/:id", handler.GetTicketByID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Ticket 1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTicketHandler_GetTicketsByUserID(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.GET("/api/tickets/v1/user/:user_id", handler.GetTicketsByUserID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE user_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Ticket 1"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/user/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE user_id = $1`)).
			WithArgs(1).
			WillReturnError(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/user/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTicketHandler_GetPendingTickets(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.GET("/api/tickets/v1/pending", handler.GetPendingTickets)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE status = $1`)).
			WithArgs("pending").
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/pending", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTicketHandler_CreateTicket(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.POST("/api/tickets/v1", handler.CreateTicket)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.CreateTicketRequest{UserID: 1, Title: "New"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		// Mock User existence check
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPost, "/api/tickets/v1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/tickets/v1", bytes.NewBufferString("bad"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTicketHandler_UpdateTicket(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.PUT("/api/tickets/v1/:id", handler.UpdateTicket)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.UpdateTicketRequest{Title: "Edited"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPut, "/api/tickets/v1/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		reqBody := models.UpdateTicketRequest{Title: "Edited"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodPut, "/api/tickets/v1/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestTicketHandler_UpdateTicketStatus(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.PATCH("/api/tickets/v1/:id/status", handler.UpdateTicketStatus)

	t.Run("Success", func(t *testing.T) {
		reqBody := UpdateStatusRequest{Status: "approved"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPatch, "/api/tickets/v1/1/status", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTicketHandler_ApproveRejectTicket(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.PATCH("/api/tickets/v1/:id/approve", handler.ApproveTicket)
	r.PATCH("/api/tickets/v1/:id/reject", handler.RejectTicket)

	t.Run("Approve", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPatch, "/api/tickets/v1/1/approve", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Reject", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPatch, "/api/tickets/v1/1/reject", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTicketHandler_BulkUpdateStatus(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.POST("/api/tickets/v1/bulk-status", handler.BulkUpdateStatus)

	t.Run("Success", func(t *testing.T) {
		reqBody := BulkUpdateStatusRequest{IDs: []int{1, 2}, Status: "completed"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPost, "/api/tickets/v1/bulk-status", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("EmptyIDs", func(t *testing.T) {
		reqBody := BulkUpdateStatusRequest{IDs: []int{}, Status: "completed"}
		jsonBody, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/api/tickets/v1/bulk-status", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestTicketHandler_DeleteTicket(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.DELETE("/api/tickets/v1/:id", handler.DeleteTicket)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		// Mock User check
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodDelete, "/api/tickets/v1/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTicketHandler_GetStatistics(t *testing.T) {
	handler, mock, r := setupTicketHandlerTest(t)
	r.GET("/api/tickets/v1/statistics", handler.GetStatistics)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT status, count(*) FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"status", "count"}).AddRow("pending", 5))

		req, _ := http.NewRequest(http.MethodGet, "/api/tickets/v1/statistics", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
