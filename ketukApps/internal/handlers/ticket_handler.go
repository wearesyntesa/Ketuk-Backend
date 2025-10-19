package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
)

type TicketHandler struct {
	ticketService *services.TicketService
}

func NewTicketHandler(ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// @Summary Get all tickets
// @Description Get a list of all tickets
// @Tags tickets
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /api/tickets/v1 [get]
func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	tickets, err := h.ticketService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve tickets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tickets retrieved successfully",
		Data:    tickets,
	})
}

// @Summary Get ticket by ID
// @Description Get a ticket by its ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/tickets/v1/{id} [get]
func (h *TicketHandler) GetTicketByID(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	ticket, err := h.ticketService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Ticket not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket retrieved successfully",
		Data:    ticket,
	})
}

// @Summary Get tickets by user ID
// @Description Get all tickets for a specific user
// @Tags tickets
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/user/{user_id} [get]
func (h *TicketHandler) GetTicketsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userIDInt, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
			Error:   "User ID must be a valid integer",
		})
		return
	}

	userID := uint(userIDInt)
	tickets, err := h.ticketService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "No tickets found for user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User tickets retrieved successfully",
		Data:    tickets,
	})
}

// @Summary Get tickets by status
// @Description Get all tickets with a specific status
// @Tags tickets
// @Produce json
// @Param status path string true "Ticket Status" Enums(pending, in_progress, approved, rejected, completed)
// @Success 200 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/status/{status} [get]
func (h *TicketHandler) GetTicketsByStatus(c *gin.Context) {
	status := c.Param("status")
	tickets, err := h.ticketService.GetByStatus(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve tickets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Tickets retrieved successfully",
		Data:    tickets,
	})
}

// @Summary Get pending tickets
// @Description Get all pending tickets
// @Tags tickets
// @Produce json
// @Success 200 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/pending [get]
func (h *TicketHandler) GetPendingTickets(c *gin.Context) {
	tickets, err := h.ticketService.GetPendingTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve pending tickets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Pending tickets retrieved successfully",
		Data:    tickets,
	})
}

// @Summary Search tickets
// @Description Search tickets by title or description
// @Tags tickets
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/search [get]
func (h *TicketHandler) SearchTickets(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Search query is required",
			Error:   "Query parameter 'q' cannot be empty",
		})
		return
	}

	tickets, err := h.ticketService.SearchTickets(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Search failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Search completed successfully",
		Data:    tickets,
	})
}

// @Summary Create a new ticket
// @Description Create a new ticket for a user
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body models.CreateTicketRequest true "Ticket data"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/tickets/v1 [post]
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var req models.CreateTicketRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	ticket, err := h.ticketService.CreateFromRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Ticket created successfully",
		Data:    ticket,
	})
}

// @Summary Update ticket
// @Description Update ticket information by ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param ticket body models.UpdateTicketRequest true "Updated ticket data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/tickets/v1/{id} [put]
func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	var req models.UpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	ticket, err := h.ticketService.UpdateFromRequest(id, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "ticket not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket updated successfully",
		Data:    ticket,
	})
}

// @Summary Update ticket status
// @Description Update the status of a ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Param status body UpdateStatusRequest true "New status"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/tickets/v1/{id}/status [patch]
func (h *TicketHandler) UpdateTicketStatus(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	ticket, err := h.ticketService.UpdateStatus(id, req.Status)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "ticket not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update ticket status",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket status updated successfully",
		Data:    ticket,
	})
}

// @Summary Approve ticket
// @Description Approve a ticket
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/{id}/approve [patch]
func (h *TicketHandler) ApproveTicket(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	ticket, err := h.ticketService.ApproveTicket(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "ticket not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to approve ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket approved successfully",
		Data:    ticket,
	})
}

// @Summary Reject ticket
// @Description Reject a ticket
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/{id}/reject [patch]
func (h *TicketHandler) RejectTicket(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	ticket, err := h.ticketService.RejectTicket(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "ticket not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to reject ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket rejected successfully",
		Data:    ticket,
	})
}

// @Summary Bulk update ticket status
// @Description Update status for multiple tickets
// @Tags tickets
// @Accept json
// @Produce json
// @Param request body BulkUpdateStatusRequest true "Bulk update request"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/tickets/v1/bulk-status [post]
func (h *TicketHandler) BulkUpdateStatus(c *gin.Context) {
	var req BulkUpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "At least one ticket ID is required",
			Error:   "IDs array cannot be empty",
		})
		return
	}

	// Convert []int to []uint
	uintIDs := make([]uint, len(req.IDs))
	for i, id := range req.IDs {
		uintIDs[i] = uint(id)
	}

	tickets, err := h.ticketService.BulkUpdateStatus(uintIDs, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Bulk update completed with some errors",
			Data:    tickets,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bulk status update completed successfully",
		Data:    tickets,
	})
}

// @Summary Delete ticket
// @Description Delete a ticket by ID
// @Tags tickets
// @Produce json
// @Param id path int true "Ticket ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/tickets/v1/{id} [delete]
func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	id := uint(idInt)
	err = h.ticketService.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Ticket deleted successfully",
	})
}

// @Summary Get ticket statistics
// @Description Get statistics about tickets (counts by status and category)
// @Tags tickets
// @Produce json
// @Success 200 {object} models.APIResponse
// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Router /api/tickets/v1/statistics [get]
func (h *TicketHandler) GetStatistics(c *gin.Context) {
	stats, err := h.ticketService.GetStatistics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve statistics",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Statistics retrieved successfully",
		Data:    stats,
	})
}

// Request structs for the handlers
type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type BulkUpdateStatusRequest struct {
	IDs    []int  `json:"ids" binding:"required"`
	Status string `json:"status" binding:"required"`
}
