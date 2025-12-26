# WiFi QR Code Generator - Frontend Specification

## Overview

This document provides detailed specifications for the Angular 21+ frontend application, including component architecture, routing, state management, and UI/UX guidelines.

**Framework**: Angular 21+
**Language**: TypeScript 5.3+
**Styling**: Tailwind CSS 3.4+
**Architecture**: Standalone Components with Signals

---

## Table of Contents

1. [Project Structure](#project-structure)
2. [Routing & Navigation](#routing--navigation)
3. [Component Specifications](#component-specifications)
4. [State Management](#state-management)
5. [Services](#services)
6. [Guards & Interceptors](#guards--interceptors)
7. [Forms & Validation](#forms--validation)
8. [UI/UX Guidelines](#uiux-guidelines)
9. [Performance Optimization](#performance-optimization)

---

## Project Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── core/                          # Core functionality (singleton services, guards, interceptors)
│   │   │   ├── guards/
│   │   │   │   ├── auth.guard.ts
│   │   │   │   └── admin.guard.ts
│   │   │   ├── interceptors/
│   │   │   │   ├── jwt.interceptor.ts
│   │   │   │   └── error.interceptor.ts
│   │   │   ├── services/
│   │   │   │   ├── auth.service.ts
│   │   │   │   ├── qr-code.service.ts
│   │   │   │   └── storage.service.ts
│   │   │   ├── models/
│   │   │   │   ├── user.model.ts
│   │   │   │   ├── qr-code.model.ts
│   │   │   │   └── api-response.model.ts
│   │   │   └── stores/
│   │   │       ├── auth.store.ts
│   │   │       └── qr-code.store.ts
│   │   ├── features/                      # Feature modules (lazy-loaded)
│   │   │   ├── auth/
│   │   │   │   ├── login/
│   │   │   │   │   ├── login.component.ts
│   │   │   │   │   ├── login.component.html
│   │   │   │   │   └── login.component.css
│   │   │   │   └── register/
│   │   │   │       ├── register.component.ts
│   │   │   │       ├── register.component.html
│   │   │   │       └── register.component.css
│   │   │   ├── dashboard/
│   │   │   │   ├── dashboard.component.ts
│   │   │   │   ├── dashboard.component.html
│   │   │   │   └── dashboard.component.css
│   │   │   ├── qr-generator/
│   │   │   │   ├── qr-generator.component.ts
│   │   │   │   ├── qr-generator.component.html
│   │   │   │   ├── qr-generator.component.css
│   │   │   │   └── components/
│   │   │   │       ├── qr-form/
│   │   │   │       │   ├── qr-form.component.ts
│   │   │   │       │   ├── qr-form.component.html
│   │   │   │       │   └── qr-form.component.css
│   │   │   │       └── qr-display/
│   │   │   │           ├── qr-display.component.ts
│   │   │   │           ├── qr-display.component.html
│   │   │   │           └── qr-display.component.css
│   │   │   ├── my-codes/
│   │   │   │   ├── my-codes.component.ts
│   │   │   │   ├── my-codes.component.html
│   │   │   │   ├── my-codes.component.css
│   │   │   │   └── components/
│   │   │   │       └── qr-card/
│   │   │   │           ├── qr-card.component.ts
│   │   │   │           ├── qr-card.component.html
│   │   │   │           └── qr-card.component.css
│   │   │   └── admin/
│   │   │       ├── credentials/
│   │   │       │   ├── credentials.component.ts
│   │   │       │   ├── credentials.component.html
│   │   │       │   └── credentials.component.css
│   │   │       └── components/
│   │   │           └── credential-table/
│   │   │               ├── credential-table.component.ts
│   │   │               ├── credential-table.component.html
│   │   │               └── credential-table.component.css
│   │   ├── shared/                        # Shared components, directives, pipes
│   │   │   ├── components/
│   │   │   │   ├── navbar/
│   │   │   │   │   ├── navbar.component.ts
│   │   │   │   │   ├── navbar.component.html
│   │   │   │   │   └── navbar.component.css
│   │   │   │   ├── loading-spinner/
│   │   │   │   │   ├── loading-spinner.component.ts
│   │   │   │   │   ├── loading-spinner.component.html
│   │   │   │   │   └── loading-spinner.component.css
│   │   │   │   └── error-message/
│   │   │   │       ├── error-message.component.ts
│   │   │   │       ├── error-message.component.html
│   │   │   │       └── error-message.component.css
│   │   │   ├── directives/
│   │   │   │   └── click-outside.directive.ts
│   │   │   └── pipes/
│   │   │       └── date-format.pipe.ts
│   │   ├── app.component.ts               # Root component
│   │   ├── app.component.html
│   │   ├── app.component.css
│   │   ├── app.config.ts                  # Application configuration
│   │   └── app.routes.ts                  # Route definitions
│   ├── assets/
│   │   ├── images/
│   │   └── icons/
│   ├── environments/
│   │   ├── environment.ts
│   │   └── environment.development.ts
│   ├── index.html
│   ├── main.ts
│   └── styles.css                         # Global Tailwind imports
├── tailwind.config.js
├── angular.json
├── package.json
├── tsconfig.json
└── tsconfig.app.json
```

---

## Routing & Navigation

### Route Configuration

```typescript
// app.routes.ts
import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';
import { AdminGuard } from './core/guards/admin.guard';

export const routes: Routes = [
  // Public routes
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login/login.component')
      .then(m => m.LoginComponent),
    title: 'Login - WiFi QR Code Generator'
  },
  {
    path: 'register',
    loadComponent: () => import('./features/auth/register/register.component')
      .then(m => m.RegisterComponent),
    title: 'Register - WiFi QR Code Generator'
  },

  // Protected routes (requires authentication)
  {
    path: 'dashboard',
    loadComponent: () => import('./features/dashboard/dashboard.component')
      .then(m => m.DashboardComponent),
    canActivate: [AuthGuard],
    title: 'Dashboard - WiFi QR Code Generator'
  },
  {
    path: 'qr-generator',
    loadComponent: () => import('./features/qr-generator/qr-generator.component')
      .then(m => m.QRGeneratorComponent),
    canActivate: [AuthGuard],
    title: 'Generate QR Code - WiFi QR Code Generator'
  },
  {
    path: 'my-codes',
    loadComponent: () => import('./features/my-codes/my-codes.component')
      .then(m => m.MyCodesComponent),
    canActivate: [AuthGuard],
    title: 'My QR Codes - WiFi QR Code Generator'
  },

  // Admin routes (requires admin role)
  {
    path: 'admin',
    canActivate: [AuthGuard, AdminGuard],
    children: [
      {
        path: 'credentials',
        loadComponent: () => import('./features/admin/credentials/credentials.component')
          .then(m => m.CredentialsComponent),
        title: 'All Credentials - Admin'
      },
      {
        path: '',
        redirectTo: 'credentials',
        pathMatch: 'full'
      }
    ]
  },

  // Redirects
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: '**', redirectTo: '/dashboard' }
];
```

### Navigation Structure

```
Public Pages (No Authentication)
├── /login
└── /register

Protected Pages (Authentication Required)
├── /dashboard
├── /qr-generator
└── /my-codes

Admin Pages (Admin Role Required)
└── /admin/credentials
```

---

## Component Specifications

### 1. App Component (Root)

**Path**: `app/app.component.ts`
**Type**: Smart Component
**Purpose**: Application shell with router outlet and navigation

```typescript
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { NavbarComponent } from './shared/components/navbar/navbar.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, NavbarComponent],
  template: `
    <div class="min-h-screen bg-gray-50">
      <app-navbar />
      <main class="container mx-auto px-4 py-8">
        <router-outlet />
      </main>
    </div>
  `
})
export class AppComponent {
  title = 'WiFi QR Code Generator';
}
```

---

### 2. Navbar Component

**Path**: `shared/components/navbar/navbar.component.ts`
**Type**: Smart Component
**Purpose**: Application navigation with authentication state

#### Features
- Show/hide navigation items based on authentication
- Display user email when authenticated
- Logout functionality
- Admin-only menu items

#### Template Structure
```html
<nav class="bg-white shadow-lg">
  <div class="container mx-auto px-4">
    <div class="flex justify-between items-center py-4">
      <!-- Logo -->
      <div class="text-2xl font-bold text-blue-600">
        <a routerLink="/">WiFi QR Generator</a>
      </div>

      <!-- Navigation Links (Authenticated) -->
      @if (authStore.isAuthenticated()) {
        <div class="flex items-center space-x-6">
          <a routerLink="/dashboard" routerLinkActive="text-blue-600" class="hover:text-blue-600">
            Dashboard
          </a>
          <a routerLink="/qr-generator" routerLinkActive="text-blue-600" class="hover:text-blue-600">
            Generate QR
          </a>
          <a routerLink="/my-codes" routerLinkActive="text-blue-600" class="hover:text-blue-600">
            My Codes
          </a>

          @if (authStore.isAdmin()) {
            <a routerLink="/admin/credentials" routerLinkActive="text-blue-600" class="hover:text-blue-600">
              Admin
            </a>
          }

          <!-- User Menu -->
          <div class="relative">
            <button (click)="toggleUserMenu()" class="flex items-center space-x-2">
              <span>{{ authStore.user()?.email }}</span>
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"/>
              </svg>
            </button>

            @if (showUserMenu()) {
              <div class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg">
                <button (click)="logout()" class="block w-full text-left px-4 py-2 hover:bg-gray-100">
                  Logout
                </button>
              </div>
            }
          </div>
        </div>
      } @else {
        <!-- Public Links -->
        <div class="flex items-center space-x-4">
          <a routerLink="/login" class="text-gray-700 hover:text-blue-600">Login</a>
          <a routerLink="/register" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
            Sign Up
          </a>
        </div>
      }
    </div>
  </div>
</nav>
```

---

### 3. Login Component

**Path**: `features/auth/login/login.component.ts`
**Type**: Smart Component
**Purpose**: User authentication

#### Component Implementation

```typescript
import { Component, inject, signal } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../../core/services/auth.service';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterLink],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {
  private readonly authService = inject(AuthService);
  private readonly fb = inject(FormBuilder);
  private readonly router = inject(Router);

  loginForm: FormGroup;
  loading = signal(false);
  error = signal<string | null>(null);

  constructor() {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', Validators.required]
    });
  }

  onSubmit(): void {
    if (this.loginForm.invalid) {
      return;
    }

    this.loading.set(true);
    this.error.set(null);

    this.authService.login(this.loginForm.value).subscribe({
      next: () => {
        this.loading.set(false);
        // Navigation handled by AuthService
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(err.error?.error?.message || 'Login failed. Please try again.');
      }
    });
  }

  get email() {
    return this.loginForm.get('email');
  }

  get password() {
    return this.loginForm.get('password');
  }
}
```

#### Template

```html
<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        Sign in to your account
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Or
        <a routerLink="/register" class="font-medium text-blue-600 hover:text-blue-500">
          create a new account
        </a>
      </p>
    </div>

    <form [formGroup]="loginForm" (ngSubmit)="onSubmit()" class="mt-8 space-y-6">
      @if (error()) {
        <div class="rounded-md bg-red-50 p-4">
          <p class="text-sm text-red-800">{{ error() }}</p>
        </div>
      }

      <div class="rounded-md shadow-sm -space-y-px">
        <div>
          <label for="email" class="sr-only">Email address</label>
          <input
            id="email"
            formControlName="email"
            type="email"
            autocomplete="email"
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
            placeholder="Email address"
          />
          @if (email?.invalid && email?.touched) {
            <p class="mt-1 text-sm text-red-600">
              @if (email?.errors?.['required']) {
                Email is required
              } @else if (email?.errors?.['email']) {
                Invalid email format
              }
            </p>
          }
        </div>

        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            formControlName="password"
            type="password"
            autocomplete="current-password"
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-blue-500 focus:border-blue-500 focus:z-10 sm:text-sm"
            placeholder="Password"
          />
          @if (password?.invalid && password?.touched) {
            <p class="mt-1 text-sm text-red-600">Password is required</p>
          }
        </div>
      </div>

      <div>
        <button
          type="submit"
          [disabled]="loginForm.invalid || loading()"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          @if (loading()) {
            Signing in...
          } @else {
            Sign in
          }
        </button>
      </div>
    </form>
  </div>
</div>
```

---

### 4. QR Generator Component

**Path**: `features/qr-generator/qr-generator.component.ts`
**Type**: Smart Component
**Purpose**: WiFi QR code creation interface

#### Component Implementation

```typescript
import { Component, inject, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { QRCodeService } from '../../core/services/qr-code.service';
import { QRCodeStore } from '../../core/stores/qr-code.store';
import { WiFiQRCode, CreateQRCodeRequest } from '../../core/models/qr-code.model';
import { QRFormComponent } from './components/qr-form/qr-form.component';
import { QRDisplayComponent } from './components/qr-display/qr-display.component';

@Component({
  selector: 'app-qr-generator',
  standalone: true,
  imports: [CommonModule, QRFormComponent, QRDisplayComponent],
  template: `
    <div class="container mx-auto px-4 py-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-8">Generate WiFi QR Code</h1>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div>
          <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-xl font-semibold mb-4">WiFi Network Details</h2>
            <app-qr-form
              [loading]="qrStore.loading()"
              (submitForm)="onCreateQRCode($event)"
            />
          </div>
        </div>

        <div>
          @if (generatedQRCode()) {
            <div class="bg-white rounded-lg shadow-md p-6">
              <h2 class="text-xl font-semibold mb-4">Generated QR Code</h2>
              <app-qr-display [qrCode]="generatedQRCode()!" />
            </div>
          } @else {
            <div class="bg-gray-50 rounded-lg border-2 border-dashed border-gray-300 p-12 text-center">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
              </svg>
              <p class="mt-2 text-sm text-gray-600">QR code will appear here</p>
            </div>
          }
        </div>
      </div>
    </div>
  `
})
export class QRGeneratorComponent {
  private readonly qrService = inject(QRCodeService);
  readonly qrStore = inject(QRCodeStore);

  readonly generatedQRCode = signal<WiFiQRCode | null>(null);

  onCreateQRCode(request: CreateQRCodeRequest): void {
    this.qrService.createQRCode(request).subscribe({
      next: (qrCode) => {
        this.generatedQRCode.set(qrCode);
      },
      error: (error) => {
        console.error('Error creating QR code:', error);
      }
    });
  }
}
```

---

### 5. QR Form Component (Dumb)

**Path**: `features/qr-generator/components/qr-form/qr-form.component.ts`
**Type**: Dumb/Presentational Component
**Purpose**: WiFi credentials input form

#### Component Implementation

```typescript
import { Component, output, input } from '@angular/core';
import { FormBuilder, FormGroup, Validators, ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { CreateQRCodeRequest } from '../../../../core/models/qr-code.model';

@Component({
  selector: 'app-qr-form',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './qr-form.component.html',
  styleUrls: ['./qr-form.component.css']
})
export class QRFormComponent {
  loading = input<boolean>(false);
  submitForm = output<CreateQRCodeRequest>();

  form: FormGroup;

  securityTypes = [
    { value: 'WPA2', label: 'WPA2 (Recommended)' },
    { value: 'WPA', label: 'WPA' },
    { value: 'WEP', label: 'WEP (Legacy)' },
    { value: 'nopass', label: 'No Password (Open)' }
  ];

  constructor(private fb: FormBuilder) {
    this.form = this.fb.group({
      ssid: ['', [Validators.required, Validators.maxLength(32)]],
      password: [''],
      securityType: ['WPA2', Validators.required],
      isHidden: [false]
    }, { validators: this.passwordRequiredValidator });
  }

  passwordRequiredValidator(form: FormGroup) {
    const securityType = form.get('securityType')?.value;
    const password = form.get('password')?.value;

    if (securityType !== 'nopass' && !password) {
      form.get('password')?.setErrors({ required: true });
      return { passwordRequired: true };
    }

    if (password && password.length > 63) {
      form.get('password')?.setErrors({ maxlength: true });
      return { passwordMaxLength: true };
    }

    form.get('password')?.setErrors(null);
    return null;
  }

  onSubmit(): void {
    if (this.form.valid) {
      this.submitForm.emit(this.form.value);
      this.form.reset({ securityType: 'WPA2', isHidden: false });
    }
  }

  get ssid() {
    return this.form.get('ssid');
  }

  get password() {
    return this.form.get('password');
  }
}
```

---

### 6. QR Display Component (Dumb)

**Path**: `features/qr-generator/components/qr-display/qr-display.component.ts`
**Type**: Dumb Component
**Purpose**: Display generated QR code

#### Component Implementation

```typescript
import { Component, input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { WiFiQRCode } from '../../../../core/models/qr-code.model';
import { QRCodeModule } from 'angularx-qrcode';

@Component({
  selector: 'app-qr-display',
  standalone: true,
  imports: [CommonModule, QRCodeModule],
  template: `
    <div class="flex flex-col items-center space-y-4">
      <div class="bg-white p-4 rounded-lg border-2 border-gray-200">
        <qrcode
          [qrdata]="qrCode().qrCodeData"
          [width]="256"
          [errorCorrectionLevel]="'M'"
        ></qrcode>
      </div>

      <div class="w-full space-y-2">
        <div class="flex justify-between text-sm">
          <span class="font-medium text-gray-700">Network Name:</span>
          <span class="text-gray-900">{{ qrCode().ssid }}</span>
        </div>

        @if (qrCode().securityType !== 'nopass') {
          <div class="flex justify-between text-sm">
            <span class="font-medium text-gray-700">Password:</span>
            <span class="text-gray-900 font-mono">{{ qrCode().password }}</span>
          </div>
        }

        <div class="flex justify-between text-sm">
          <span class="font-medium text-gray-700">Security:</span>
          <span class="text-gray-900">{{ qrCode().securityType }}</span>
        </div>

        <div class="flex justify-between text-sm">
          <span class="font-medium text-gray-700">Hidden Network:</span>
          <span class="text-gray-900">{{ qrCode().isHidden ? 'Yes' : 'No' }}</span>
        </div>
      </div>

      <div class="w-full pt-4 border-t border-gray-200">
        <p class="text-xs text-gray-500 text-center">
          Scan this QR code with your mobile device to connect to the WiFi network
        </p>
      </div>
    </div>
  `
})
export class QRDisplayComponent {
  qrCode = input.required<WiFiQRCode>();
}
```

---

### 7. My Codes Component

**Path**: `features/my-codes/my-codes.component.ts`
**Type**: Smart Component
**Purpose**: Display user's QR codes with pagination

#### Features
- List all user's QR codes
- Pagination support
- Delete functionality
- View QR code details

---

### 8. Admin Credentials Component

**Path**: `features/admin/credentials/credentials.component.ts`
**Type**: Smart Component
**Purpose**: Admin view of all credentials

#### Features
- View all WiFi credentials across users
- Search by SSID
- Filter by user
- Pagination
- Display user email alongside credentials

---

## State Management

### Auth Store

**Path**: `core/stores/auth.store.ts`

```typescript
import { computed, Injectable, signal } from '@angular/core';
import { User } from '../models/user.model';

@Injectable({ providedIn: 'root' })
export class AuthStore {
  private readonly _user = signal<User | null>(null);
  private readonly _token = signal<string | null>(null);

  readonly user = this._user.asReadonly();
  readonly token = this._token.asReadonly();

  readonly isAuthenticated = computed(() => !!this._token());
  readonly isAdmin = computed(() => this._user()?.role === 'admin');

  setAuth(user: User, token: string): void {
    this._user.set(user);
    this._token.set(token);
  }

  clearAuth(): void {
    this._user.set(null);
    this._token.set(null);
  }

  updateUser(user: User): void {
    this._user.set(user);
  }
}
```

### QR Code Store

**Path**: `core/stores/qr-code.store.ts`

```typescript
import { computed, Injectable, signal } from '@angular/core';
import { WiFiQRCode } from '../models/qr-code.model';

@Injectable({ providedIn: 'root' })
export class QRCodeStore {
  private readonly _qrCodes = signal<WiFiQRCode[]>([]);
  private readonly _loading = signal<boolean>(false);
  private readonly _error = signal<string | null>(null);

  readonly qrCodes = this._qrCodes.asReadonly();
  readonly loading = this._loading.asReadonly();
  readonly error = this._error.asReadonly();
  readonly qrCodeCount = computed(() => this._qrCodes().length);

  setQRCodes(codes: WiFiQRCode[]): void {
    this._qrCodes.set(codes);
  }

  addQRCode(code: WiFiQRCode): void {
    this._qrCodes.update(codes => [code, ...codes]);
  }

  removeQRCode(id: string): void {
    this._qrCodes.update(codes => codes.filter(c => c.id !== id));
  }

  setLoading(loading: boolean): void {
    this._loading.set(loading);
  }

  setError(error: string | null): void {
    this._error.set(error);
  }

  clear(): void {
    this._qrCodes.set([]);
    this._error.set(null);
  }
}
```

---

## Services

### Auth Service

**Path**: `core/services/auth.service.ts`

```typescript
import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthStore } from '../stores/auth.store';
import { StorageService } from './storage.service';

interface LoginRequest {
  email: string;
  password: string;
}

interface AuthResponse {
  success: boolean;
  data: {
    token: string;
    user: {
      id: string;
      email: string;
      role: string;
      createdAt: string;
    };
  };
  message: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly authStore = inject(AuthStore);
  private readonly storage = inject(StorageService);

  private readonly baseUrl = `${environment.apiUrl}/auth`;

  login(credentials: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/login`, credentials).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  register(data: LoginRequest): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(`${this.baseUrl}/register`, data).pipe(
      tap(response => this.handleAuthSuccess(response))
    );
  }

  logout(): void {
    this.authStore.clearAuth();
    this.storage.remove('token');
    this.router.navigate(['/login']);
  }

  initializeAuth(): void {
    const token = this.storage.get('token');
    if (token) {
      try {
        const payload = this.decodeToken(token);
        this.authStore.setAuth(
          {
            id: payload.userId,
            email: payload.email,
            role: payload.role,
            createdAt: ''
          },
          token
        );
      } catch {
        this.storage.remove('token');
      }
    }
  }

  private handleAuthSuccess(response: AuthResponse): void {
    this.authStore.setAuth(response.data.user, response.data.token);
    this.storage.set('token', response.data.token);
    this.router.navigate(['/dashboard']);
  }

  private decodeToken(token: string): any {
    const payload = token.split('.')[1];
    return JSON.parse(atob(payload));
  }
}
```

---

## Guards & Interceptors

### Auth Guard

**Path**: `core/guards/auth.guard.ts`

```typescript
import { inject } from '@angular/core';
import { Router, CanActivateFn } from '@angular/router';
import { AuthStore } from '../stores/auth.store';

export const AuthGuard: CanActivateFn = () => {
  const authStore = inject(AuthStore);
  const router = inject(Router);

  if (authStore.isAuthenticated()) {
    return true;
  }

  router.navigate(['/login']);
  return false;
};
```

### JWT Interceptor

**Path**: `core/interceptors/jwt.interceptor.ts`

```typescript
import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthStore } from '../stores/auth.store';

export const jwtInterceptor: HttpInterceptorFn = (req, next) => {
  const authStore = inject(AuthStore);
  const token = authStore.token();

  if (token) {
    req = req.clone({
      setHeaders: {
        Authorization: `Bearer ${token}`
      }
    });
  }

  return next(req);
};
```

---

## Forms & Validation

### Validation Rules

#### Email
- Required
- Valid email format
- Max 255 characters

#### Password (Registration)
- Required
- Min 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 number

#### SSID
- Required
- Max 32 characters (IEEE 802.11 spec)

#### WiFi Password
- Conditional: Required if security type is not 'nopass'
- Max 63 characters

---

## UI/UX Guidelines

### Color Scheme (Tailwind)

- **Primary**: Blue (blue-600, blue-700)
- **Success**: Green (green-500, green-600)
- **Error**: Red (red-500, red-600)
- **Warning**: Yellow (yellow-500, yellow-600)
- **Neutral**: Gray (gray-50 to gray-900)

### Typography
- **Headings**: Font weight 700 (font-bold)
- **Body**: Font weight 400 (normal)
- **H1**: text-3xl
- **H2**: text-2xl
- **H3**: text-xl
- **Body**: text-base

### Spacing
- Container: `container mx-auto px-4`
- Section padding: `py-8`
- Card padding: `p-6`
- Element spacing: `space-y-4`, `space-x-4`

### Components
- **Cards**: White background, rounded-lg, shadow-md
- **Buttons**: Rounded, padding py-2 px-4, hover states
- **Forms**: Border, rounded, focus ring
- **Input fields**: Border-gray-300, focus:border-blue-500

---

## Performance Optimization

### Lazy Loading
- All feature routes lazy-loaded
- Components loaded on-demand

### Change Detection
- Use OnPush where appropriate
- Signals for reactive state

### Bundle Optimization
- Tree-shaking via standalone components
- Code splitting per route
- Production build with optimization

### Best Practices
- Avoid unnecessary re-renders
- Use computed signals for derived state
- Unsubscribe from observables (use async pipe or takeUntilDestroyed)
- Optimize images and assets

---

**Document Version**: 1.0
**Last Updated**: 2025-01-15
**Maintainer**: Frontend Team
