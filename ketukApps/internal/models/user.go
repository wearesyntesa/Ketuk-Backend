package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	GoogleSub string    `json:"google_sub" gorm:"uniqueIndex;size:255;not null"`
	Name      string    `json:"name" binding:"required" gorm:"column:full_name;size:255;not null"`
	Email     string    `json:"email" binding:"required,email" gorm:"uniqueIndex;size:255;not null"`
	Role      string    `json:"role" gorm:"type:user_role;default:user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	GoogleSub string `json:"google_sub" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}
