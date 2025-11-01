package services
import (
	"errors"
	"time"

	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type ScheduleRegulerService struct {
	db *gorm.DB
}

func NewScheduleRegulerService(db *gorm.DB) *ScheduleRegulerService {
	return &ScheduleRegulerService{
		db: db,
	}
}

// ScheduleReguler methods

// GetAllScheduleRegulers returns all schedule regulers
func (s *ScheduleRegulerService) GetAllScheduleRegulers() ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Preload("Activities").Find(&schedules)
	return schedules, result.Error
}

// GetScheduleRegulerByID returns a schedule reguler by its ID
func (s *ScheduleRegulerService) GetScheduleRegulerByID(id int) (*models.ScheduleReguler, error) {
	var schedule models.ScheduleReguler
	result := s.db.Preload("User").Preload("Activities").First(&schedule, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("schedule reguler not found")
	}
	return &schedule, result.Error
}

// GetScheduleRegulersByUserID returns all schedule regulers for a specific user
func (s *ScheduleRegulerService) GetScheduleRegulersByUserID(userID int) ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Preload("Activities").Where("user_id = ?", userID).Find(&schedules)
	return schedules, result.Error
}

// GetScheduleRegulersByCategory returns all schedule regulers with a specific category
func (s *ScheduleRegulerService) GetScheduleRegulersByCategory(category models.Category) ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Preload("Activities").Where("kategori = ?", string(category)).Find(&schedules)
	return schedules, result.Error
}

// GetScheduleRegulersByDateRange returns schedule regulers within a date range
func (s *ScheduleRegulerService) GetScheduleRegulersByDateRange(startDate, endDate time.Time) ([]models.ScheduleReguler, error) {
	var schedules []models.ScheduleReguler
	result := s.db.Preload("User").Preload("Activities").
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Find(&schedules)
	return schedules, result.Error
}
// CreateScheduleReguler creates a new schedule reguler
func (s *ScheduleRegulerService) CreateScheduleReguler(schedule *models.ScheduleReguler) (*models.ScheduleReguler, error){
	result := s.db.Create(schedule)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with activities data
	s.db.Preload("Activities").First(schedule, schedule.IDSchedule)
	return schedule, nil
}

// UpdateScheduleReguler updates an existing schedule reguler
func (s *ScheduleRegulerService) UpdateScheduleReguler(schedule *models.ScheduleReguler) error {
	result := s.db.Save(schedule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("schedule reguler not found")
	}
	return nil
}

// DeleteScheduleReguler deletes a schedule reguler by its ID
func (s *ScheduleRegulerService) DeleteScheduleReguler(id int) error {
	result := s.db.Delete(&models.ScheduleReguler{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("schedule reguler not found")
	}
	return nil
}

// Additional methods can be added as needed
func (s *ScheduleRegulerService) DeleteScheduleRegulersByUserID(userID int) error {
	result := s.db.Where("user_id = ?", userID).Delete(&models.ScheduleReguler{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *ScheduleRegulerService) DeleteScheduleRegulersByDateRange(startDate, endDate time.Time) error {
	result := s.db.Where("start_date >= ? AND end_date <= ?", startDate, endDate).Delete(&models.ScheduleReguler{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
