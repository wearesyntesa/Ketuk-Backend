package models

import (
	"time"
)

// User represents a user in the system
// @Description User account information
type User struct {
	ID uint `json:"id" gorm:"primaryKey" example:"1"`
	// GoogleSub is kept for backwards compatibility with existing data
	GoogleSub string    `json:"google_sub,omitempty" gorm:"uniqueIndex;size:255" example:"google-oauth2|123456789"`
	Name      string    `json:"name" binding:"required" gorm:"column:full_name;size:255;not null" example:"John Doe"`
	Email     string    `json:"email" binding:"required,email" gorm:"uniqueIndex;size:255;not null" example:"john.doe@example.com"`
	Password  string    `json:"-" gorm:"size:255"` // Password hash, not included in JSON responses
	Role      string    `json:"role" gorm:"type:user_role;default:user" example:"user"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUserRequest represents the request body for creating a new user
// @Description Request body for creating a new user
type CreateUserRequest struct {
	GoogleSub string `json:"google_sub" binding:"required" example:"google-oauth2|123456789"`
	Name      string `json:"name" binding:"required" example:"John Doe"`
	Email     string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}

// UpdateUserRequest represents the request body for updating a user
// @Description Request body for updating user information
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" example:"Jane Doe"`
	Email string `json:"email,omitempty" example:"jane.doe@example.com"`
	Role  string `json:"role,omitempty" example:"admin"`
}

// APIResponse represents a standard API response
// @Description Standard API response format
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty" example:""`
}

// HealthResponse represents a health check response
// @Description Health check response format
type HealthResponse struct {
	Status    string `json:"status" example:"healthy"`
	Timestamp string `json:"timestamp" example:"2023-01-01T00:00:00Z"`
	Version   string `json:"version" example:"1.0.0"`
}

// RefreshTokenRequest represents the request body for refreshing token
// @Description Request body for refreshing access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh_token_123456789"`
}
