package models

import "time"

// TicketEventAction defines the type of action performed on a ticket
type TicketEventAction string

const (
	EventCreated       TicketEventAction = "created"
	EventUpdated       TicketEventAction = "updated"
	EventStatusChanged TicketEventAction = "status_changed"
	EventDeleted       TicketEventAction = "deleted"
	EventAssigned      TicketEventAction = "assigned"
	EventCommented     TicketEventAction = "commented"
	EventApproved      TicketEventAction = "approved"
	EventRejected      TicketEventAction = "rejected"
)

// TicketEventLog represents the ticket_event_log table for audit trails
// @Description Ticket event log for audit trail
type TicketEventLog struct {
	ID        int               `json:"id" gorm:"primaryKey;column:id"`
	TicketID  int               `json:"ticketId" gorm:"column:ticket_id;not null"`
	UserID    *int              `json:"userId,omitempty" gorm:"column:user_id"`
	Action    TicketEventAction `json:"action" gorm:"column:action;type:ticket_event_action;not null"`
	OldValue  *string           `json:"oldValue,omitempty" gorm:"column:old_value;type:jsonb"`
	NewValue  *string           `json:"newValue,omitempty" gorm:"column:new_value;type:jsonb"`
	Changes   *string           `json:"changes,omitempty" gorm:"column:changes;type:jsonb"`
	IPAddress *string           `json:"ipAddress,omitempty" gorm:"column:ip_address;size:45"`
	UserAgent *string           `json:"userAgent,omitempty" gorm:"column:user_agent;type:text"`
	Notes     *string           `json:"notes,omitempty" gorm:"column:notes;type:text"`
	CreatedAt time.Time         `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
	Ticket    *Ticket           `json:"ticket,omitempty" gorm:"foreignKey:TicketID"`
	User      *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName overrides the table name for TicketEventLog
func (TicketEventLog) TableName() string {
	return "ticket_event_log"
}
