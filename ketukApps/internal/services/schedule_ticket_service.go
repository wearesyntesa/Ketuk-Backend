package services

import (
	"errors"
	"time"

	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type ScheduleService struct {
	db *gorm.DB
}

func NewScheduleService(db *gorm.DB) *ScheduleService {
	return &ScheduleService{
		db: db,
	}
}

// ScheduleTicket methods

// GetAllScheduleTickets returns all schedule tickets
func (s *ScheduleService) GetAllScheduleTickets() ([]models.ScheduleTicket, error) {
	var schedules []models.ScheduleTicket
	result := s.db.Preload("User").Preload("Tickets").Find(&schedules)
	return schedules, result.Error
}

// GetScheduleTicketByID returns a schedule ticket by its ID
func (s *ScheduleService) GetScheduleTicketByID(id uint) (*models.ScheduleTicket, error) {
	var schedule models.ScheduleTicket
	result := s.db.Preload("User").Preload("Tickets").First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("schedule ticket not found")
	}
	return &schedule, result.Error
}

// GetScheduleTicketsByUserID returns all schedule tickets for a specific user
func (s *ScheduleService) GetScheduleTicketsByUserID(userID uint) ([]models.ScheduleTicket, error) {
	var schedules []models.ScheduleTicket
	result := s.db.Preload("User").Preload("Tickets").Where("user_id = ?", userID).Find(&schedules)
	return schedules, result.Error
}

// GetScheduleTicketsByCategory returns all schedule tickets with a specific category
func (s *ScheduleService) GetScheduleTicketsByCategory(category models.Category) ([]models.ScheduleTicket, error) {
	var schedules []models.ScheduleTicket
	result := s.db.Preload("User").Preload("Tickets").Where("kategori = ?", string(category)).Find(&schedules)
	return schedules, result.Error
}

// GetScheduleTicketsByDateRange returns schedule tickets within a date range
func (s *ScheduleService) GetScheduleTicketsByDateRange(startDate, endDate time.Time) ([]models.ScheduleTicket, error) {
	var schedules []models.ScheduleTicket
	result := s.db.Preload("User").Preload("Tickets").
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Find(&schedules)
	return schedules, result.Error
}

// CreateScheduleTicket creates a new schedule ticket
func (s *ScheduleService) CreateScheduleTicket(schedule *models.ScheduleTicket) (*models.ScheduleTicket, error) {
	if schedule.Title == "" {
		return nil, errors.New("title is required")
	}
	if schedule.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	result := s.db.Create(schedule)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with user data
	s.db.Preload("User").Preload("Tickets").First(schedule, schedule.IDSchedule)
	return schedule, nil
}

// UpdateScheduleTicket updates a schedule ticket
func (s *ScheduleService) UpdateScheduleTicket(id uint, updates map[string]interface{}) (*models.ScheduleTicket, error) {
	var schedule models.ScheduleTicket
	if err := s.db.First(&schedule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule ticket not found")
		}
		return nil, err
	}

	if len(updates) > 0 {
		result := s.db.Model(&schedule).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("User").Preload("Tickets").First(&schedule, id)
	return &schedule, nil
}

// DeleteScheduleTicket removes a schedule ticket
func (s *ScheduleService) DeleteScheduleTicket(id uint) error {
	result := s.db.Delete(&models.ScheduleTicket{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("schedule ticket not found")
	}
	return nil
}

// ScheduleReguler methods

// GetAllScheduleReguler returns all regular schedules
func (s *ScheduleService) GetAllScheduleReguler() ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Find(&schedules)
	return schedules, result.Error
}

// GetScheduleRegulerByID returns a regular schedule by its ID
func (s *ScheduleService) GetScheduleRegulerByID(id uint) (*models.ScheduleReguler, error) {
	var schedule models.ScheduleReguler
	result := s.db.Preload("User").First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("schedule reguler not found")
	}
	return &schedule, result.Error
}

