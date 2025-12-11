package services

import (
	"errors"
	"fmt"
	"time"

	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type TicketService struct {
	db           *gorm.DB
	auditService *AuditService
}

func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{
		db:           db,
		auditService: NewAuditService(db),
	}
}

// GetAll returns all tickets
func (s *TicketService) GetAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	result := s.db.Preload("User").Find(&tickets)
	return tickets, result.Error
}

// GetByID returns a ticket by its ID
func (s *TicketService) GetByID(id uint) (*models.Ticket, error) {
	var ticket models.Ticket
	result := s.db.Preload("User").First(&ticket, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("ticket not found")
	}
	return &ticket, result.Error
}

// GetByUserID returns all tickets for a specific user
func (s *TicketService) GetByUserID(userID uint) ([]models.Ticket, error) {
	var tickets []models.Ticket
	result := s.db.Preload("User").Where("user_id = ?", userID).Find(&tickets)
	if result.Error != nil {
		return nil, result.Error
	}
	if len(tickets) == 0 {
		return nil, errors.New("no tickets found for user")
	}
	return tickets, nil
}

// GetByStatus returns all tickets with a specific status
func (s *TicketService) GetByStatus(status string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	result := s.db.Preload("User").Where("status = ?", status).Find(&tickets)
	return tickets, result.Error
}

// Create creates a new ticket using the new model structure
func (s *TicketService) Create(userID uint, title, description string) (*models.Ticket, error) {
	if title == "" {
		return nil, errors.New("title is required")
	}
	if description == "" {
		return nil, errors.New("description is required")
	}

	ticket := models.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "pending",
	}

	result := s.db.Create(&ticket)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with user data
	s.db.Preload("User").First(&ticket, ticket.ID)

	// Log audit trail
	userIDInt := int(userID)
	s.auditService.LogTicketEvent(
		int(ticket.ID),
		&userIDInt,
		models.EventCreated,
		nil,
		ticket,
		nil,
		nil,
		nil,
		nil,
	)

	return &ticket, nil
}

// CreateFromRequest creates a new ticket from CreateTicketRequest
func (s *TicketService) CreateFromRequest(req models.CreateTicketRequest) (*models.Ticket, error) {
	return s.Create(req.UserID, req.Title, req.Description)
}

// CreateFromModel creates a new ticket from a models.Ticket (used for queue processing)
func (s *TicketService) CreateFromModel(ticket *models.Ticket) (*models.Ticket, error) {
	if ticket.Title == "" {
		return nil, errors.New("title is required")
	}
	if ticket.Description == "" {
		return nil, errors.New("description is required")
	}

	// Reset ID to 0 so PostgreSQL can auto-generate it
	ticket.ID = 0

	// Set default status if not provided
	if ticket.Status == "" {
		ticket.Status = "pending"
	}

	result := s.db.Create(ticket)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with user data
	s.db.Preload("User").First(ticket, ticket.ID)

	// Log audit trail
	userIDInt := int(ticket.UserID)
	s.auditService.LogTicketEvent(
		int(ticket.ID),
		&userIDInt,
		models.EventCreated,
		nil,
		ticket,
		nil,
		nil,
		nil,
		nil,
	)

	return ticket, nil
}

// UpdateStatus updates the status of a ticket
func (s *TicketService) UpdateStatus(id uint, status, reason string) (*models.Ticket, error) {
	validStatuses := []string{"pending", "accepted", "rejected"}

	// Validate status
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, errors.New("invalid status")
	}

	var ticket models.Ticket
	if err := s.db.First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	// Store old ticket state for audit
	oldTicket := ticket

	updates := map[string]interface{}{
		"status": status,
		"reason": reason,
	}

	// Set approved_at time if status is accepted
	if status == "accepted" {
		now := time.Now()
		updates["approved_at"] = now
	}

	result := s.db.Model(&ticket).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// Reload with updated data
	s.db.Preload("User").First(&ticket, id)

	// Log audit trail for status change
	changes := map[string]interface{}{
		"status": map[string]string{
			"old": string(oldTicket.Status),
			"new": string(ticket.Status),
		},
		"reason": map[string]string{
			"old": oldTicket.Reason,
			"new": ticket.Reason,
		},
	}

	action := models.EventStatusChanged
	if status == "accepted" {
		action = models.EventApproved
	} else if status == "rejected" {
		action = models.EventRejected
	}

	userIDInt := int(ticket.UserID)
	s.auditService.LogTicketEvent(
		int(ticket.ID),
		&userIDInt,
		action,
		oldTicket,
		ticket,
		changes,
		nil,
		nil,
		nil,
	)

	return &ticket, nil
}

