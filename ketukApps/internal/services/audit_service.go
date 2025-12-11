package services

import (
	"encoding/json"
	"fmt"
	"ketukApps/internal/models"

	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{
		db: db,
	}
}

// LogTicketEvent logs a ticket event to the audit trail
func (s *AuditService) LogTicketEvent(
	ticketID int,
	userID *int,
	action models.TicketEventAction,
	oldValue interface{},
	newValue interface{},
	changes map[string]interface{},
	ipAddress *string,
	userAgent *string,
	notes *string,
) error {
	var oldValueJSON, newValueJSON, changesJSON *string

	// Convert old value to JSON
	if oldValue != nil {
		data, err := json.Marshal(oldValue)
		if err != nil {
			return fmt.Errorf("failed to marshal old value: %w", err)
		}
		str := string(data)
		oldValueJSON = &str
	}

	// Convert new value to JSON
	if newValue != nil {
		data, err := json.Marshal(newValue)
		if err != nil {
			return fmt.Errorf("failed to marshal new value: %w", err)
		}
		str := string(data)
		newValueJSON = &str
	}

	// Convert changes to JSON
	if changes != nil && len(changes) > 0 {
		data, err := json.Marshal(changes)
		if err != nil {
			return fmt.Errorf("failed to marshal changes: %w", err)
		}
		str := string(data)
		changesJSON = &str
	}

	// Create event log entry
	eventLog := models.TicketEventLog{
		TicketID:  ticketID,
		UserID:    userID,
		Action:    action,
		OldValue:  oldValueJSON,
		NewValue:  newValueJSON,
		Changes:   changesJSON,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Notes:     notes,
	}

	if err := s.db.Create(&eventLog).Error; err != nil {
		return fmt.Errorf("failed to create event log: %w", err)
	}

	return nil
}

// GetTicketEventLogs retrieves all event logs for a specific ticket
func (s *AuditService) GetTicketEventLogs(ticketID int) ([]models.TicketEventLog, error) {
	var logs []models.TicketEventLog
	err := s.db.Where("ticket_id = ?", ticketID).
		Order("created_at DESC").
		Preload("User").
		Find(&logs).Error
	return logs, err
}

// GetEventLogsByUser retrieves all event logs by a specific user
func (s *AuditService) GetEventLogsByUser(userID int) ([]models.TicketEventLog, error) {
	var logs []models.TicketEventLog
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Preload("Ticket").
		Find(&logs).Error
	return logs, err
}

// CompareTickets compares two ticket objects and returns the changes
func (s *AuditService) CompareTickets(oldTicket, newTicket *models.Ticket) map[string]interface{} {
	changes := make(map[string]interface{})

	if oldTicket.Title != newTicket.Title {
		changes["title"] = map[string]string{
			"old": oldTicket.Title,
			"new": newTicket.Title,
		}
	}

	if oldTicket.Description != newTicket.Description {
		changes["description"] = map[string]string{
			"old": oldTicket.Description,
			"new": newTicket.Description,
		}
	}

	if oldTicket.Status != newTicket.Status {
		changes["status"] = map[string]string{
			"old": string(oldTicket.Status),
			"new": string(newTicket.Status),
		}
	}

	if oldTicket.Reason != newTicket.Reason {
		changes["reason"] = map[string]string{
			"old": oldTicket.Reason,
			"new": newTicket.Reason,
		}
	}

	// Compare schedule ID (nullable)
	if (oldTicket.IDSchedule == nil && newTicket.IDSchedule != nil) ||
		(oldTicket.IDSchedule != nil && newTicket.IDSchedule == nil) ||
		(oldTicket.IDSchedule != nil && newTicket.IDSchedule != nil && *oldTicket.IDSchedule != *newTicket.IDSchedule) {
		changes["idSchedule"] = map[string]interface{}{
			"old": oldTicket.IDSchedule,
			"new": newTicket.IDSchedule,
		}
	}

	return changes
}
