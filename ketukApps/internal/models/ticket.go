package models

import "time"

// Ticket represents the tickets table with all fields flattened
// @Description Ticket information
type Ticket struct {
	ID             uint         `json:"id" gorm:"primaryKey;column:id" example:"1"`
	UserID         uint         `json:"userId" gorm:"column:user_id;not null" example:"1"`
	User           User         `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Title          string       `json:"title" gorm:"column:title;size:100;not null" example:"Room Booking Request"`
	Description    string       `json:"description" gorm:"column:description;type:text" example:"Need to book conference room for meeting"`
	Category       Category     `json:"category" gorm:"column:category;type:ticket_category;not null;default:'Lainnya'" example:"Kelas"`
	StartDate      time.Time    `json:"startDate" gorm:"column:start_date;not null" example:"2023-12-01T09:00:00Z"`
	EndDate        time.Time    `json:"endDate" gorm:"column:end_date;not null" example:"2023-12-01T17:00:00Z"`
	Email          string       `json:"email" gorm:"column:email;size:100;not null;default:''" example:"user@example.com"`
	Phone          string       `json:"phone" gorm:"column:phone;size:15;not null;default:''" example:"08123456789"`
	PIC            string       `json:"pic" gorm:"column:pic;size:100;not null;default:''" example:"John Doe"`
	Status         TicketStatus `json:"status" gorm:"column:status;type:ticket_status;default:pending" example:"pending"`
	IDSchedule     *int         `json:"idSchedule,omitempty" gorm:"column:id_schedule" example:"1"`
	CreatedAt      time.Time    `json:"createdAt" gorm:"column:created_at;autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt      time.Time    `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" example:"2023-01-01T00:00:00Z"`
	ApprovedAt     *time.Time   `json:"approvedAt,omitempty" gorm:"column:approved_at" example:"2023-01-02T00:00:00Z"`
	Reason         string       `json:"reason" gorm:"column:reason;type:text;not null;default:'No reason provided'" example:"No reason provided"`
	FormLinkID     *uint        `json:"formLinkId,omitempty" gorm:"column:form_link_id" example:"1"`
	SubmitterName  string       `json:"submitterName,omitempty" gorm:"column:submitter_name;size:100" example:"John Student"`
	SubmitterEmail string       `json:"submitterEmail,omitempty" gorm:"column:submitter_email;size:100" example:"student@university.ac.id"`
	SubmitterPhone string       `json:"submitterPhone,omitempty" gorm:"column:submitter_phone;size:20" example:"08123456789"`
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
// @Description Request body for creating a new ticket
type CreateTicketRequest struct {
	UserID      uint      `json:"userId" binding:"required" example:"1"`
	Title       string    `json:"title" binding:"required" example:"Room Booking Request"`
	Description string    `json:"description" example:"Need to book conference room for meeting"`
	Category    Category  `json:"category" binding:"required" example:"Kelas"`
	StartDate   time.Time `json:"startDate" binding:"required" example:"2023-12-01T09:00:00Z"`
	EndDate     time.Time `json:"endDate" binding:"required" example:"2023-12-01T17:00:00Z"`
}

// UpdateTicketRequest is the request body for updating a ticket
// @Description Request body for updating a ticket
type UpdateTicketRequest struct {
	Title       string `json:"title,omitempty" example:"Updated Room Booking Request"`
	Description string `json:"description,omitempty" example:"Updated description for the booking"`
}

// TicketResponse represents a ticket response
// @Description Ticket response format
type TicketResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Ticket operation completed successfully"`
	Ticket  Ticket `json:"ticket,omitempty"`
	Error   string `json:"error,omitempty" example:""`
}
