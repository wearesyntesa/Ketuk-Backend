package services

import (
	"regexp"
	"testing"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTicketService_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		title := "Test Ticket"
		desc := "Description"

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets" ("user_id","title","description","status","id_schedule","created_at","updated_at","approved_at","reason") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
			WithArgs(userID, title, desc, "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Reload with User
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "status"}).AddRow(1, userID, title, "pending"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "full_name"}).AddRow(userID, "Test User"))

		ticket, err := service.Create(userID, title, desc)
		assert.NoError(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, title, ticket.Title)
		assert.Equal(t, models.StatusPending, ticket.Status)
	})

	t.Run("EmptyTitle", func(t *testing.T) {
		_, err := service.Create(1, "", "desc")
		assert.Error(t, err)
		assert.Equal(t, "title is required", err.Error())
	})
}

func TestTicketService_GetByID(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title"}).AddRow(1, 1, "Ticket 1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "full_name"}).AddRow(1, "User 1"))

		ticket, err := service.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, uint(1), ticket.ID)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := service.GetByID(999)
		assert.Error(t, err)
		assert.Equal(t, "ticket not found", err.Error())
	})
}

func TestTicketService_UpdateStatus(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("SuccessAccepted", func(t *testing.T) {
		// 1. Check exists
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "pending"))

		// 2. Update
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets" SET "approved_at"=$1,"reason"=$2,"status"=$3,"updated_at"=$4 WHERE "id" = $5`)).
			WithArgs(sqlmock.AnyArg(), "Reason", "accepted", sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 3. Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status", "user_id"}).AddRow(1, "accepted", 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		ticket, err := service.UpdateStatus(1, "accepted", "Reason")
		assert.NoError(t, err)
		assert.Equal(t, "accepted", string(ticket.Status))
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		_, err := service.UpdateStatus(1, "invalid", "")
		assert.Error(t, err)
	})
}

func TestTicketService_GetByUserID(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE user_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1).AddRow(2, 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		tickets, err := service.GetByUserID(1)
		assert.NoError(t, err)
		assert.Len(t, tickets, 2)
	})
}

func TestTicketService_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.Delete(1)
		assert.NoError(t, err)
	})
}

func TestTicketService_GetStatistics(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		statuses := []string{"pending", "accepted", "rejected"}
		for _, status := range statuses {
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "tickets" WHERE status = $1`)).
				WithArgs(status).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		}

		stats, err := service.GetStatistics()
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(10), stats["total"])
	})
}

func TestTicketService_BulkUpdateStatus(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewTicketService(db)

	t.Run("Success", func(t *testing.T) {
		ids := []uint{1}

		for _, id := range ids {
			// Update logic repeated for each ID
			// 1. Check exists
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
				WithArgs(id, 1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(id, "pending"))

			// 2. Update
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets" SET "reason"=$1,"status"=$2,"updated_at"=$3 WHERE "id" = $4`)).
				WithArgs("Bulk Reason", "rejected", sqlmock.AnyArg(), id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()

			// 3. Reload
			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1 ORDER BY "tickets"."id" LIMIT $2`)).
				WithArgs(id, 1).
				WillReturnRows(sqlmock.NewRows([]string{"id", "status", "user_id"}).AddRow(id, "rejected", 1))

			mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		}

		tickets, err := service.BulkUpdateStatus(ids, "rejected", "Bulk Reason")
		assert.NoError(t, err)
		assert.Len(t, tickets, 1)
	})
}
