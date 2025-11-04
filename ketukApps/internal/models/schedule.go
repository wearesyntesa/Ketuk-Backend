package models

import "time"

// ScheduleTicket represents the schedule_ticket table (schedule from tickets)
type ScheduleTicket struct {
	IDSchedule  uint      `json:"idSchedule" gorm:"primaryKey;column:id_schedule"`
	Title       string    `json:"title" gorm:"column:title;size:255;not null"`
	StartDate   time.Time `json:"startDate" gorm:"column:start_date;not null"`
	EndDate     time.Time `json:"endDate" gorm:"column:end_date;not null"`
	UserID      uint      `json:"userId" gorm:"column:user_id;not null"`
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
	ID        uint             `json:"id" gorm:"primaryKey;column:id"`
	Tahun     int              `json:"tahun" gorm:"column:tahun;not null"`
	Semester  SemesterCategory `json:"semester" gorm:"column:semester;type:semester_category;not null"`
	CreatedAt time.Time        `json:"createdAt" gorm:"column:created_at;not null"`
	StartDate time.Time        `json:"startDate" gorm:"column:start_date;not null"`
	UserID    uint             `json:"userId" gorm:"column:user_id;not null"`
	EndDate   time.Time        `json:"endDate" gorm:"column:end_date;not null"`
	User      *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// ScheduleReguler represents the schedule_reguler table
type ScheduleReguler struct {
	IDSchedule uint      `json:"idSchedule" gorm:"primaryKey;column:id_schedule"`
	Title      string    `json:"title" gorm:"column:title;size:255;not null"`
	StartDate  time.Time `json:"startDate" gorm:"column:start_date;not null"`
	EndDate    time.Time `json:"endDate" gorm:"column:end_date;not null"`
	UserID     uint      `json:"userId" gorm:"column:user_id;not null"`
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
