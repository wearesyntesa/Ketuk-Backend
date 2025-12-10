package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
)

type ScheduleHandler struct {
	scheduleService *services.ScheduleService
}

func NewScheduleHandler(scheduleService *services.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: scheduleService,
	}
}

// ScheduleTicket Handlers

// GetAllScheduleTickets handles GET /api/schedules/tickets
func (h *ScheduleHandler) GetAllScheduleTickets(c *gin.Context) {
	schedules, err := h.scheduleService.GetAllScheduleTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve schedule tickets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule tickets retrieved successfully",
		Data:    schedules,
	})
}

// GetScheduleTicketByID handles GET /api/schedules/tickets/:id
func (h *ScheduleHandler) GetScheduleTicketByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	schedule, err := h.scheduleService.GetScheduleTicketByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Schedule ticket not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule ticket retrieved successfully",
		Data:    schedule,
	})
}

// GetScheduleTicketsByUserID handles GET /api/schedules/tickets/user/:user_id
func (h *ScheduleHandler) GetScheduleTicketsByUserID(c *gin.Context) {
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

	schedules, err := h.scheduleService.GetScheduleTicketsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve user schedule tickets",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User schedule tickets retrieved successfully",
		Data:    schedules,
	})
}

// GetScheduleTicketsByCategory handles GET /api/schedules/tickets/category/:category
func (h *ScheduleHandler) GetScheduleTicketsByCategory(c *gin.Context) {
	categoryParam := c.Param("category")
	category := models.Category(categoryParam)

	schedules, err := h.scheduleService.GetScheduleTicketsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve schedule tickets by category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule tickets retrieved successfully",
		Data:    schedules,
	})
}

// CreateScheduleTicket handles POST /api/schedules/tickets
func (h *ScheduleHandler) CreateScheduleTicket(c *gin.Context) {
	var schedule models.ScheduleTicket

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	createdSchedule, err := h.scheduleService.CreateScheduleTicket(&schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create schedule ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Schedule ticket created successfully",
		Data:    createdSchedule,
	})
}

// UpdateScheduleTicket handles PUT /api/schedules/tickets/:id
func (h *ScheduleHandler) UpdateScheduleTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	schedule, err := h.scheduleService.UpdateScheduleTicket(id, updates)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "schedule ticket not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update schedule ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule ticket updated successfully",
		Data:    schedule,
	})
}

// DeleteScheduleTicket handles DELETE /api/schedules/tickets/:id
func (h *ScheduleHandler) DeleteScheduleTicket(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule ticket ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	err = h.scheduleService.DeleteScheduleTicket(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete schedule ticket",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule ticket deleted successfully",
	})
}

// ScheduleReguler Handlers

// @Summary Get all regular schedules
// @Description Get a list of all regular schedules
// @Tags schedule-reguler
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.ScheduleReguler}
// @Failure 500 {object} models.APIResponse
// @Router /api/schedules/reguler/v1 [get]
func (h *ScheduleHandler) GetAllScheduleReguler(c *gin.Context) {
	schedules, err := h.scheduleService.GetAllScheduleReguler()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve regular schedules",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Regular schedules retrieved successfully",
		Data:    schedules,
	})
}

// @Summary Get regular schedule by ID
// @Description Get a regular schedule by its ID
// @Tags schedule-reguler
// @Security BearerAuth
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.APIResponse{data=models.ScheduleReguler}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/schedules/reguler/v1/{id} [get]
func (h *ScheduleHandler) GetScheduleRegulerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule reguler ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	schedule, err := h.scheduleService.GetScheduleRegulerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Schedule reguler not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Schedule reguler retrieved successfully",
		Data:    schedule,
	})
}

// @Summary Get regular schedules by user ID
// @Description Get all regular schedules for a specific user
// @Tags schedule-reguler
// @Security BearerAuth
// @Produce json
// @Param user_id path int true "User ID"
// @Success 200 {object} models.APIResponse{data=[]models.ScheduleReguler}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/schedules/reguler/v1/user/{user_id} [get]
func (h *ScheduleHandler) GetScheduleRegulerByUserID(c *gin.Context) {
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

	schedules, err := h.scheduleService.GetScheduleRegulerByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve user regular schedules",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User regular schedules retrieved successfully",
		Data:    schedules,
	})
}

