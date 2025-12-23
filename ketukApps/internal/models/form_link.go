package models

import "time"

// FormLink represents a shareable form link for public schedule requests
// @Description FormLink information for public forms
type FormLink struct {
	ID              uint      `json:"id" gorm:"primaryKey;column:id" example:"1"`
	Code            string    `json:"code" gorm:"column:code;size:20;uniqueIndex;not null" example:"abc123"`
	Title           string    `json:"title" gorm:"column:title;size:200;not null" example:"Lab Reservation - Semester Ganjil 2025"`
	Description     string    `json:"description" gorm:"column:description;type:text" example:"Form untuk reservasi lab semester ganjil"`
	CreatedBy       uint      `json:"createdBy" gorm:"column:created_by;not null" example:"1"`
	Creator         User      `json:"creator" gorm:"foreignKey:CreatedBy;references:ID"`
	PICName         string    `json:"picName" gorm:"column:pic_name;size:100;not null" example:"Dr. John Doe"`
	PICEmail        string    `json:"picEmail" gorm:"column:pic_email;size:100;not null" example:"john.doe@university.ac.id"`
	PICPhone        string    `json:"picPhone" gorm:"column:pic_phone;size:20" example:"08123456789"`
	ExpiresAt       time.Time `json:"expiresAt" gorm:"column:expires_at;not null" example:"2025-01-31T23:59:59Z"`
	MaxSubmissions  *int      `json:"maxSubmissions,omitempty" gorm:"column:max_submissions" example:"50"`
	SubmissionCount int       `json:"submissionCount" gorm:"column:submission_count;default:0" example:"15"`
	IsActive        bool      `json:"isActive" gorm:"column:is_active;default:true" example:"true"`
	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime" example:"2025-01-01T00:00:00Z"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime" example:"2025-01-01T00:00:00Z"`
}

// TableName specifies the table name for FormLink
func (FormLink) TableName() string {
	return "form_links"
}

// IsExpired checks if the form link has expired
func (f *FormLink) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// IsAvailable checks if the form link is available for submissions
func (f *FormLink) IsAvailable() bool {
	if !f.IsActive {
		return false
	}
	if f.IsExpired() {
		return false
	}
	if f.MaxSubmissions != nil && f.SubmissionCount >= *f.MaxSubmissions {
		return false
	}
	return true
}

// GetRemainingSlots returns the number of remaining submission slots (nil if unlimited)
func (f *FormLink) GetRemainingSlots() *int {
	if f.MaxSubmissions == nil {
		return nil
	}
	remaining := *f.MaxSubmissions - f.SubmissionCount
	if remaining < 0 {
		remaining = 0
	}
	return &remaining
}

// CreateFormLinkRequest is the request body for creating a new form link
// @Description Request body for creating a new form link
type CreateFormLinkRequest struct {
	Title          string    `json:"title" binding:"required" example:"Lab Reservation - Semester Ganjil 2025"`
	Description    string    `json:"description" example:"Form untuk reservasi lab semester ganjil"`
	PICName        string    `json:"picName" binding:"required" example:"Dr. John Doe"`
	PICEmail       string    `json:"picEmail" binding:"required,email" example:"john.doe@university.ac.id"`
	PICPhone       string    `json:"picPhone" example:"08123456789"`
	ExpiresAt      time.Time `json:"expiresAt" binding:"required" example:"2025-01-31T23:59:59Z"`
	MaxSubmissions *int      `json:"maxSubmissions" example:"50"`
}

// UpdateFormLinkRequest is the request body for updating a form link
// @Description Request body for updating a form link
type UpdateFormLinkRequest struct {
	Title          string     `json:"title,omitempty" example:"Updated Title"`
	Description    string     `json:"description,omitempty" example:"Updated description"`
	PICName        string     `json:"picName,omitempty" example:"Dr. Jane Doe"`
	PICEmail       string     `json:"picEmail,omitempty" example:"jane.doe@university.ac.id"`
	PICPhone       string     `json:"picPhone,omitempty" example:"08987654321"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty" example:"2025-02-28T23:59:59Z"`
	MaxSubmissions *int       `json:"maxSubmissions,omitempty" example:"100"`
	IsActive       *bool      `json:"isActive,omitempty" example:"false"`
}

// PublicFormSubmitRequest is the request body for submitting a public form
// @Description Request body for submitting a schedule request via public form
type PublicFormSubmitRequest struct {
	SubmitterName  string    `json:"submitterName" binding:"required" example:"John Student"`
	SubmitterEmail string    `json:"submitterEmail" binding:"required,email" example:"student@university.ac.id"`
	SubmitterPhone string    `json:"submitterPhone" example:"08123456789"`
	Title          string    `json:"title" binding:"required" example:"Praktikum Basis Data"`
	Description    string    `json:"description" example:"Praktikum untuk kelas TI-2A"`
	Category       Category  `json:"category" binding:"required" example:"Praktikum"`
	StartDate      time.Time `json:"startDate" binding:"required" example:"2025-01-15T08:00:00Z"`
	EndDate        time.Time `json:"endDate" binding:"required" example:"2025-01-15T10:00:00Z"`
}

// PublicFormResponse represents the public form configuration response
// @Description Public form configuration for display
type PublicFormResponse struct {
	Code           string    `json:"code" example:"abc123"`
	Title          string    `json:"title" example:"Lab Reservation - Semester Ganjil 2025"`
	Description    string    `json:"description" example:"Form untuk reservasi lab semester ganjil"`
	PICName        string    `json:"picName" example:"Dr. John Doe"`
	PICEmail       string    `json:"picEmail" example:"john.doe@university.ac.id"`
	PICPhone       string    `json:"picPhone" example:"08123456789"`
	ExpiresAt      time.Time `json:"expiresAt" example:"2025-01-31T23:59:59Z"`
	RemainingSlots *int      `json:"remainingSlots,omitempty" example:"35"`
	IsAvailable    bool      `json:"isAvailable" example:"true"`
}

// PublicFormSubmitResponse represents the response after submitting a public form
// @Description Response after successful public form submission
type PublicFormSubmitResponse struct {
	ConfirmationCode string `json:"confirmationCode" example:"KTK-2025-001234"`
	Title            string `json:"title" example:"Praktikum Basis Data"`
	StartDate        string `json:"startDate" example:"2025-01-15T08:00:00Z"`
	EndDate          string `json:"endDate" example:"2025-01-15T10:00:00Z"`
	PICName          string `json:"picName" example:"Dr. John Doe"`
	PICEmail         string `json:"picEmail" example:"john.doe@university.ac.id"`
	PICPhone         string `json:"picPhone" example:"08123456789"`
	Message          string `json:"message" example:"Permintaan Anda telah diterima dan sedang ditinjau"`
}
