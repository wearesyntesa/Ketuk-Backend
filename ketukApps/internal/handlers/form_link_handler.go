package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
)

type FormLinkHandler struct {
	formLinkService *services.FormLinkService
}

func NewFormLinkHandler(formLinkService *services.FormLinkService) *FormLinkHandler {
	return &FormLinkHandler{
		formLinkService: formLinkService,
	}
}

// @Summary Create a new form link
// @Description Create a new shareable form link for public schedule requests
// @Tags form-links
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param formLink body models.CreateFormLinkRequest true "Form link data"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/form-links/v1 [post]
func (h *FormLinkHandler) CreateFormLink(c *gin.Context) {
	var req models.CreateFormLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get the user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not found in context",
		})
		return
	}

	userModel := user.(models.User)

	formLink, err := h.formLinkService.Create(req, userModel.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create form link",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Form link created successfully",
		Data:    formLink,
	})
}

// @Summary Get all form links
// @Description Get a list of all form links
// @Tags form-links
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse
// @Router /api/form-links/v1 [get]
func (h *FormLinkHandler) GetAllFormLinks(c *gin.Context) {
	formLinks, err := h.formLinkService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve form links",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form links retrieved successfully",
		Data:    formLinks,
	})
}

// @Summary Get form link by ID
// @Description Get a form link by its ID
// @Tags form-links
// @Security BearerAuth
// @Produce json
// @Param id path int true "Form Link ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id} [get]
func (h *FormLinkHandler) GetFormLinkByID(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	formLink, err := h.formLinkService.GetByID(uint(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Form link not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form link retrieved successfully",
		Data:    formLink,
	})
}

// @Summary Update form link
// @Description Update a form link by ID
// @Tags form-links
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Form Link ID"
// @Param formLink body models.UpdateFormLinkRequest true "Updated form link data"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id} [put]
func (h *FormLinkHandler) UpdateFormLink(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var req models.UpdateFormLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	formLink, err := h.formLinkService.Update(uint(idInt), req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "form link not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update form link",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form link updated successfully",
		Data:    formLink,
	})
}

// @Summary Delete form link
// @Description Delete a form link by ID
// @Tags form-links
// @Security BearerAuth
// @Produce json
// @Param id path int true "Form Link ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id} [delete]
func (h *FormLinkHandler) DeleteFormLink(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	if err := h.formLinkService.Delete(uint(idInt)); err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete form link",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form link deleted successfully",
	})
}

// @Summary Deactivate form link
// @Description Deactivate a form link (disable submissions)
// @Tags form-links
// @Security BearerAuth
// @Produce json
// @Param id path int true "Form Link ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id}/deactivate [patch]
func (h *FormLinkHandler) DeactivateFormLink(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	formLink, err := h.formLinkService.Deactivate(uint(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to deactivate form link",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form link deactivated successfully",
		Data:    formLink,
	})
}

// @Summary Clone form link
// @Description Clone an existing form link with a new code and expiration
// @Tags form-links
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Form Link ID to clone"
// @Param request body CloneFormLinkRequest true "Clone request with new expiration"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id}/clone [post]
func (h *FormLinkHandler) CloneFormLink(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	var req CloneFormLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Get the user from context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "User not found in context",
		})
		return
	}

	userModel := user.(models.User)

	formLink, err := h.formLinkService.Clone(uint(idInt), req.ExpiresAt, userModel.ID)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "form link not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to clone form link",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Form link cloned successfully",
		Data:    formLink,
	})
}

// @Summary Get public form by code
// @Description Get public form configuration by code (no authentication required)
// @Tags public
// @Produce json
// @Param code path string true "Form Link Code"
// @Success 200 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/public/form/{code} [get]
func (h *FormLinkHandler) GetPublicForm(c *gin.Context) {
	code := c.Param("code")

	formResponse, err := h.formLinkService.GetPublicFormByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Form not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Form retrieved successfully",
		Data:    formResponse,
	})
}

// @Summary Submit public form
// @Description Submit a schedule request via public form (no authentication required)
// @Tags public
// @Accept json
// @Produce json
// @Param code path string true "Form Link Code"
// @Param request body models.PublicFormSubmitRequest true "Form submission data"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/public/form/{code}/submit [post]
func (h *FormLinkHandler) SubmitPublicForm(c *gin.Context) {
	code := c.Param("code")

	var req models.PublicFormSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.formLinkService.SubmitPublicForm(code, req)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "form link not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to submit form",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Form submitted successfully",
		Data:    response,
	})
}

// @Summary Get submissions by form link
// @Description Get all submissions for a specific form link
// @Tags form-links
// @Security BearerAuth
// @Produce json
// @Param id path int true "Form Link ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/form-links/v1/{id}/submissions [get]
func (h *FormLinkHandler) GetSubmissions(c *gin.Context) {
	idParam := c.Param("id")
	idInt, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid form link ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	// First check if form link exists
	_, err = h.formLinkService.GetByID(uint(idInt))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Form link not found",
			Error:   err.Error(),
		})
		return
	}

	tickets, err := h.formLinkService.GetSubmissionsByFormLinkID(uint(idInt))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve submissions",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Submissions retrieved successfully",
		Data:    tickets,
	})
}

// CloneFormLinkRequest is the request body for cloning a form link
type CloneFormLinkRequest struct {
	ExpiresAt time.Time `json:"expiresAt" binding:"required"`
}
