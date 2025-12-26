package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"gin-quickstart/internal/models"
	"gin-quickstart/internal/repositories"

	"github.com/google/uuid"
)

var (
	ErrWifiNotFound       = errors.New("WiFi credential not found")
	ErrUnauthorizedAccess = errors.New("unauthorized access to WiFi credential")
)

// WifiService handles WiFi credential business logic
type WifiService struct {
	wifiRepo      *repositories.WifiRepository
	qrCodeService *QRCodeService
	encryptionKey []byte
}

// NewWifiService creates a new WiFi service
func NewWifiService(wifiRepo *repositories.WifiRepository, qrCodeService *QRCodeService, encryptionKey string) *WifiService {
	return &WifiService{
		wifiRepo:      wifiRepo,
		qrCodeService: qrCodeService,
		encryptionKey: []byte(encryptionKey), // Must be 32 bytes for AES-256
	}
}

// CreateWifiRequest represents a request to create WiFi credential
type CreateWifiRequest struct {
	SSID         string               `json:"ssid" binding:"required,min=1,max=32"`
	Password     string               `json:"password" binding:"max=63"`
	SecurityType models.SecurityType  `json:"security_type" binding:"required"`
	IsHidden     bool                 `json:"is_hidden"`
}

// UpdateWifiRequest represents a request to update WiFi credential
type UpdateWifiRequest struct {
	SSID         string               `json:"ssid" binding:"omitempty,min=1,max=32"`
	Password     string               `json:"password" binding:"omitempty,max=63"`
	SecurityType models.SecurityType  `json:"security_type" binding:"omitempty"`
	IsHidden     *bool                `json:"is_hidden"`
}

// Create creates a new WiFi credential with QR code
func (s *WifiService) Create(userID uuid.UUID, req *CreateWifiRequest) (*models.WifiCredential, error) {
	// Validate security type
	if !models.IsValidSecurityType(string(req.SecurityType)) {
		return nil, fmt.Errorf("invalid security type: %s", req.SecurityType)
	}

	// Validate password requirement
	if req.SecurityType != models.SecurityNone && req.Password == "" {
		return nil, errors.New("password is required for secured networks")
	}

	// Encrypt password
	encryptedPassword, err := s.encryptPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Generate QR code
	qrCodeData, err := s.qrCodeService.GenerateWiFiQRCode(
		req.SSID,
		req.Password,
		req.SecurityType,
		req.IsHidden,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Create credential
	credential := &models.WifiCredential{
		UserID:            userID,
		SSID:              req.SSID,
		EncryptedPassword: encryptedPassword,
		SecurityType:      req.SecurityType,
		IsHidden:          req.IsHidden,
		QRCodeData:        qrCodeData,
	}

	if err := s.wifiRepo.Create(credential); err != nil {
		return nil, fmt.Errorf("failed to create WiFi credential: %w", err)
	}

	return credential, nil
}

// GetByID retrieves a WiFi credential by ID
func (s *WifiService) GetByID(id uuid.UUID, userID uuid.UUID, isAdmin bool) (*models.WifiCredential, error) {
	credential, err := s.wifiRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get WiFi credential: %w", err)
	}
	if credential == nil {
		return nil, ErrWifiNotFound
	}

	// Check authorization (user can only access their own credentials unless admin)
	if !isAdmin && credential.UserID != userID {
		return nil, ErrUnauthorizedAccess
	}

	return credential, nil
}

// GetAllByUser retrieves all WiFi credentials for a user
func (s *WifiService) GetAllByUser(userID uuid.UUID) ([]models.WifiCredential, error) {
	credentials, err := s.wifiRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get WiFi credentials: %w", err)
	}
	return credentials, nil
}

// GetAll retrieves all WiFi credentials (admin only)
func (s *WifiService) GetAll() ([]models.WifiCredential, error) {
	credentials, err := s.wifiRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all WiFi credentials: %w", err)
	}
	return credentials, nil
}

// Delete deletes a WiFi credential
func (s *WifiService) Delete(id uuid.UUID, userID uuid.UUID, isAdmin bool) error {
	// Get credential to check ownership
	credential, err := s.wifiRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to get WiFi credential: %w", err)
	}
	if credential == nil {
		return ErrWifiNotFound
	}

	// Check authorization
	if !isAdmin && credential.UserID != userID {
		return ErrUnauthorizedAccess
	}

	if err := s.wifiRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete WiFi credential: %w", err)
	}

	return nil
}

// encryptPassword encrypts a password using AES-256-GCM
func (s *WifiService) encryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts an encrypted password
func (s *WifiService) DecryptPassword(encryptedPassword string) (string, error) {
	if encryptedPassword == "" {
		return "", nil
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
