import { Component, signal, inject, OnInit } from '@angular/core';
import { DatePipe } from '@angular/common';
import { WiFiService } from '../../../core/services/wifi.service';
import { WiFiCredentialWithUser } from '../../../core/models/wifi.model';

@Component({
  selector: 'app-credentials',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './credentials.component.html',
  styleUrls: ['./credentials.component.scss']
})
export class CredentialsComponent implements OnInit {
  private readonly wifiService = inject(WiFiService);

  readonly credentials = signal<WiFiCredentialWithUser[]>([]);
  readonly isLoading = signal(false);
  readonly errorMessage = signal<string | null>(null);
  readonly selectedCredential = signal<WiFiCredentialWithUser | null>(null);

  ngOnInit(): void {
    this.loadAllCredentials();
  }

  loadAllCredentials(): void {
    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.wifiService.getAllCredentials().subscribe({
      next: (credentials) => {
        this.isLoading.set(false);
        this.credentials.set(credentials);
      },
      error: (error) => {
        this.isLoading.set(false);
        this.errorMessage.set(error.message);
      }
    });
  }

  viewQRCode(credential: WiFiCredentialWithUser): void {
    this.selectedCredential.set(credential);
  }

  closeQRCode(): void {
    this.selectedCredential.set(null);
  }

  getQRCodeUrl(qrCodeData: string): string {
    return this.wifiService.getQRCodeUrl(qrCodeData);
  }
}
