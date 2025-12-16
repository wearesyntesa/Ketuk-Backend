package models

import "time"

// ScheduleTicket represents the schedule_ticket table (schedule from tickets)
type ScheduleTicket struct {
	IDSchedule  int       `json:"idSchedule" gorm:"primaryKey;column:id_schedule"`
	Title       string    `json:"title" gorm:"column:title;size:255;not null"`
	StartDate   time.Time `json:"startDate" gorm:"column:start_date;not null"`
	EndDate     time.Time `json:"endDate" gorm:"column:end_date;not null"`
	UserID      int       `json:"userId" gorm:"column:user_id;not null"`
	Kategori    Category  `json:"kategori" gorm:"column:kategori;type:ticket_category;not null"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime"`
	User        *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tickets     []Ticket  `json:"tickets,omitempty" gorm:"foreignKey:IDSchedule;references:IDSchedule"`
}

// SemesterCategory defines the semester type
type SemesterCategory string

const (
	SemesterGanjil SemesterCategory = "Ganjil"
	SemesterGenap  SemesterCategory = "Genap"
)

// Unblocking represents the unblocking table (semester unblocking)
type Unblocking struct {
	ID        int              `json:"id" gorm:"primaryKey;column:id"`
	Tahun     int              `json:"tahun" gorm:"column:tahun;not null"`
	Semester  SemesterCategory `json:"semester" gorm:"column:semester;type:semester_category;not null"`
	CreatedAt time.Time        `json:"createdAt" gorm:"column:created_at;not null"`
	StartDate time.Time        `json:"startDate" gorm:"column:start_date;not null"`
	UserID    int              `json:"userId" gorm:"column:user_id;not null"`
	EndDate   time.Time        `json:"endDate" gorm:"column:end_date;not null"`
	User      *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ScheduleReguler represents the schedule_reguler table
type ScheduleReguler struct {
	IDSchedule int       `json:"idSchedule" gorm:"primaryKey;column:id_schedule"`
	Title      string    `json:"title" gorm:"column:title;size:255;not null"`
	StartDate  time.Time `json:"startDate" gorm:"column:start_date;not null"`
	EndDate    time.Time `json:"endDate" gorm:"column:end_date;not null"`
	UserID     int       `json:"userId" gorm:"column:user_id;not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName overrides the table name for ScheduleTicket
func (ScheduleTicket) TableName() string {
	return "schedule_ticket"
}

// TableName overrides the table name for Unblocking
func (Unblocking) TableName() string {
	return "unblocking"
}

// TableName overrides the table name for ScheduleReguler
func (ScheduleReguler) TableName() string {
	return "schedule_reguler"
}

// UnblockingResponse represents a unblocking response
// @Description Unblocking response format
type UnblockingResponse struct {
	Success    bool       `json:"success" example:"true"`
	Message    string     `json:"message" example:"Unblocking operation completed successfully"`
	Unblocking Unblocking `json:"unblocking,omitempty"`
}

type UnblockingsResponse struct {
	Success     bool         `json:"success" example:"true"`
	Message     string       `json:"message" example:"Unblocking operation completed successfully"`
	Unblockings []Unblocking `json:"unblockings,omitempty"`
}

type CreateUnblockingRequest struct {
	Tahun     int              `json:"tahun" binding:"required" example:"2023"`
	Semester  SemesterCategory `json:"semester" binding:"required,oneof=Ganjil Genap" example:"Ganjil"`
	StartDate time.Time        `json:"startDate" binding:"required" example:"2023-09-01T00:00:00Z"`
	EndDate   time.Time        `json:"endDate" binding:"required" example:"2023-12-31T00:00:00Z"`
	UserID    int              `json:"userId" binding:"required" example:"1"`
}

// CreateScheduleRegulerRequest represents request to create schedule reguler
type CreateScheduleRegulerRequest struct {
	Title     string    `json:"title" binding:"required" example:"Regular Maintenance Schedule"`
	StartDate time.Time `json:"startDate" binding:"required" example:"2023-12-01T09:00:00Z"`
	EndDate   time.Time `json:"endDate" binding:"required" example:"2023-12-01T17:00:00Z"`
	UserID    int       `json:"userId" binding:"required" example:"1"`
}

// UpdateScheduleRegulerRequest represents request to update schedule reguler
type UpdateScheduleRegulerRequest struct {
	Title     *string    `json:"title,omitempty" example:"Updated Schedule Title"`
	StartDate *time.Time `json:"startDate,omitempty" example:"2023-12-01T09:00:00Z"`
	EndDate   *time.Time `json:"endDate,omitempty" example:"2023-12-01T17:00:00Z"`
}

// CreateScheduleTicketRequest represents request to create schedule ticket
type CreateScheduleTicketRequest struct {
	Title       string    `json:"title" binding:"required" example:"Network Maintenance"`
	StartDate   time.Time `json:"startDate" binding:"required" example:"2023-12-01T09:00:00Z"`
	EndDate     time.Time `json:"endDate" binding:"required" example:"2023-12-01T17:00:00Z"`
	UserID      int       `json:"userId" binding:"required" example:"1"`
	Kategori    Category  `json:"kategori" binding:"required" example:"barang"`
	Description string    `json:"description" example:"Scheduled network maintenance"`
}

// UpdateScheduleTicketRequest represents request to update schedule ticket
type UpdateScheduleTicketRequest struct {
	Title       *string    `json:"title,omitempty" example:"Updated Schedule"`
	StartDate   *time.Time `json:"startDate,omitempty" example:"2023-12-01T09:00:00Z"`
	EndDate     *time.Time `json:"endDate,omitempty" example:"2023-12-01T17:00:00Z"`
	Kategori    *Category  `json:"kategori,omitempty" example:"barang"`
	Description *string    `json:"description,omitempty" example:"Updated description"`
}
