package models

import "time"

// ItemCategory represents the items_category table
type ItemCategory struct {
	ID            int       `json:"id" gorm:"primaryKey;column:id"`
	CategoryName  string    `json:"categoryName" gorm:"column:category_name;size:255;not null"`
	Specification string    `json:"specification" gorm:"column:specification;type:text"`
	Quantity      int       `json:"quantity" gorm:"column:quantity;default:0"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	Items         []Item    `json:"items,omitempty" gorm:"foreignKey:CategoryID"`
}

// Item represents the items table
type Item struct {
	ID         int           `json:"id" gorm:"primaryKey;column:id"`
	Name       string        `json:"name" gorm:"column:name;size:255;not null"`
	Year       *int          `json:"year,omitempty" gorm:"column:year"`
	Kondisi    string        `json:"kondisi" gorm:"column:kondisi;size:100"`
	Note       string        `json:"note" gorm:"column:note;type:text"`
	CategoryID int           `json:"categoryId" gorm:"column:category_id;not null"`
	CreatedAt  time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
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
