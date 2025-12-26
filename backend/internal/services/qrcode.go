package services

import (
	"encoding/base64"
	"fmt"

	"gin-quickstart/internal/models"

	qrcode "github.com/skip2/go-qrcode"
)

// QRCodeService handles QR code generation
type QRCodeService struct{}

// NewQRCodeService creates a new QR code service
func NewQRCodeService() *QRCodeService {
	return &QRCodeService{}
}

// GenerateWiFiQRCode generates a QR code for WiFi credentials
// Format: WIFI:T:<security>;S:<ssid>;P:<password>;H:<hidden>;;
func (s *QRCodeService) GenerateWiFiQRCode(ssid string, password string, security models.SecurityType, hidden bool) (string, error) {
	// Build WiFi QR code string according to specification
	// Reference: https://github.com/zxing/zxing/wiki/Barcode-Contents#wi-fi-network-config-android-ios-11
	wifiString := s.buildWiFiString(ssid, password, security, hidden)

	// Generate QR code as PNG bytes
	// Using 256x256 pixels with medium error correction
	pngBytes, err := qrcode.Encode(wifiString, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Encode to base64 for easy storage and transmission
	base64String := base64.StdEncoding.EncodeToString(pngBytes)

	return base64String, nil
}

// buildWiFiString constructs the WiFi configuration string for QR code
func (s *QRCodeService) buildWiFiString(ssid string, password string, security models.SecurityType, hidden bool) string {
	// Escape special characters in SSID and password
	escapedSSID := escapeWiFiString(ssid)
	escapedPassword := escapeWiFiString(password)

	// Hidden flag: "true" if hidden, empty otherwise
	hiddenFlag := ""
	if hidden {
		hiddenFlag = "true"
	}

	// Security type mapping
	securityStr := string(security)
	if security == models.SecurityNone {
		securityStr = "nopass"
		escapedPassword = "" // No password for open networks
	}

	// Build the WiFi string
	// Format: WIFI:T:<security>;S:<ssid>;P:<password>;H:<hidden>;;
	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;H:%s;;", securityStr, escapedSSID, escapedPassword, hiddenFlag)
}

// escapeWiFiString escapes special characters in WiFi QR code strings
// Special characters that need escaping: \ ; , : "
func escapeWiFiString(s string) string {
	result := ""
	for _, char := range s {
		switch char {
		case '\\', ';', ',', ':', '"':
			result += "\\" + string(char)
		default:
			result += string(char)
		}
	}
	return result
}

// ValidateQRCodeData validates that QR code data is properly formatted base64
func (s *QRCodeService) ValidateQRCodeData(data string) error {
	_, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("invalid base64 QR code data: %w", err)
	}
	return nil
}
