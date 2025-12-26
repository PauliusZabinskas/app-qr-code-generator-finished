import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, catchError, throwError } from 'rxjs';
import { WiFiCredential, CreateWiFiRequest, WiFiCredentialWithUser, BackendWiFiCredential } from '../models/wifi.model';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class WiFiService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = 'http://localhost:8080/api';

  createWiFiCredential(data: CreateWiFiRequest): Observable<WiFiCredential> {
    const backendData = {
      ssid: data.ssid,
      password: data.password,
      security_type: data.securityType,
      is_hidden: data.hidden
    };

    return this.http.post<BackendWiFiCredential>(`${this.apiUrl}/wifi`, backendData).pipe(
      map(cred => this.mapFromBackendCredential(cred)),
      catchError(this.handleError('Failed to create WiFi credential'))
    );
  }

  getMyCredentials(): Observable<WiFiCredential[]> {
    return this.http.get<BackendWiFiCredential[]>(`${this.apiUrl}/wifi`).pipe(
      map(creds => creds.map(c => this.mapFromBackendCredential(c))),
      catchError(this.handleError('Failed to load your WiFi credentials'))
    );
  }

  getCredentialById(id: string): Observable<WiFiCredential> {
    return this.http.get<BackendWiFiCredential>(`${this.apiUrl}/wifi/${id}`).pipe(
      map(cred => this.mapFromBackendCredential(cred)),
      catchError(this.handleError('Failed to load WiFi credential'))
    );
  }

  deleteCredential(id: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/wifi/${id}`).pipe(
      catchError(this.handleError('Failed to delete WiFi credential'))
    );
  }

  getAllCredentials(): Observable<WiFiCredentialWithUser[]> {
    return this.http.get<any[]>(`${this.apiUrl}/admin/credentials`).pipe(
      map(creds => creds.map(c => this.mapFromBackendCredential(c) as WiFiCredentialWithUser)),
      catchError(this.handleError('Failed to load all WiFi credentials'))
    );
  }

  getQRCodeUrl(qrCodeData: string): string {
    // If it's already a data URI, return it
    if (qrCodeData?.startsWith('data:')) {
      return qrCodeData;
    }
    // Otherwise assume it's valid base64 PNG data
    return `data:image/png;base64,${qrCodeData}`;
  }

  private mapFromBackendCredential(backend: BackendWiFiCredential): WiFiCredential {
    return {
      id: backend.id,
      userId: backend.user_id,
      ssid: backend.ssid,
      securityType: backend.security_type,
      hidden: backend.is_hidden,
      qrCodeData: backend.qr_code_data,
      createdAt: backend.created_at,
      updatedAt: backend.updated_at
    };
  }

  private handleError(message: string) {
    return (error: any): Observable<never> => {
      console.error(message, error);
      const errorMessage = error.error?.message || error.message || message;
      return throwError(() => new Error(errorMessage));
    };
  }
}
