import { Component, signal, inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { WiFiService } from '../../core/services/wifi.service';
import { WiFiCredential, SecurityType } from '../../core/models/wifi.model';

@Component({
  selector: 'app-qr-generator',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './qr-generator.component.html',
  styleUrls: ['./qr-generator.component.scss']
})
export class QrGeneratorComponent {
  private readonly fb = inject(FormBuilder);
  private readonly wifiService = inject(WiFiService);

  readonly wifiForm: FormGroup;
  readonly isLoading = signal(false);
  readonly errorMessage = signal<string | null>(null);
  readonly successMessage = signal<string | null>(null);
  readonly generatedCredential = signal<WiFiCredential | null>(null);

  readonly securityTypes: SecurityType[] = ['WPA', 'WPA2', 'WEP', 'nopass'];

  constructor() {
    this.wifiForm = this.fb.group({
      ssid: ['', [Validators.required, Validators.minLength(1)]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      securityType: ['WPA2', [Validators.required]],
      hidden: [false]
    });

    // Update password validation based on security type
    this.wifiForm.get('securityType')?.valueChanges.subscribe(type => {
      const passwordControl = this.wifiForm.get('password');
      if (type === 'nopass') {
        passwordControl?.clearValidators();
        passwordControl?.setValue('');
      } else {
        passwordControl?.setValidators([Validators.required, Validators.minLength(8)]);
      }
      passwordControl?.updateValueAndValidity();
    });
  }

  onSubmit(): void {
    if (this.wifiForm.invalid) {
      this.wifiForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set(null);
    this.successMessage.set(null);
    this.generatedCredential.set(null);

    this.wifiService.createWiFiCredential(this.wifiForm.value).subscribe({
      next: (credential) => {
        this.isLoading.set(false);
        this.generatedCredential.set(credential);
        this.successMessage.set('WiFi QR code generated successfully!');
        this.wifiForm.reset({ securityType: 'WPA2', hidden: false });
      },
      error: (error) => {
        this.isLoading.set(false);
        this.errorMessage.set(error.message);
      }
    });
  }

  downloadQRCode(): void {
    const credential = this.generatedCredential();
    if (!credential) return;

    const link = document.createElement('a');
    link.href = this.getQRCodeUrl();
    link.download = `wifi-qr-${credential.ssid.replace(/\s+/g, '-').toLowerCase()}.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  getQRCodeUrl(): string {
    const credential = this.generatedCredential();
    return credential ? this.wifiService.getQRCodeUrl(credential.qrCodeData) : '';
  }

  resetForm(): void {
    this.generatedCredential.set(null);
    this.successMessage.set(null);
    this.errorMessage.set(null);
  }

  get ssid() {
    return this.wifiForm.get('ssid');
  }

  get password() {
    return this.wifiForm.get('password');
  }

  get securityType() {
    return this.wifiForm.get('securityType');
  }

  get isPasswordRequired(): boolean {
    return this.securityType?.value !== 'nopass';
  }
}
