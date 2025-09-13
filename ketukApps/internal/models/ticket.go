package models

import (
	"time"
)

type Ticket struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Title       string     `json:"title" gorm:"size:255;not null"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status" gorm:"type:ticket_status;default:pending"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
}

type Category string

const (
	CategoryLainnya  Category = "Lainnya"
	CategoryClass    Category = "Kelas"
	CategoryPractice Category = "Praktikum"
	CategoryThesis   Category = "Skripsi"
)

// For backward compatibility with existing handlers
type RequestData struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description,omitempty"`
	Category    Category `json:"category,omitempty"`
	CreatedAt   string   `json:"created_at,omitempty"`
	ApprovedAt  string   `json:"approved_at,omitempty"`
}

type TicketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Ticket  Ticket `json:"ticket,omitempty"`
	Error   string `json:"error,omitempty"`
}
