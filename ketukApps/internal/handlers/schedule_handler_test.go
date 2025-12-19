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

func setupScheduleHandlerTest(t *testing.T) (*ScheduleHandler, sqlmock.Sqlmock, *gin.Engine) {
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

	service := services.NewScheduleService(gormDB)
	handler := NewScheduleHandler(service)
	r := gin.Default()

	return handler, mock, r
}

// --- Schedule Ticket Tests ---

func TestScheduleHandler_GetAllScheduleTickets(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.GET("/api/schedules/tickets", handler.GetAllScheduleTickets)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Sched 1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id_schedule" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "id_schedule"}).AddRow(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(0).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/tickets", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket"`)).
			WillReturnError(errors.New("db error"))

		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/tickets", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestScheduleHandler_GetScheduleTicketByID(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.GET("/api/schedules/tickets/:id", handler.GetScheduleTicketByID)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Old"))

		// Preloads
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/tickets/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/tickets/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/tickets/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestScheduleHandler_CreateScheduleTicket(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.POST("/api/schedules/tickets", handler.CreateScheduleTicket)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.ScheduleTicket{Title: "New", UserID: 1}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedule_ticket"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/api/schedules/tickets", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("BadRequest", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/schedules/tickets", bytes.NewBufferString("bad"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		reqBody := models.ScheduleTicket{Title: "New", UserID: 1}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedule_ticket"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodPost, "/api/schedules/tickets", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestScheduleHandler_UpdateScheduleTicket(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.PUT("/api/schedules/tickets/:id", handler.UpdateScheduleTicket)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{"title": "Updated"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Old"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schedule_ticket"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Updated"))
		// Preloads...
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPut, "/api/schedules/tickets/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		reqBody := map[string]interface{}{"title": "Updated"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		req, _ := http.NewRequest(http.MethodPut, "/api/schedules/tickets/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/api/schedules/tickets/abc", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestScheduleHandler_DeleteScheduleTicket(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.DELETE("/api/schedules/tickets/:id", handler.DeleteScheduleTicket)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedule_ticket"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodDelete, "/api/schedules/tickets/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("DBError", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedule_ticket"`)).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		req, _ := http.NewRequest(http.MethodDelete, "/api/schedules/tickets/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// --- Schedule Reguler Tests ---

func TestScheduleHandler_GetAllScheduleReguler(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.GET("/api/schedules/reguler", handler.GetAllScheduleReguler)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_reguler"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodGet, "/api/schedules/reguler", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestScheduleHandler_CreateScheduleReguler(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.POST("/api/schedules/reguler", handler.CreateScheduleReguler)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.ScheduleReguler{Title: "Meeting", UserID: 1}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedule_reguler"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/api/schedules/reguler", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestScheduleHandler_UpdateScheduleReguler(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.PUT("/api/schedules/reguler/:id", handler.UpdateScheduleReguler)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{"title": "Updated Meeting"}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_reguler"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Meeting"))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schedule_reguler"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_reguler"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Updated Meeting"))
		// Preload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req, _ := http.NewRequest(http.MethodPut, "/api/schedules/reguler/1", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestScheduleHandler_DeleteScheduleReguler(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.DELETE("/api/schedules/reguler/:id", handler.DeleteScheduleReguler)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "schedule_reguler"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodDelete, "/api/schedules/reguler/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// --- Unblocking Tests (ScheduleHandler Implementation) ---

func TestScheduleHandler_GetAllUnblocking(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.GET("/api/unblocking", handler.GetAllUnblocking)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		req, _ := http.NewRequest(http.MethodGet, "/api/unblocking", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestScheduleHandler_CreateUnblocking(t *testing.T) {
	handler, mock, r := setupScheduleHandlerTest(t)
	r.POST("/api/unblocking", handler.CreateUnblocking)

	t.Run("Success", func(t *testing.T) {
		reqBody := models.Unblocking{Tahun: 2023}
		jsonBody, _ := json.Marshal(reqBody)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "unblocking"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		req, _ := http.NewRequest(http.MethodPost, "/api/unblocking", bytes.NewBuffer(jsonBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})
}
