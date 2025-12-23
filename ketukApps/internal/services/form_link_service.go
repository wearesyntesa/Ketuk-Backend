package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"ketukApps/internal/models"
)

type FormLinkService struct {
	db *gorm.DB
}

func NewFormLinkService(db *gorm.DB) *FormLinkService {
	return &FormLinkService{db: db}
}

// generateCode generates a unique short code for the form link
func (s *FormLinkService) generateCode() (string, error) {
	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:8], nil
}

// Create creates a new form link
func (s *FormLinkService) Create(req models.CreateFormLinkRequest, createdBy uint) (*models.FormLink, error) {
	// Validate expiration date
	if req.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date must be in the future")
	}

	// Generate unique code
	code, err := s.generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	// Ensure code is unique
	for i := 0; i < 5; i++ {
		var count int64
		s.db.Model(&models.FormLink{}).Where("code = ?", code).Count(&count)
		if count == 0 {
			break
		}
		code, err = s.generateCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique code: %w", err)
		}
	}

	formLink := &models.FormLink{
		Code:           code,
		Title:          req.Title,
		Description:    req.Description,
		CreatedBy:      createdBy,
		PICName:        req.PICName,
		PICEmail:       req.PICEmail,
		PICPhone:       req.PICPhone,
		ExpiresAt:      req.ExpiresAt,
		MaxSubmissions: req.MaxSubmissions,
		IsActive:       true,
	}

	if err := s.db.Create(formLink).Error; err != nil {
		return nil, fmt.Errorf("failed to create form link: %w", err)
	}

	// Reload with creator
	if err := s.db.Preload("Creator").First(formLink, formLink.ID).Error; err != nil {
		return nil, err
	}

	return formLink, nil
}

// GetAll retrieves all form links
func (s *FormLinkService) GetAll() ([]models.FormLink, error) {
	var formLinks []models.FormLink
	if err := s.db.Preload("Creator").Order("created_at DESC").Find(&formLinks).Error; err != nil {
		return nil, err
	}
	return formLinks, nil
}

// GetByID retrieves a form link by ID
func (s *FormLinkService) GetByID(id uint) (*models.FormLink, error) {
	var formLink models.FormLink
	if err := s.db.Preload("Creator").First(&formLink, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("form link not found")
		}
		return nil, err
	}
	return &formLink, nil
}

// GetByCode retrieves a form link by code (for public access)
func (s *FormLinkService) GetByCode(code string) (*models.FormLink, error) {
	var formLink models.FormLink
	if err := s.db.Where("code = ?", code).First(&formLink).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("form link not found")
		}
		return nil, err
	}
	return &formLink, nil
}

// GetPublicFormByCode retrieves a form link for public display
func (s *FormLinkService) GetPublicFormByCode(code string) (*models.PublicFormResponse, error) {
	formLink, err := s.GetByCode(code)
	if err != nil {
		return nil, err
	}

	response := &models.PublicFormResponse{
		Code:           formLink.Code,
		Title:          formLink.Title,
		Description:    formLink.Description,
		PICName:        formLink.PICName,
		PICEmail:       formLink.PICEmail,
		PICPhone:       formLink.PICPhone,
		ExpiresAt:      formLink.ExpiresAt,
		RemainingSlots: formLink.GetRemainingSlots(),
		IsAvailable:    formLink.IsAvailable(),
	}

	return response, nil
}

// Update updates a form link
func (s *FormLinkService) Update(id uint, req models.UpdateFormLinkRequest) (*models.FormLink, error) {
	formLink, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Title != "" {
		formLink.Title = req.Title
	}
	if req.Description != "" {
		formLink.Description = req.Description
	}
	if req.PICName != "" {
		formLink.PICName = req.PICName
	}
	if req.PICEmail != "" {
		formLink.PICEmail = req.PICEmail
	}
	if req.PICPhone != "" {
		formLink.PICPhone = req.PICPhone
	}
	if req.ExpiresAt != nil {
		if req.ExpiresAt.Before(time.Now()) {
			return nil, errors.New("expiration date must be in the future")
		}
		formLink.ExpiresAt = *req.ExpiresAt
	}
	if req.MaxSubmissions != nil {
		formLink.MaxSubmissions = req.MaxSubmissions
	}
	if req.IsActive != nil {
		formLink.IsActive = *req.IsActive
	}

	if err := s.db.Save(formLink).Error; err != nil {
		return nil, fmt.Errorf("failed to update form link: %w", err)
	}

	return formLink, nil
}