// @Summary Create a new regular schedule
// @Description Create a new regular schedule
// @Tags schedule-reguler
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param schedule body models.CreateScheduleRegulerRequest true "Schedule data"
// @Success 201 {object} models.APIResponse{data=models.ScheduleReguler}
// @Failure 400 {object} models.APIResponse
// @Router /api/schedules/reguler/v1 [post]
func (h *ScheduleHandler) CreateScheduleReguler(c *gin.Context) {
	var req models.CreateScheduleRegulerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	schedule := models.ScheduleReguler{
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		UserID:    req.UserID,
	}

	createdSchedule, err := h.scheduleService.CreateScheduleReguler(&schedule)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create regular schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Regular schedule created successfully",
		Data:    createdSchedule,
	})
}

// @Summary Update regular schedule
// @Description Update regular schedule information by ID. All fields are optional.
// @Tags schedule-reguler
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param updates body models.UpdateScheduleRegulerRequest true "Updated schedule data"
// @Success 200 {object} models.APIResponse{data=models.ScheduleReguler}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/schedules/reguler/v1/{id} [put]
func (h *ScheduleHandler) UpdateScheduleReguler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule reguler ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	schedule, err := h.scheduleService.UpdateScheduleReguler(id, updates)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "schedule reguler not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update regular schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Regular schedule updated successfully",
		Data:    schedule,
	})
}

// @Summary Delete regular schedule
// @Description Delete a regular schedule by ID
// @Tags schedule-reguler
// @Security BearerAuth
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/schedules/reguler/v1/{id} [delete]
func (h *ScheduleHandler) DeleteScheduleReguler(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid schedule reguler ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	err = h.scheduleService.DeleteScheduleReguler(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete regular schedule",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Regular schedule deleted successfully",
	})
}

// Unblocking Handlers

// GetAllUnblocking handles GET /api/unblocking
func (h *ScheduleHandler) GetAllUnblocking(c *gin.Context) {
	unblockings, err := h.scheduleService.GetAllUnblocking()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve unblocking records",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Unblocking records retrieved successfully",
		Data:    unblockings,
	})
}

// GetUnblockingByID handles GET /api/unblocking/:id
func (h *ScheduleHandler) GetUnblockingByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid unblocking ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	unblocking, err := h.scheduleService.GetUnblockingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Unblocking record not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Unblocking record retrieved successfully",
		Data:    unblocking,
	})
}

// GetUnblockingByUserID handles GET /api/unblocking/user/:user_id
func (h *ScheduleHandler) GetUnblockingByUserID(c *gin.Context) {
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

	unblockings, err := h.scheduleService.GetUnblockingByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve user unblocking records",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User unblocking records retrieved successfully",
		Data:    unblockings,
	})
}

// GetUnblockingBySemester handles GET /api/unblocking/semester/:tahun/:semester
func (h *ScheduleHandler) GetUnblockingBySemester(c *gin.Context) {
	tahunParam := c.Param("tahun")
	tahun, err := strconv.Atoi(tahunParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid tahun",
			Error:   "Tahun must be a valid integer",
		})
		return
	}

	semesterParam := c.Param("semester")
	semester := models.SemesterCategory(semesterParam)

	unblockings, err := h.scheduleService.GetUnblockingBySemester(tahun, semester)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve semester unblocking records",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Semester unblocking records retrieved successfully",
		Data:    unblockings,
	})
}

// CreateUnblocking handles POST /api/unblocking
func (h *ScheduleHandler) CreateUnblocking(c *gin.Context) {
	var unblocking models.Unblocking

	if err := c.ShouldBindJSON(&unblocking); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Set created_at to current time
	unblocking.CreatedAt = time.Now()

	createdUnblocking, err := h.scheduleService.CreateUnblocking(&unblocking)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create unblocking record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Unblocking record created successfully",
		Data:    createdUnblocking,
	})
}

// UpdateUnblocking handles PUT /api/unblocking/:id
func (h *ScheduleHandler) UpdateUnblocking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid unblocking ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	unblocking, err := h.scheduleService.UpdateUnblocking(id, updates)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "unblocking not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update unblocking record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Unblocking record updated successfully",
		Data:    unblocking,
	})
}

// DeleteUnblocking handles DELETE /api/unblocking/:id
func (h *ScheduleHandler) DeleteUnblocking(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid unblocking ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	err = h.scheduleService.DeleteUnblocking(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete unblocking record",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Unblocking record deleted successfully",
	})
}
