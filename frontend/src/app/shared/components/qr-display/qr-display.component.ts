import { Component, input } from '@angular/core';

@Component({
  selector: 'app-qr-display',
  standalone: true,
  imports: [],
  templateUrl: './qr-display.component.html',
  styleUrls: ['./qr-display.component.scss']
})
export class QrDisplayComponent {
  readonly qrCodeUrl = input.required<string>();
  readonly ssid = input.required<string>();
  readonly securityType = input<string>('');
}