// Delete deletes a form link
func (s *FormLinkService) Delete(id uint) error {
	result := s.db.Delete(&models.FormLink{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("form link not found")
	}
	return nil
}

// Deactivate deactivates a form link
func (s *FormLinkService) Deactivate(id uint) (*models.FormLink, error) {
	formLink, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	formLink.IsActive = false
	if err := s.db.Save(formLink).Error; err != nil {
		return nil, err
	}

	return formLink, nil
}

// Clone creates a copy of an existing form link with a new code and expiration
func (s *FormLinkService) Clone(id uint, newExpiresAt time.Time, createdBy uint) (*models.FormLink, error) {
	original, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if newExpiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date must be in the future")
	}

	// Generate new unique code
	code, err := s.generateCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code: %w", err)
	}

	cloned := &models.FormLink{
		Code:           code,
		Title:          original.Title,
		Description:    original.Description,
		CreatedBy:      createdBy,
		PICName:        original.PICName,
		PICEmail:       original.PICEmail,
		PICPhone:       original.PICPhone,
		ExpiresAt:      newExpiresAt,
		MaxSubmissions: original.MaxSubmissions,
		IsActive:       true,
	}

	if err := s.db.Create(cloned).Error; err != nil {
		return nil, fmt.Errorf("failed to clone form link: %w", err)
	}

	// Reload with creator
	if err := s.db.Preload("Creator").First(cloned, cloned.ID).Error; err != nil {
		return nil, err
	}

	return cloned, nil
}

// SubmitPublicForm handles a public form submission
func (s *FormLinkService) SubmitPublicForm(code string, req models.PublicFormSubmitRequest) (*models.PublicFormSubmitResponse, error) {
	formLink, err := s.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Check if form is available
	if !formLink.IsAvailable() {
		if !formLink.IsActive {
			return nil, errors.New("this form is no longer active")
		}
		if formLink.IsExpired() {
			return nil, errors.New("this form has expired")
		}
		if formLink.MaxSubmissions != nil && formLink.SubmissionCount >= *formLink.MaxSubmissions {
			return nil, errors.New("this form has reached its maximum number of submissions")
		}
		return nil, errors.New("this form is not available")
	}

	// Validate dates
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Create ticket with a system user or null user
	// We need to handle this - for now use user ID 1 (admin) as a placeholder
	// In production, you might want a dedicated "system" user
	ticket := &models.Ticket{
		UserID:         1, // System/admin user - public submissions are attributed to system
		Title:          req.Title,
		Description:    req.Description,
		Category:       req.Category,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         models.StatusPending,
		FormLinkID:     &formLink.ID,
		SubmitterName:  req.SubmitterName,
		SubmitterEmail: req.SubmitterEmail,
		SubmitterPhone: req.SubmitterPhone,
		Email:          req.SubmitterEmail,
		Phone:          req.SubmitterPhone,
		PIC:            req.SubmitterName,
	}

	// Use transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create ticket
	if err := tx.Create(ticket).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	// Increment submission count
	if err := tx.Model(formLink).Update("submission_count", gorm.Expr("submission_count + 1")).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update submission count: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Generate confirmation code
	confirmationCode := fmt.Sprintf("KTK-%d-%06d", time.Now().Year(), ticket.ID)

	response := &models.PublicFormSubmitResponse{
		ConfirmationCode: confirmationCode,
		Title:            ticket.Title,
		StartDate:        ticket.StartDate.Format(time.RFC3339),
		EndDate:          ticket.EndDate.Format(time.RFC3339),
		PICName:          formLink.PICName,
		PICEmail:         formLink.PICEmail,
		PICPhone:         formLink.PICPhone,
		Message:          "Permintaan Anda telah diterima dan sedang ditinjau. Kami akan menghubungi Anda segera.",
	}

	return response, nil
}

// GetSubmissionsByFormLinkID retrieves all tickets submitted through a form link
func (s *FormLinkService) GetSubmissionsByFormLinkID(formLinkID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	if err := s.db.Where("form_link_id = ?", formLinkID).Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}
