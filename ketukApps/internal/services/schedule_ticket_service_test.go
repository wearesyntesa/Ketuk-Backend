package services

import (
	"regexp"
	"testing"
	"time"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestScheduleService_CreateScheduleTicket(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewScheduleService(db)

	t.Run("Success", func(t *testing.T) {
		schedule := &models.ScheduleTicket{
			Title:     "Ticket Schedule",
			UserID:    1,
			Kategori:  models.Kelas,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(2 * time.Hour),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "schedule_ticket" ("title","start_date","end_date","user_id","kategori","description","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id_schedule"`)).
			WithArgs(schedule.Title, sqlmock.AnyArg(), sqlmock.AnyArg(), schedule.UserID, schedule.Kategori, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))
		mock.ExpectCommit()

		// Reload with User and Tickets
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1 ORDER BY "schedule_ticket"."id_schedule" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule", "title"}).AddRow(1, "Ticket Schedule"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(0).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id_schedule" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

		created, err := service.CreateScheduleTicket(schedule)
		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, "Ticket Schedule", created.Title)
	})

	t.Run("EmptyTitle", func(t *testing.T) {
		s := &models.ScheduleTicket{UserID: 1}
		_, err := service.CreateScheduleTicket(s)
		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())
	})
}

func TestScheduleService_UpdateScheduleTicket(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewScheduleService(db)

	t.Run("Success", func(t *testing.T) {
		// Check exists
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1 ORDER BY "schedule_ticket"."id_schedule" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))

		// Update
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "schedule_ticket" SET "title"=$1,"updated_at"=$2 WHERE "id_schedule" = $3`)).
			WithArgs("Updated", sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Reload (Select + User + Tickets)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "schedule_ticket" WHERE "schedule_ticket"."id_schedule" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id_schedule"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		updates := map[string]interface{}{"title": "Updated"}
		_, err := service.UpdateScheduleTicket(1, updates)
		assert.NoError(t, err)
	})
}

func TestScheduleService_CreateUnblocking(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewScheduleService(db)

	t.Run("Success", func(t *testing.T) {
		ub := &models.Unblocking{
			Tahun:    2023,
			UserID:   1,
			Semester: models.SemesterGanjil,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "unblocking" ("tahun","semester","created_at","start_date","user_id","end_date") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(ub.Tahun, ub.Semester, sqlmock.AnyArg(), sqlmock.AnyArg(), ub.UserID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking" WHERE "unblocking"."id" = $1 ORDER BY "unblocking"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		created, err := service.CreateUnblocking(ub)
		assert.NoError(t, err)
		assert.NotNil(t, created)
	})

	t.Run("EmptyTahun", func(t *testing.T) {
		ub := &models.Unblocking{UserID: 1}
		_, err := service.CreateUnblocking(ub)
		assert.Error(t, err)
		assert.Equal(t, "tahun is required", err.Error())
	})
}
