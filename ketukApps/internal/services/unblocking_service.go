package services

import (
	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type UnblockingService struct {
	db *gorm.DB
}

func NewUnblockingService(db *gorm.DB) *UnblockingService {
	return &UnblockingService{
		db: db,
	}
}

// Create unblocking request
func (s *UnblockingService) Create(unblocking *models.Unblocking) (*models.Unblocking, error) {
	result := s.db.Create(unblocking)
	return unblocking, result.Error
}

// GetByID returns an unblocking request by its ID
func (s *UnblockingService) GetByID(id int) (*models.Unblocking, error) {
	var unblocking models.Unblocking
	result := s.db.Preload("User").First(&unblocking, id)
	return &unblocking, result.Error
}

// GetByUserID returns all unblocking requests for a specific user
func (s *UnblockingService) GetByUserID(userID int) ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Preload("User").Where("user_id = ?", userID).Find(&unblockings)
	return unblockings, result.Error
}

// GetAll returns all unblocking requests
func (s *UnblockingService) GetAll() ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Preload("User").Find(&unblockings)
	return unblockings, result.Error
}

// Delete an unblocking request by its ID
func (s *UnblockingService) Delete(id int) error {
	result := s.db.Delete(&models.Unblocking{}, id)
	return result.Error
}

// Get unblocked schedules that are due for unblocking
func (s *UnblockingService) GetDueUnblockings() ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Where("scheduled_at <= NOW()").Find(&unblockings)
	return unblockings, result.Error
}
