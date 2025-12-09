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

// CreateUnblocking godoc
// @Summary Create a new unblocking request
// @Description Create a new unblocking request for semester unblocking
// @Tags unblocking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param unblocking body models.CreateUnblockingRequest true "Unblocking data"
// @Success 201 {object} models.UnblockingResponse{unblocking=models.Unblocking}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/unblockings/v1 [post]
func (h *UnblockingHandler) CreateUnblocking(c *gin.Context) {
	var req models.CreateUnblockingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	unblockingData := models.Unblocking{
		Tahun:     req.Tahun,
		Semester:  req.Semester,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		UserID:    req.UserID,
	}


	unblocking, err := h.unblockingService.Create(&unblockingData)
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

// GetUnblockingByID godoc
// @Summary Get unblocking request by ID
// @Description Get an unblocking request by its ID
// @Tags unblocking
// @Produce json
// @Security BearerAuth
// @Param id path int true "Unblocking ID"
// @Success 200 {object} models.UnblockingResponse{unblocking=models.Unblocking}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/unblockings/v1/{id} [get]
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

// GetUnblockingsByUserID godoc
// @Summary Get unblocking requests by user ID
// @Description Get all unblocking requests for a specific user
// @Tags unblocking
// @Produce json
// @Security BearerAuth
// @Param user_id path int true "User ID"
// @Success 200 {object} models.UnblockingsResponse{unblockings=[]models.Unblocking}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/unblockings/v1/user/{user_id} [get]
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

// GetAllUnblockings godoc
// @Summary Get all unblocking requests
// @Description Get a list of all unblocking requests
// @Tags unblocking
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UnblockingsResponse{unblockings=[]models.Unblocking}
// @Failure 500 {object} models.APIResponse
// @Router /api/unblockings/v1 [get]
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
