package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SecurityType defines WiFi security types
type SecurityType string

const (
	SecurityWPA  SecurityType = "WPA"
	SecurityWPA2 SecurityType = "WPA2"
	SecurityWEP  SecurityType = "WEP"
	SecurityNone SecurityType = "nopass"
)

// WifiCredential represents a WiFi credential with QR code
type WifiCredential struct {
	ID                uuid.UUID    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID    `gorm:"type:uuid;not null;index" json:"user_id"`
	SSID              string       `gorm:"column:ssid;not null;size:255" json:"ssid"`
	EncryptedPassword string       `gorm:"not null" json:"-"` // Never expose encrypted password
	SecurityType      SecurityType `gorm:"type:varchar(20);not null" json:"security_type"`
	IsHidden          bool         `gorm:"default:false" json:"is_hidden"`
	QRCodeData        string       `gorm:"type:text" json:"qr_code_data"` // Base64 encoded PNG
	CreatedAt         time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"user,omitempty"`
}

// BeforeCreate hook to generate UUID if not set
func (w *WifiCredential) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}

// TableName specifies the table name for WifiCredential model
func (WifiCredential) TableName() string {
	return "wifi_qr_codes"
}

// PublicWifiCredential represents WiFi credential data safe for public consumption
type PublicWifiCredential struct {
	ID           uuid.UUID    `json:"id"`
	UserID       uuid.UUID    `json:"user_id"`
	SSID         string       `json:"ssid"`
	SecurityType SecurityType `json:"security_type"`
	IsHidden     bool         `json:"is_hidden"`
	QRCodeData   string       `json:"qr_code_data"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// ToPublic converts WifiCredential to PublicWifiCredential
func (w *WifiCredential) ToPublic() *PublicWifiCredential {
	return &PublicWifiCredential{
		ID:           w.ID,
		UserID:       w.UserID,
		SSID:         w.SSID,
		SecurityType: w.SecurityType,
		IsHidden:     w.IsHidden,
		QRCodeData:   w.QRCodeData,
		CreatedAt:    w.CreatedAt,
		UpdatedAt:    w.UpdatedAt,
	}
}

// IsValidSecurityType checks if the security type is valid
func IsValidSecurityType(st string) bool {
	switch SecurityType(st) {
	case SecurityWPA, SecurityWPA2, SecurityWEP, SecurityNone:
		return true
	default:
		return false
	}
}
