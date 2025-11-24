package services

import (
	"errors"

	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetAll() ([]models.User, error) {
	var users []models.User
	result := s.db.Find(&users)
	return users, result.Error
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &user, result.Error
}

func (s *UserService) Create(user *models.User) (*models.User, error) {
	// Check if email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "user"
	}

	result := s.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *UserService) Update(id uint, req models.UpdateUserRequest) (*models.User, error) {
	var user models.User

	// Check if user exists
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// Check if new email already exists (if email is being updated)
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["full_name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	if len(updates) > 0 {
		result := s.db.Model(&user).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload user to get updated data
	s.db.First(&user, id)
	return &user, nil
}

func (s *UserService) Delete(id uint) error {
	result := s.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
