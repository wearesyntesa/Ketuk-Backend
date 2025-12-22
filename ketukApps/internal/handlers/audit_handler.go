package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
)

type AuditHandler struct {
	auditService *services.AuditService
}

func NewAuditHandler(auditService *services.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// @Summary Get ticket event logs
// @Description Get all event logs for a specific ticket
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param ticket_id path int true "Ticket ID"
// @Success 200 {object} models.APIResponse{data=[]models.TicketEventLog}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/audit/tickets/{ticket_id}/logs [get]
func (h *AuditHandler) GetTicketEventLogs(c *gin.Context) {
	ticketIDParam := c.Param("ticket_id")
	ticketID, err := strconv.Atoi(ticketIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "Ticket ID must be a valid integer",
		})
		return
	}

	logs, err := h.auditService.GetTicketEventLogs(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve event logs",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Event logs retrieved successfully",
		Data:    logs,
	})
}

// @Summary Get event logs by user
// @Description Get all event logs created by a specific user
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.APIResponse{data=[]models.TicketEventLog}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/audit/users/{user_id}/logs [get]
func (h *AuditHandler) GetEventLogsByUser(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "User ID must be a valid integer",
		})
		return
	}

	logs, err := h.auditService.GetEventLogsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve event logs",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Event logs retrieved successfully",
		Data:    logs,
	})
}

// @Summary Get all event logs
// @Description Get all event logs in the system
// @Tags audit
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.TicketEventLog}
// @Failure 500 {object} models.APIResponse
// @Router /api/audit/logs [get]
func (h *AuditHandler) GetAllEventLogs(c *gin.Context) {
	logs, err := h.auditService.GetAllEventLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve event logs",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Event logs retrieved successfully",
		Data:    logs,
	})
}
