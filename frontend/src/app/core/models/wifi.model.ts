export type SecurityType = 'WPA' | 'WPA2' | 'WEP' | 'nopass';

export interface WiFiCredential {
  id: string;
  userId: string;
  ssid: string;
  password?: string; // Not usually returned
  securityType: SecurityType;
  hidden: boolean;
  qrCodeData: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateWiFiRequest {
  ssid: string;
  password: string;
  securityType: SecurityType;
  hidden: boolean;
}

export interface WiFiCredentialWithUser extends WiFiCredential {
  userEmail?: string;
}

// Backend DTOs to help with mapping
export interface BackendWiFiCredential {
  id: string;
  user_id: string;
  ssid: string;
  security_type: SecurityType;
  is_hidden: boolean;
  qr_code_data: string;
  created_at: string;
  updated_at: string;
}
