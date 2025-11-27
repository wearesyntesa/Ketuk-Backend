package handlers

import (
	"net/http"
	"strconv"

	"ketukApps/internal/models"
	"ketukApps/internal/services"

	"github.com/gin-gonic/gin"
)

type UnblockingHandler struct {
	unblockingService *services.UnblockingService
}

func NewUnblockingHandler(unblockingService *services.UnblockingService) *UnblockingHandler {
	return &UnblockingHandler{
		unblockingService: unblockingService,
	}
}

// CreateUnblocking handles the creation of a new unblocking request
func (h *UnblockingHandler) CreateUnblocking(c *gin.Context) {
	var req models.Unblocking
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unblocking, err := h.unblockingService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.UnblockingResponse{
		Success:    true,
		Message:    "Unblocking operation completed successfully",
		Unblocking: *unblocking,
	})
}

// GetUnblockingByID handles fetching an unblocking request by its ID
func (h *UnblockingHandler) GetUnblockingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid unblocking ID"})
		return
	}

	unblocking, err := h.unblockingService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UnblockingResponse{
		Success:    true,
		Message:    "Unblocking operation completed successfully",
		Unblocking: *unblocking,
	})
}

// GetUnblockingsByUserID handles fetching all unblocking requests for a specific user
func (h *UnblockingHandler) GetUnblockingsByUserID(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	unblockings, err := h.unblockingService.GetByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UnblockingsResponse{
		Success:     true,
		Message:     "Unblocking operation completed successfully",
		Unblockings: unblockings,
	})
}

// GetAllUnblockings handles fetching all unblocking requests
func (h *UnblockingHandler) GetAllUnblockings(c *gin.Context) {
	unblockings, err := h.unblockingService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.UnblockingsResponse{
		Success:     true,
		Message:     "Unblocking operation completed successfully",
		Unblockings: unblockings,
	})
}