// GetScheduleRegulerByUserID returns all regular schedules for a specific user
func (s *ScheduleService) GetScheduleRegulerByUserID(userID uint) ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Where("user_id = ?", userID).Find(&schedules)
	return schedules, result.Error
}

// GetScheduleRegulerByDateRange returns regular schedules within a date range
func (s *ScheduleService) GetScheduleRegulerByDateRange(startDate, endDate time.Time) ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Find(&schedules)
	return schedules, result.Error
}

// CreateScheduleReguler creates a new regular schedule
func (s *ScheduleService) CreateScheduleReguler(schedule *models.ScheduleReguler) (*models.ScheduleReguler, error) {
	if schedule.Title == "" {
		return nil, errors.New("title is required")
	}
	if schedule.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	result := s.db.Create(schedule)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with user data
	s.db.Preload("User").First(schedule, schedule.IDSchedule)
	return schedule, nil
}

// UpdateScheduleReguler updates a regular schedule
func (s *ScheduleService) UpdateScheduleReguler(id uint, updates map[string]interface{}) (*models.ScheduleReguler, error) {
	var schedule models.ScheduleReguler
	if err := s.db.First(&schedule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule reguler not found")
		}
		return nil, err
	}

	if len(updates) > 0 {
		result := s.db.Model(&schedule).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("User").First(&schedule, id)
	return &schedule, nil
}

// DeleteScheduleReguler removes a regular schedule
func (s *ScheduleService) DeleteScheduleReguler(id uint) error {
	result := s.db.Delete(&models.ScheduleReguler{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("schedule reguler not found")
	}
	return nil
}

// Unblocking methods

// GetAllUnblocking returns all unblocking records
func (s *ScheduleService) GetAllUnblocking() ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Preload("User").Find(&unblockings)
	return unblockings, result.Error
}

// GetUnblockingByID returns an unblocking record by its ID
func (s *ScheduleService) GetUnblockingByID(id uint) (*models.Unblocking, error) {
	var unblocking models.Unblocking
	result := s.db.Preload("User").First(&unblocking, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("unblocking not found")
	}
	return &unblocking, result.Error
}

// GetUnblockingByUserID returns all unblocking records for a specific user
func (s *ScheduleService) GetUnblockingByUserID(userID uint) ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Preload("User").Where("user_id = ?", userID).Find(&unblockings)
	return unblockings, result.Error
}

// GetUnblockingBySemester returns all unblocking records for a specific semester
func (s *ScheduleService) GetUnblockingBySemester(tahun int, semester models.SemesterCategory) ([]models.Unblocking, error) {
	var unblockings []models.Unblocking
	result := s.db.Preload("User").
		Where("tahun = ? AND semester = ?", tahun, string(semester)).
		Find(&unblockings)
	return unblockings, result.Error
}

// CreateUnblocking creates a new unblocking record
func (s *ScheduleService) CreateUnblocking(unblocking *models.Unblocking) (*models.Unblocking, error) {
	if unblocking.Tahun == 0 {
		return nil, errors.New("tahun is required")
	}
	if unblocking.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	result := s.db.Create(unblocking)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with user data
	s.db.Preload("User").First(unblocking, unblocking.ID)
	return unblocking, nil
}

// UpdateUnblocking updates an unblocking record
func (s *ScheduleService) UpdateUnblocking(id uint, updates map[string]interface{}) (*models.Unblocking, error) {
	var unblocking models.Unblocking
	if err := s.db.First(&unblocking, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("unblocking not found")
		}
		return nil, err
	}

	if len(updates) > 0 {
		result := s.db.Model(&unblocking).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("User").First(&unblocking, id)
	return &unblocking, nil
}

// DeleteUnblocking removes an unblocking record
func (s *ScheduleService) DeleteUnblocking(id uint) error {
	result := s.db.Delete(&models.Unblocking{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("unblocking not found")
	}
	return nil
}
