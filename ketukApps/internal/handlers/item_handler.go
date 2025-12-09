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

// @Summary Get all item categories
// @Description Get a list of all item categories
// @Tags item-categories
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.ItemCategory}
// @Router /api/item-categories/v1 [get]
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

// @Summary Get item category by ID
// @Description Get an item category by its ID
// @Tags item-categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.APIResponse{data=models.ItemCategory}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/item-categories/v1/{id} [get]
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

// @Summary Create a new item category
// @Description Create a new item category
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body models.CreateItemCategoryRequest true "Item category data"
// @Success 201 {object} models.APIResponse{data=models.ItemCategory}
// @Failure 400 {object} models.APIResponse
// @Router /api/item-categories/v1 [post]
func (h *ItemHandler) CreateItemCategory(c *gin.Context) {
	var category models.CreateItemCategoryRequest

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	data := models.ItemCategory{
		CategoryName:  category.CategoryName,
		Specification: category.Specification,
	}

	createdCategory, err := h.itemService.CreateItemCategory(&data)
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

// @Summary Update item category
// @Description Update item category information by ID. All fields are optional.
// @Tags item-categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param updates body map[string]interface{} true "Updated category data (categoryName, specification)"
// @Success 200 {object} models.APIResponse{data=models.ItemCategory}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/item-categories/v1/{id} [put]
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

// @Summary Delete item category
// @Description Delete an item category by ID
// @Tags item-categories
// @Security BearerAuth
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/item-categories/v1/{id} [delete]
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

// @Summary Get all items
// @Description Get a list of all items
// @Tags items
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Item}
// @Router /api/items/v1 [get]
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

// @Summary Get item by ID
// @Description Get an item by its ID
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.APIResponse{data=models.Item}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/items/v1/{id} [get]
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

// @Summary Get items by category ID
// @Description Get all items for a specific category
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param category_id path int true "Category ID"
// @Success 200 {object} models.APIResponse{data=[]models.Item}
// @Failure 400 {object} models.APIResponse
// @Router /api/items/v1/category/{category_id} [get]
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

// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Summary Search items
// //@Description Search items by name or note
// //@Tags items
// //@Produce json
// //@Param q query string true "Search query"
// //@Success 200 {object} models.APIResponse
// //@Router /api/items/v1/search [get]
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

// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Summary Get items by kondisi
// //@Description Get all items with a specific condition
// //@Tags items
// //@Produce json
// //@Param kondisi path string true "Item Kondisi (Baik, Rusak Ringan, Rusak Berat)"
// //@Success 200 {object} models.APIResponse
// //@Router /api/items/v1/kondisi/{kondisi} [get]
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

// @Summary Create a new item
// @Description Create a new item
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param item body models.CreateItemRequest true "Item data"
// @Success 201 {object} models.APIResponse{data=models.Item}
// @Failure 400 {object} models.APIResponse
// @Router /api/items/v1 [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item models.CreateItemRequest

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	itemCreate := models.Item{
		Name:       item.Name,
		Year:       item.Year,
		Kondisi:    item.Kondisi,
		Note:       item.Note,
		CategoryID: item.CategoryID,
	}

	createdItem, err := h.itemService.CreateItem(&itemCreate)
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

// @Summary Update item
// @Description Update item information by ID. All fields are optional.
// @Tags items
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param updates body map[string]interface{} true "Updated item data (name, year, kondisi, note, categoryId)"
// @Success 200 {object} models.APIResponse{data=models.Item}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/items/v1/{id} [put]
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

// @Summary Delete item
// @Description Delete an item by ID
// @Tags items
// @Security BearerAuth
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/items/v1/{id} [delete]
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

// NOTE: This endpoint is not registered in the router - comment out to hide from Swagger
// //@Summary Get item statistics
// //@Description Get statistics about items (counts by category and kondisi)
// //@Tags items
// //@Produce json
// //@Success 200 {object} models.APIResponse
// //@Router /api/items/v1/statistics [get]
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
