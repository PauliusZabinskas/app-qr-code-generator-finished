package handlers

import (
	"errors"
	"net/http"

	"gin-quickstart/internal/middleware"
	"gin-quickstart/internal/models"
	"gin-quickstart/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WifiHandler handles WiFi credential endpoints
type WifiHandler struct {
	wifiService *services.WifiService
}

// NewWifiHandler creates a new WiFi handler
func NewWifiHandler(wifiService *services.WifiService) *WifiHandler {
	return &WifiHandler{
		wifiService: wifiService,
	}
}

// Create handles creating a new WiFi credential
// @Summary Create WiFi credential
// @Tags wifi
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateWifiRequest true "WiFi credential details"
// @Success 201 {object} models.PublicWifiCredential
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/wifi [post]
func (h *WifiHandler) Create(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	var req services.CreateWifiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	credential, err := h.wifiService.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create WiFi credential",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, credential.ToPublic())
}

// GetAll handles retrieving all WiFi credentials for the current user
// @Summary Get user's WiFi credentials
// @Tags wifi
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PublicWifiCredential
// @Failure 401 {object} ErrorResponse
// @Router /api/wifi [get]
func (h *WifiHandler) GetAll(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	credentials, err := h.wifiService.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve WiFi credentials",
			Message: err.Error(),
		})
		return
	}

	// Convert to public format
	// Convert to public format
	publicCredentials := make([]*models.PublicWifiCredential, 0, len(credentials))
	for _, cred := range credentials {
		publicCredentials = append(publicCredentials, cred.ToPublic())
	}

	c.JSON(http.StatusOK, publicCredentials)
}

// GetByID handles retrieving a specific WiFi credential
// @Summary Get WiFi credential by ID
// @Tags wifi
// @Produce json
// @Security BearerAuth
// @Param id path string true "WiFi credential ID"
// @Success 200 {object} models.PublicWifiCredential
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/wifi/{id} [get]
func (h *WifiHandler) GetByID(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)

	credential, err := h.wifiService.GetByID(id, userID, isAdmin)
	if err != nil {
		if errors.Is(err, services.ErrWifiNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "WiFi credential not found",
			})
			return
		}
		if errors.Is(err, services.ErrUnauthorizedAccess) {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to access this WiFi credential",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve WiFi credential",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, credential.ToPublic())
}

// Delete handles deleting a WiFi credential
// @Summary Delete WiFi credential
// @Tags wifi
// @Security BearerAuth
// @Param id path string true "WiFi credential ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/wifi/{id} [delete]
func (h *WifiHandler) Delete(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	isAdmin := middleware.IsAdmin(c)

	if err := h.wifiService.Delete(id, userID, isAdmin); err != nil {
		if errors.Is(err, services.ErrWifiNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "WiFi credential not found",
			})
			return
		}
		if errors.Is(err, services.ErrUnauthorizedAccess) {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "You don't have permission to delete this WiFi credential",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete WiFi credential",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}
