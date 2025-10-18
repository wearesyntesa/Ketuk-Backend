package models

import "time"

// Ticket represents the tickets table with all fields flattened
type Ticket struct {
	ID          uint       `json:"id" gorm:"primaryKey;column:id"`
	UserID      uint       `json:"userId" gorm:"column:user_id;not null"`
	User        User       `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Title       string     `json:"title" gorm:"column:title;size:100;not null"`
	Description string     `json:"description" gorm:"column:description;type:text"`
	Status      string     `json:"status" gorm:"column:status;type:ticket_status;default:pending"`
	IDSchedule  *int       `json:"idSchedule,omitempty" gorm:"column:id_schedule"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	ApprovedAt  *time.Time `json:"approvedAt,omitempty" gorm:"column:approved_at"`
}

// Category defines the type of request
type Category string

const (
	Kelas     Category = "Kelas"
	Lainnya   Category = "Lainnya"
	Praktikum Category = "Praktikum"
	Skripsi   Category = "Skripsi"
)

// TicketStatus defines the status of a ticket
type TicketStatus string

const (
	StatusPending  TicketStatus = "pending"
	StatusAccepted TicketStatus = "accepted"
	StatusRejected TicketStatus = "rejected"
)

// CreateTicketRequest is the request body for creating a new ticket
type CreateTicketRequest struct {
	UserID      uint   `json:"userId" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// UpdateTicketRequest is the request body for updating a ticket
type UpdateTicketRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type TicketResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Ticket  Ticket `json:"ticket,omitempty"`
	Error   string `json:"error,omitempty"`
}