// Update updates ticket details
func (s *TicketService) Update(id uint, title, description string) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := s.db.First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, err
	}

	// Store old ticket state for audit
	oldTicket := ticket

	updates := make(map[string]interface{})
	if title != "" {
		updates["title"] = title
	}
	if description != "" {
		updates["description"] = description
	}

	if len(updates) > 0 {
		result := s.db.Model(&ticket).Updates(updates)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	// Reload with updated data
	s.db.Preload("User").First(&ticket, id)

	// Log audit trail if there were changes
	if len(updates) > 0 {
		changes := s.auditService.CompareTickets(&oldTicket, &ticket)
		userIDInt := int(ticket.UserID)
		s.auditService.LogTicketEvent(
			int(ticket.ID),
			&userIDInt,
			models.EventUpdated,
			oldTicket,
			ticket,
			changes,
			nil,
			nil,
			nil,
		)
	}

	return &ticket, nil
}

// UpdateFromRequest updates ticket from UpdateTicketRequest
func (s *TicketService) UpdateFromRequest(id uint, req models.UpdateTicketRequest) (*models.Ticket, error) {
	return s.Update(id, req.Title, req.Description)
}

// Delete removes a ticket
func (s *TicketService) Delete(id uint) error {
	// Get ticket before deletion for audit
	var ticket models.Ticket
	if err := s.db.First(&ticket, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ticket not found")
		}
		return err
	}

	// Log audit trail before deletion
	userIDInt := int(ticket.UserID)
	s.auditService.LogTicketEvent(
		int(ticket.ID),
		&userIDInt,
		models.EventDeleted,
		ticket,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	result := s.db.Delete(&models.Ticket{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("ticket not found")
	}
	return nil
}

// GetStatistics returns statistics about tickets
func (s *TicketService) GetStatistics() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Count total tickets
	var total int64
	s.db.Model(&models.Ticket{}).Count(&total)
	stats["total"] = total

	// Count by status
	statusCount := make(map[string]int64)
	statuses := []string{"pending", "accepted", "rejected"}
	for _, status := range statuses {
		var count int64
		s.db.Model(&models.Ticket{}).Where("status = ?", status).Count(&count)
		statusCount[status] = count
	}
	stats["by_status"] = statusCount

	return stats, nil
}

// SearchTickets searches tickets by title or description
func (s *TicketService) SearchTickets(query string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	searchTerm := fmt.Sprintf("%%%s%%", query)
	result := s.db.Preload("User").Where("title ILIKE ? OR description ILIKE ?", searchTerm, searchTerm).Find(&tickets)
	return tickets, result.Error
}

// GetPendingTickets returns all pending tickets
func (s *TicketService) GetPendingTickets() ([]models.Ticket, error) {
	return s.GetByStatus("pending")
}

// ApproveTicket approves a ticket
func (s *TicketService) ApproveTicket(id uint) (*models.Ticket, error) {
	return s.UpdateStatus(id, "accepted", "")
}

// RejectTicket rejects a ticket
func (s *TicketService) RejectTicket(id uint) (*models.Ticket, error) {
	return s.UpdateStatus(id, "rejected", "")
}

// BulkUpdateStatus updates status for multiple tickets
func (s *TicketService) BulkUpdateStatus(ids []uint, status string, reason string) ([]models.Ticket, error) {
	var updatedTickets []models.Ticket
	var errors []string

	for _, id := range ids {
		ticket, err := s.UpdateStatus(id, status, reason)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Failed to update ticket %d: %s", id, err.Error()))
		} else {
			updatedTickets = append(updatedTickets, *ticket)
		}
	}

	if len(errors) > 0 {
		return updatedTickets, fmt.Errorf("some updates failed: %v", errors)
	}

	return updatedTickets, nil
}
