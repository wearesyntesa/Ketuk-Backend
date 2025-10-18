package services

import (
	"context"
	"ketukApps/internal/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// EnhancedTicketService extends TicketService with enhanced functionality
type EnhancedTicketService struct {
	*TicketService
}

// NewEnhancedTicketService creates a new enhanced ticket service
func NewEnhancedTicketService(db *gorm.DB) *EnhancedTicketService {
	return &EnhancedTicketService{
		TicketService: NewTicketService(db),
	}
}

// CreateWithEvents creates a ticket and logs events using new model structure
func (s *EnhancedTicketService) CreateWithEvents(ctx context.Context, userID uint, title, description string) (*models.Ticket, error) {
	// Create ticket using base service
	ticket, err := s.TicketService.Create(userID, title, description)
	if err != nil {
		return nil, err
	}

	// Log ticket created event
	log.Printf("Ticket created: ID=%d, Title='%s', User=%d",
		ticket.ID, ticket.Title, ticket.UserID)

	// Send notification to user (simplified logging for now)
	go s.sendTicketCreatedNotification(ctx, ticket)

	return ticket, nil
}

// CreateFromRequest creates a ticket using CreateTicketRequest structure and publishes events
func (s *EnhancedTicketService) CreateFromRequest(ctx context.Context, req models.CreateTicketRequest) (*models.Ticket, error) {
	return s.CreateWithEvents(ctx, req.UserID, req.Title, req.Description)
}

// UpdateStatusWithEvents updates ticket status and logs events
func (s *EnhancedTicketService) UpdateStatusWithEvents(ctx context.Context, id uint, status string) (*models.Ticket, error) {
	// Get current ticket for comparison
	currentTicket, err := s.TicketService.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update status using base service
	ticket, err := s.TicketService.UpdateStatus(id, status)
	if err != nil {
		return nil, err
	}

	// Log ticket updated event
	log.Printf("Ticket status updated: ID=%d, OldStatus='%s', NewStatus='%s'",
		ticket.ID, currentTicket.Status, status)

	// Send status change notification
	go s.sendStatusChangeNotification(ctx, ticket, currentTicket.Status)

	// If ticket is accepted, send acceptance email
	if status == "accepted" {
		go s.sendTicketAcceptedEmail(ctx, ticket)
	}

	return ticket, nil
}

// UpdateWithEvents updates ticket details and logs events
func (s *EnhancedTicketService) UpdateWithEvents(ctx context.Context, id uint, req models.UpdateTicketRequest) (*models.Ticket, error) {
	// Get current ticket for comparison
	currentTicket, err := s.TicketService.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update using base service
	ticket, err := s.TicketService.UpdateFromRequest(id, req)
	if err != nil {
		return nil, err
	}

	// Log changes
	hasChanges := false
	if req.Title != "" && req.Title != currentTicket.Title {
		log.Printf("Ticket title changed: ID=%d, Old='%s', New='%s'", ticket.ID, currentTicket.Title, req.Title)
		hasChanges = true
	}
	if req.Description != "" && req.Description != currentTicket.Description {
		log.Printf("Ticket description changed: ID=%d, Old='%s', New='%s'", ticket.ID, currentTicket.Description, req.Description)
		hasChanges = true
	}

	// Send update notification if there were changes
	if hasChanges {
		go s.sendTicketUpdatedNotification(ctx, ticket)
	}

	return ticket, nil
}

// DeleteWithEvents deletes a ticket and logs events
func (s *EnhancedTicketService) DeleteWithEvents(ctx context.Context, id uint) error {
	// Get ticket before deletion
	ticket, err := s.TicketService.GetByID(id)
	if err != nil {
		return err
	}

	// Delete using base service
	if err := s.TicketService.Delete(id); err != nil {
		return err
	}

	// Log ticket deleted event
	log.Printf("Ticket deleted: ID=%d, Title='%s'",
		ticket.ID, ticket.Title)

	// Send deletion notification
	go s.sendTicketDeletedNotification(ctx, ticket)

	return nil
}

// Notification methods (simplified with logging)
func (s *EnhancedTicketService) sendTicketCreatedNotification(ctx context.Context, ticket *models.Ticket) {
	log.Printf("NOTIFICATION: Ticket created - Ticket %d: '%s' has been created successfully",
		ticket.ID, ticket.Title)
}

func (s *EnhancedTicketService) sendStatusChangeNotification(ctx context.Context, ticket *models.Ticket, oldStatus string) {
	var title string
	switch ticket.Status {
	case "accepted":
		title = "Ticket Accepted"
	case "rejected":
		title = "Ticket Rejected"
	default:
		title = "Ticket Status Updated"
	}

	log.Printf("NOTIFICATION: %s - Ticket %d: '%s' status changed from %s to %s",
		title, ticket.ID, ticket.Title, oldStatus, ticket.Status)
}

func (s *EnhancedTicketService) sendTicketUpdatedNotification(ctx context.Context, ticket *models.Ticket) {
	log.Printf("NOTIFICATION: Ticket updated - User %d: '%s' has been updated",
		ticket.UserID, ticket.Title)
}

func (s *EnhancedTicketService) sendTicketDeletedNotification(ctx context.Context, ticket *models.Ticket) {
	log.Printf("NOTIFICATION: Ticket deleted - User %d: '%s' has been deleted",
		ticket.UserID, ticket.Title)
}

func (s *EnhancedTicketService) sendTicketAcceptedEmail(ctx context.Context, ticket *models.Ticket) {
	log.Printf("EMAIL: Ticket Accepted - Sending email to user %d for ticket '%s' (Request ID: %d)",
		ticket.UserID, ticket.Title, ticket.ID)

	approvedAtStr := "N/A"
	if ticket.ApprovedAt != nil {
		approvedAtStr = ticket.ApprovedAt.Format(time.RFC3339)
	}

	log.Printf("EMAIL CONTENT: Ticket %s accepted. Description: %s, Approved: %s",
		ticket.Title, ticket.Description, approvedAtStr)
}

// BatchProcessTickets processes multiple tickets with events
func (s *EnhancedTicketService) BatchProcessTickets(ctx context.Context, ticketIDs []uint, action string) error {
	for _, ticketID := range ticketIDs {
		switch action {
		case "accept":
			if _, err := s.UpdateStatusWithEvents(ctx, ticketID, "accepted"); err != nil {
				log.Printf("Failed to accept ticket %d: %v", ticketID, err)
			}
		case "reject":
			if _, err := s.UpdateStatusWithEvents(ctx, ticketID, "rejected"); err != nil {
				log.Printf("Failed to reject ticket %d: %v", ticketID, err)
			}
		}
	}

	return nil
}
