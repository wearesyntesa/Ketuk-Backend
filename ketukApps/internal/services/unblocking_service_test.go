package services

import (
	"regexp"
	"testing"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUnblockingService_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUnblockingService(db)

	t.Run("Success", func(t *testing.T) {
		ub := &models.Unblocking{
			Tahun:  2023,
			UserID: 1,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "unblocking" ("tahun","semester","created_at","start_date","user_id","end_date") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`)).
			WithArgs(ub.Tahun, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), ub.UserID, sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		_, err := service.Create(ub)
		assert.NoError(t, err)
	})
}

func TestUnblockingService_GetByID(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUnblockingService(db)

	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking" WHERE "unblocking"."id" = $1 ORDER BY "unblocking"."id" LIMIT $2`)).
			WithArgs(1, 1). // id, limit
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		res, err := service.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestUnblockingService_GetDueUnblockings(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUnblockingService(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "unblocking" WHERE scheduled_at <= NOW()`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		res, err := service.GetDueUnblockings()
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
