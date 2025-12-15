package services

import (
	"context"
	"regexp"
	"testing"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestEnhancedTicketService_CreateWithEvents(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewEnhancedTicketService(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		title := "Enhanced Ticket"
		desc := "Description"

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "tickets" ("user_id","title","description","status","id_schedule","created_at","updated_at","approved_at","reason") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
			WithArgs(userID, title, desc, "pending", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, title))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "full_name"}))

		ticket, err := service.CreateWithEvents(ctx, userID, title, desc)
		assert.NoError(t, err)
		assert.NotNil(t, ticket)
		assert.Equal(t, title, ticket.Title)
	})
}

func TestEnhancedTicketService_UpdateStatusWithEvents(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewEnhancedTicketService(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// GetByID (Current Ticket)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "pending"))

		// UpdateStatus (TicketService)
		// 1. GetByID (Check existence)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "pending"))
		// 2. Update
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		// 3. Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "accepted"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		ticket, err := service.UpdateStatusWithEvents(ctx, 1, "accepted", "Reason")
		assert.NoError(t, err)
		assert.Equal(t, "accepted", string(ticket.Status))
	})
}

func TestEnhancedTicketService_UpdateWithEvents(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewEnhancedTicketService(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// GetByID (Current Ticket)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Old Title"))

		// Update (TicketService)
		// 1. GetByID
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Old Title"))
		// 2. Update
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tickets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
		// 3. Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "New Title"))
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		req := models.UpdateTicketRequest{Title: "New Title"}
		ticket, err := service.UpdateWithEvents(ctx, 1, req)
		assert.NoError(t, err)
		assert.Equal(t, "New Title", ticket.Title)
	})
}

func TestEnhancedTicketService_DeleteWithEvents(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewEnhancedTicketService(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// GetByID
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "ToDelete"))

		// Delete (TicketService)
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "tickets" WHERE "tickets"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeleteWithEvents(ctx, 1)
		assert.NoError(t, err)
	})
}
