import { Component, signal, inject, OnInit } from '@angular/core';
import { DatePipe } from '@angular/common';
import { WiFiService } from '../../core/services/wifi.service';
import { WiFiCredential } from '../../core/models/wifi.model';

@Component({
  selector: 'app-my-codes',
  standalone: true,
  imports: [DatePipe],
  templateUrl: './my-codes.component.html',
  styleUrls: ['./my-codes.component.scss']
})
export class MyCodesComponent implements OnInit {
  private readonly wifiService = inject(WiFiService);

  readonly credentials = signal<WiFiCredential[]>([]);
  readonly selectedCredential = signal<WiFiCredential | null>(null);
  readonly isLoading = signal(false);
  readonly errorMessage = signal<string | null>(null);

  ngOnInit(): void {
    this.loadCredentials();
  }

  loadCredentials(): void {
    this.isLoading.set(true);
    this.errorMessage.set(null);

    this.wifiService.getMyCredentials().subscribe({
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

  viewQRCode(credential: WiFiCredential): void {
    this.selectedCredential.set(credential);
  }

  closeQRCode(): void {
    this.selectedCredential.set(null);
  }

  deleteCredential(id: string, event: Event): void {
    event.stopPropagation();

    if (!confirm('Are you sure you want to delete this WiFi credential?')) {
      return;
    }

    this.wifiService.deleteCredential(id).subscribe({
      next: () => {
        this.credentials.update(credentials =>
          credentials.filter(c => c.id !== id)
        );
        if (this.selectedCredential()?.id === id) {
          this.closeQRCode();
        }
      },
      error: (error) => {
        this.errorMessage.set(error.message);
      }
    });
  }

  getQRCodeUrl(qrCodeData: string): string {
    return this.wifiService.getQRCodeUrl(qrCodeData);
  }
}
