package handlers

import (
	"net/http"

	"gin-quickstart/internal/repositories"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles admin-only endpoints
type AdminHandler struct {
	userRepo *repositories.UserRepository
	wifiRepo *repositories.WifiRepository
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(userRepo *repositories.UserRepository, wifiRepo *repositories.WifiRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
		wifiRepo: wifiRepo,
	}
}

// GetAllUsers handles retrieving all users (admin only)
// @Summary Get all users
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PublicUser
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/admin/users [get]
func (h *AdminHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve users",
			Message: err.Error(),
		})
		return
	}

	// Convert to public format
	publicUsers := make([]interface{}, 0, len(users))
	for _, user := range users {
		publicUsers = append(publicUsers, user.ToPublic())
	}

	c.JSON(http.StatusOK, publicUsers)
}

// GetAllCredentials handles retrieving all WiFi credentials (admin only)
// @Summary Get all WiFi credentials
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PublicWifiCredential
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/admin/credentials [get]
func (h *AdminHandler) GetAllCredentials(c *gin.Context) {
	credentials, err := h.wifiRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve WiFi credentials",
			Message: err.Error(),
		})
		return
	}

	// Convert to public format with user info
	type CredentialWithUser struct {
		ID           string `json:"id"`
		UserID       string `json:"user_id"`
		UserEmail    string `json:"user_email"`
		SSID         string `json:"ssid"`
		SecurityType string `json:"security_type"`
		IsHidden     bool   `json:"is_hidden"`
		CreatedAt    string `json:"created_at"`
	}

	publicCredentials := make([]CredentialWithUser, 0, len(credentials))
	for _, cred := range credentials {
		userEmail := ""
		if cred.User != nil {
			userEmail = cred.User.Email
		}

		publicCredentials = append(publicCredentials, CredentialWithUser{
			ID:           cred.ID.String(),
			UserID:       cred.UserID.String(),
			UserEmail:    userEmail,
			SSID:         cred.SSID,
			SecurityType: string(cred.SecurityType),
			IsHidden:     cred.IsHidden,
			CreatedAt:    cred.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, publicCredentials)
}

// GetStats handles retrieving system statistics (admin only)
// @Summary Get system statistics
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /api/admin/stats [get]
func (h *AdminHandler) GetStats(c *gin.Context) {
	// Get user count
	users, err := h.userRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve statistics",
			Message: err.Error(),
		})
		return
	}

	// Get credential count
	credentialCount, err := h.wifiRepo.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to retrieve statistics",
			Message: err.Error(),
		})
		return
	}

	stats := StatsResponse{
		TotalUsers:       len(users),
		TotalCredentials: int(credentialCount),
	}

	c.JSON(http.StatusOK, stats)
}

// StatsResponse represents system statistics
type StatsResponse struct {
	TotalUsers       int `json:"total_users"`
	TotalCredentials int `json:"total_credentials"`
}
