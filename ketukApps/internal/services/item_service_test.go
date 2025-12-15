package services

import (
	"regexp"
	"testing"

	"ketukApps/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestItemService_CreateItemCategory(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Success", func(t *testing.T) {
		cat := &models.ItemCategory{
			CategoryName: "Electronics",
			Quantity:     10,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items_category" ("category_name","specification","quantity","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(cat.CategoryName, sqlmock.AnyArg(), cat.Quantity, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		// Reload
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1 ORDER BY "items_category"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_name"}).AddRow(1, "Electronics"))

		// Preload Items
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items" WHERE "items"."category_id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		createdCat, err := service.CreateItemCategory(cat)
		assert.NoError(t, err)
		assert.NotNil(t, createdCat)
		assert.Equal(t, "Electronics", createdCat.CategoryName)
	})

	t.Run("EmptyName", func(t *testing.T) {
		cat := &models.ItemCategory{Quantity: 5}
		_, err := service.CreateItemCategory(cat)
		assert.Error(t, err)
		assert.Equal(t, "category name is required", err.Error())
	})
}

func TestItemService_GetItemCategoryByID(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1 ORDER BY "items_category"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_name"}).AddRow(1, "Test Cat"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items" WHERE "items"."category_id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "category_id"}).AddRow(100, "Item 1", 1))

		cat, err := service.GetItemCategoryByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, cat)
		assert.Equal(t, "Test Cat", cat.CategoryName)
		assert.Len(t, cat.Items, 1)
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1 ORDER BY "items_category"."id" LIMIT $2`)).
			WithArgs(999, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		cat, err := service.GetItemCategoryByID(999)
		assert.Error(t, err)
		assert.Equal(t, "item category not found", err.Error())
		assert.Nil(t, cat)
	})
}

func TestItemService_CreateItem(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Success", func(t *testing.T) {
		item := &models.Item{
			Name:       "Laptop",
			CategoryID: 1,
			Kondisi:    "Good",
		}

		// Verify category exists
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1 ORDER BY "items_category"."id" LIMIT $2`)).
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// Create Item
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "items" ("name","year","kondisi","note","category_id","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
			WithArgs(item.Name, sqlmock.AnyArg(), item.Kondisi, sqlmock.AnyArg(), item.CategoryID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(101))
		mock.ExpectCommit()

		// Reload Item with Category
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items" WHERE "items"."id" = $1 ORDER BY "items"."id" LIMIT $2`)).
			WithArgs(101, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "category_id"}).AddRow(101, "Laptop", 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "category_name"}).AddRow(1, "Electronics"))

		createdItem, err := service.CreateItem(item)
		assert.NoError(t, err)
		assert.NotNil(t, createdItem)
		assert.Equal(t, "Laptop", createdItem.Name)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		item := &models.Item{Name: "Laptop", CategoryID: 99}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "items_category" WHERE "items_category"."id" = $1 ORDER BY "items_category"."id" LIMIT $2`)).
			WithArgs(99, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		_, err := service.CreateItem(item)
		assert.Error(t, err)
		assert.Equal(t, "category not found", err.Error())
	})
}

func TestItemService_DeleteItem(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "items" WHERE "items"."id" = $1`)).
			WithArgs(101).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeleteItem(101)
		assert.NoError(t, err)
	})
}

func TestItemService_DeleteItemCategory(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Success", func(t *testing.T) {
		// Check for existing items
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "items" WHERE category_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "items_category" WHERE "items_category"."id" = $1`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := service.DeleteItemCategory(1)
		assert.NoError(t, err)
	})

	t.Run("HasItems", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "items" WHERE category_id = $1`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))

		err := service.DeleteItemCategory(1)
		assert.Error(t, err)
		assert.Equal(t, "cannot delete category with existing items", err.Error())
	})
}

func TestItemService_GetItemStatistics(t *testing.T) {
	db, mock := setupTestDB(t)
	service := NewItemService(db)

	t.Run("Success", func(t *testing.T) {
		// Count Items
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "items"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

		// Count Categories
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "items_category"`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

		// Group Stats
		rows := sqlmock.NewRows([]string{"category_id", "category_name", "item_count"}).
			AddRow(1, "Cat A", 50).
			AddRow(2, "Cat B", 50)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT items.category_id, items_category.category_name, COUNT(items.id) as item_count FROM "items" LEFT JOIN items_category ON items.category_id = items_category.id GROUP BY items.category_id, items_category.category_name`)).
			WillReturnRows(rows)

		stats, err := service.GetItemStatistics()
		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, int64(100), stats["total_items"])
		assert.Equal(t, int64(10), stats["total_categories"])
	})
}
