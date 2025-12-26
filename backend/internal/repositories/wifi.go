package repositories

import (
	"errors"
	"fmt"

	"gin-quickstart/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// WifiRepository handles database operations for WiFi credentials
type WifiRepository struct {
	db *gorm.DB
}

// NewWifiRepository creates a new WiFi repository
func NewWifiRepository(db *gorm.DB) *WifiRepository {
	return &WifiRepository{db: db}
}

// Create creates a new WiFi credential
func (r *WifiRepository) Create(credential *models.WifiCredential) error {
	if err := r.db.Create(credential).Error; err != nil {
		return fmt.Errorf("failed to create WiFi credential: %w", err)
	}
	return nil
}

// FindByID finds a WiFi credential by ID
func (r *WifiRepository) FindByID(id uuid.UUID) (*models.WifiCredential, error) {
	var credential models.WifiCredential
	err := r.db.First(&credential, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find WiFi credential by ID: %w", err)
	}
	return &credential, nil
}

// FindByUserID retrieves all WiFi credentials for a specific user
func (r *WifiRepository) FindByUserID(userID uuid.UUID) ([]models.WifiCredential, error) {
	var credentials []models.WifiCredential
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&credentials).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find WiFi credentials by user ID: %w", err)
	}
	return credentials, nil
}

// GetAll retrieves all WiFi credentials (admin functionality)
func (r *WifiRepository) GetAll() ([]models.WifiCredential, error) {
	var credentials []models.WifiCredential
	err := r.db.Preload("User").
		Order("created_at DESC").
		Find(&credentials).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all WiFi credentials: %w", err)
	}
	return credentials, nil
}

// Update updates a WiFi credential
func (r *WifiRepository) Update(credential *models.WifiCredential) error {
	if err := r.db.Save(credential).Error; err != nil {
		return fmt.Errorf("failed to update WiFi credential: %w", err)
	}
	return nil
}

// Delete deletes a WiFi credential
func (r *WifiRepository) Delete(id uuid.UUID) error {
	result := r.db.Delete(&models.WifiCredential{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete WiFi credential: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteByUserID deletes all WiFi credentials for a user
func (r *WifiRepository) DeleteByUserID(userID uuid.UUID) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&models.WifiCredential{}).Error; err != nil {
		return fmt.Errorf("failed to delete WiFi credentials by user ID: %w", err)
	}
	return nil
}

// Count returns the total number of WiFi credentials
func (r *WifiRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.WifiCredential{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count WiFi credentials: %w", err)
	}
	return count, nil
}

// CountByUserID returns the number of WiFi credentials for a user
func (r *WifiRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.WifiCredential{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count WiFi credentials by user ID: %w", err)
	}
	return count, nil
}
