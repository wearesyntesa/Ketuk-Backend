package models

import "time"

// ItemCategory represents the items_category table
// @Description Item category information
type ItemCategory struct {
	ID            int       `json:"id" gorm:"primaryKey;column:id" example:"1"`
	CategoryName  string    `json:"categoryName" gorm:"column:category_name;size:255;not null" example:"Komputer"`
	Specification string    `json:"specification" gorm:"column:specification;type:text" example:"Komputer Desktop Intel Core i5"`
	Quantity      int       `json:"quantity" gorm:"column:quantity;default:0" example:"10"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" example:"2023-01-01T00:00:00Z"`
	Items         []Item    `json:"items,omitempty" gorm:"foreignKey:CategoryID"`
}

// Item represents the items table
// @Description Item information
type Item struct {
	ID         int           `json:"id" gorm:"primaryKey;column:id" example:"1"`
	Name       string        `json:"name" gorm:"column:name;size:255;not null" example:"PC-001"`
	Year       *int          `json:"year,omitempty" gorm:"column:year" example:"2023"`
	Kondisi    string        `json:"kondisi" gorm:"column:kondisi;size:100" example:"Baik"`
	Note       string        `json:"note" gorm:"column:note;type:text" example:"Kondisi normal, ready to use"`
	CategoryID int           `json:"categoryId" gorm:"column:category_id;not null" example:"1"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" example:"2023-01-01T00:00:00Z"`
	Category   *ItemCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName overrides the table name for ItemCategory
func (ItemCategory) TableName() string {
	return "items_category"
}

// TableName overrides the table name for Item
func (Item) TableName() string {
	return "items"
}
