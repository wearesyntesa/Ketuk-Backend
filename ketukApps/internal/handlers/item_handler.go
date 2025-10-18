package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ketukApps/internal/models"
	"ketukApps/internal/services"
)

type ItemHandler struct {
	itemService *services.ItemService
}

func NewItemHandler(itemService *services.ItemService) *ItemHandler {
	return &ItemHandler{
		itemService: itemService,
	}
}

// ItemCategory Handlers

// GetAllItemCategories handles GET /api/items/categories
func (h *ItemHandler) GetAllItemCategories(c *gin.Context) {
	categories, err := h.itemService.GetAllItemCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve item categories",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item categories retrieved successfully",
		Data:    categories,
	})
}

// GetItemCategoryByID handles GET /api/items/categories/:id
func (h *ItemHandler) GetItemCategoryByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	category, err := h.itemService.GetItemCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Item category not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item category retrieved successfully",
		Data:    category,
	})
}

// CreateItemCategory handles POST /api/items/categories
func (h *ItemHandler) CreateItemCategory(c *gin.Context) {
	var category models.ItemCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	createdCategory, err := h.itemService.CreateItemCategory(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create item category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Item category created successfully",
		Data:    createdCategory,
	})
}

// UpdateItemCategory handles PUT /api/items/categories/:id
func (h *ItemHandler) UpdateItemCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
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

	category, err := h.itemService.UpdateItemCategory(id, updates)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "item category not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update item category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item category updated successfully",
		Data:    category,
	})
}

// DeleteItemCategory handles DELETE /api/items/categories/:id
func (h *ItemHandler) DeleteItemCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	err = h.itemService.DeleteItemCategory(id)
	if err != nil {
		status := http.StatusNotFound
		if err.Error() == "cannot delete category with existing items" {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to delete item category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item category deleted successfully",
	})
}

// Item Handlers

// GetAllItems handles GET /api/items
func (h *ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.itemService.GetAllItems()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve items",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Data:    items,
	})
}

// GetItemByID handles GET /api/items/:id
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid item ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	item, err := h.itemService.GetItemByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Item not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item retrieved successfully",
		Data:    item,
	})
}

// GetItemsByCategoryID handles GET /api/items/category/:category_id
func (h *ItemHandler) GetItemsByCategoryID(c *gin.Context) {
	categoryIDParam := c.Param("category_id")
	categoryID, err := strconv.Atoi(categoryIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   "Category ID must be a valid integer",
		})
		return
	}

	items, err := h.itemService.GetItemsByCategoryID(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve items for category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Data:    items,
	})
}

// SearchItems handles GET /api/items/search?q=query
func (h *ItemHandler) SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Search query is required",
			Error:   "Query parameter 'q' cannot be empty",
		})
		return
	}

	items, err := h.itemService.SearchItems(query)
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
		Data:    items,
	})
}

// GetItemsByKondisi handles GET /api/items/kondisi/:kondisi
func (h *ItemHandler) GetItemsByKondisi(c *gin.Context) {
	kondisi := c.Param("kondisi")

	items, err := h.itemService.GetItemsByKondisi(kondisi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve items by condition",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Items retrieved successfully",
		Data:    items,
	})
}

// CreateItem handles POST /api/items
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	createdItem, err := h.itemService.CreateItem(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Failed to create item",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Item created successfully",
		Data:    createdItem,
	})
}

// UpdateItem handles PUT /api/items/:id
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid item ID",
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

	item, err := h.itemService.UpdateItem(id, updates)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "item not found" {
			status = http.StatusNotFound
		}

		c.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update item",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item updated successfully",
		Data:    item,
	})
}

// DeleteItem handles DELETE /api/items/:id
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid item ID",
			Error:   "ID must be a valid integer",
		})
		return
	}

	err = h.itemService.DeleteItem(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Failed to delete item",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Item deleted successfully",
	})
}

// GetItemStatistics handles GET /api/items/statistics
func (h *ItemHandler) GetItemStatistics(c *gin.Context) {
	stats, err := h.itemService.GetItemStatistics()
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
