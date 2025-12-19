package services

import (
	"regexp"
	"testing"
	"time"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialect := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})
	db, err := gorm.Open(dialect, &gorm.Config{})
	assert.NoError(t, err)

	return db, mock
}

func TestUserService_GetByID(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUserService(db)

	t.Run("UserFound", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "full_name"}).
			AddRow(1, "test@example.com", "Test User")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(rows)

		user, err := service.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := service.GetByID(999)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserService_Create(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			Email:    "new@example.com",
			Password: "hashedpassword",
			Name:     "New User",
		}

		// Check overlap email
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs("new@example.com", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Create
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("google_sub","full_name","email","password","role","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
			WithArgs(sqlmock.AnyArg(), user.Name, user.Email, user.Password, "user", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		createdUser, err := service.Create(user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, uint(1), createdUser.ID)
		assert.Equal(t, "user", createdUser.Role) // Default role
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		user := &models.User{
			Email: "existing@example.com",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs("existing@example.com", 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		_, err := service.Create(user)
		assert.Error(t, err)
		assert.Equal(t, "email already exists", err.Error())
	})
}

func TestUserService_Update(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		userID := uint(1)
		req := models.UpdateUserRequest{
			Name: "Updated Name",
		}

		// 1. Check if user exists
		rows := sqlmock.NewRows([]string{"id", "email", "full_name", "created_at", "updated_at"}).
			AddRow(userID, "test@example.com", "Old Name", time.Now(), time.Now())
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		// 2. Update
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "full_name"=$1,"updated_at"=$2 WHERE "id" = $3`)).
			WithArgs("Updated Name", sqlmock.AnyArg(), userID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// 3. Reload (The implementation does a reload)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(userID, 1).
			WillReturnRows(rows)

		updatedUser, err := service.Update(userID, req)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
	})
}

func TestUserService_Delete(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.Delete(1)
		assert.NoError(t, err)
	})
}
