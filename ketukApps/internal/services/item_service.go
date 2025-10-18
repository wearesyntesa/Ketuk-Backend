package services

import (
	"errors"

	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type ItemService struct {
	db *gorm.DB
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{
		db: db,
	}
}

// ItemCategory methods

// GetAllItemCategories returns all item categories
func (s *ItemService) GetAllItemCategories() ([]models.ItemCategory, error) {
	var categories []models.ItemCategory
	result := s.db.Preload("Items").Find(&categories)
	return categories, result.Error
}

// GetItemCategoryByID returns an item category by its ID
func (s *ItemService) GetItemCategoryByID(id int) (*models.ItemCategory, error) {
	var category models.ItemCategory
	result := s.db.Preload("Items").First(&category, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("item category not found")
	}
	return &category, result.Error
}

// CreateItemCategory creates a new item category
func (s *ItemService) CreateItemCategory(category *models.ItemCategory) (*models.ItemCategory, error) {
	if category.CategoryName == "" {
		return nil, errors.New("category name is required")
	}

	// Initialize quantity to 0 if not set
	if category.Quantity == 0 {
		category.Quantity = 0
	}

	result := s.db.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with items data
	s.db.Preload("Items").First(category, category.ID)
	return category, nil
}

// UpdateItemCategory updates an item category
func (s *ItemService) UpdateItemCategory(id int, updates map[string]interface{}) (*models.ItemCategory, error) {
	var category models.ItemCategory
	if err := s.db.First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item category not found")
		}
		return nil, err
	}

	if len(updates) > 0 {
		result := s.db.Model(&category).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("Items").First(&category, id)
	return &category, nil
}

// DeleteItemCategory removes an item category
func (s *ItemService) DeleteItemCategory(id int) error {
	// Check if there are any items in this category
	var count int64
	s.db.Model(&models.Item{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("cannot delete category with existing items")
	}

	result := s.db.Delete(&models.ItemCategory{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item category not found")
	}
	return nil
}

// Item methods

// GetAllItems returns all items
func (s *ItemService) GetAllItems() ([]models.Item, error) {
	var items []models.Item
	result := s.db.Preload("Category").Find(&items)
	return items, result.Error
}

// GetItemByID returns an item by its ID
func (s *ItemService) GetItemByID(id int) (*models.Item, error) {
	var item models.Item
	result := s.db.Preload("Category").First(&item, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("item not found")
	}
	return &item, result.Error
}

// GetItemsByCategoryID returns all items for a specific category
func (s *ItemService) GetItemsByCategoryID(categoryID int) ([]models.Item, error) {
	var items []models.Item
	result := s.db.Preload("Category").Where("category_id = ?", categoryID).Find(&items)
	return items, result.Error
}

// SearchItems searches items by name or note
func (s *ItemService) SearchItems(query string) ([]models.Item, error) {
	var items []models.Item
	searchTerm := "%" + query + "%"
	result := s.db.Preload("Category").
		Where("name ILIKE ? OR note ILIKE ?", searchTerm, searchTerm).
		Find(&items)
	return items, result.Error
}

// GetItemsByKondisi returns all items with a specific condition
func (s *ItemService) GetItemsByKondisi(kondisi string) ([]models.Item, error) {
	var items []models.Item
	result := s.db.Preload("Category").Where("kondisi = ?", kondisi).Find(&items)
	return items, result.Error
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(item *models.Item) (*models.Item, error) {
	if item.Name == "" {
		return nil, errors.New("item name is required")
	}
	if item.CategoryID == 0 {
		return nil, errors.New("category ID is required")
	}

	// Verify category exists
	var category models.ItemCategory
	if err := s.db.First(&category, item.CategoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	result := s.db.Create(item)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with category data
	s.db.Preload("Category").First(item, item.ID)
	return item, nil
}

// UpdateItem updates an item
func (s *ItemService) UpdateItem(id int, updates map[string]interface{}) (*models.Item, error) {
	var item models.Item
	if err := s.db.First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("item not found")
		}
		return nil, err
	}

	// If category_id is being updated, verify the new category exists
	if newCategoryID, ok := updates["category_id"].(int); ok {
		var category models.ItemCategory
		if err := s.db.First(&category, newCategoryID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("new category not found")
			}
			return nil, err
		}
	}

	if len(updates) > 0 {
		result := s.db.Model(&item).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("Category").First(&item, id)
	return &item, nil
}

// DeleteItem removes an item
func (s *ItemService) DeleteItem(id int) error {
	result := s.db.Delete(&models.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}

// GetItemStatistics returns statistics about items
func (s *ItemService) GetItemStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total items
	var totalItems int64
	s.db.Model(&models.Item{}).Count(&totalItems)
	stats["total_items"] = totalItems

	// Count total categories
	var totalCategories int64
	s.db.Model(&models.ItemCategory{}).Count(&totalCategories)
	stats["total_categories"] = totalCategories

	// Get items per category
	var categoryStats []struct {
		CategoryID   int    `json:"category_id"`
		CategoryName string `json:"category_name"`
		ItemCount    int64  `json:"item_count"`
	}

	s.db.Model(&models.Item{}).
		Select("items.category_id, items_category.category_name, COUNT(items.id) as item_count").
		Joins("LEFT JOIN items_category ON items.category_id = items_category.id").
		Group("items.category_id, items_category.category_name").
		Scan(&categoryStats)

	stats["items_per_category"] = categoryStats

	return stats, nil
}
