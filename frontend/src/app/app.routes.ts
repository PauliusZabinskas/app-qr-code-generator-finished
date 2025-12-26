import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { adminGuard } from './core/guards/admin.guard';
import { LoginComponent } from './features/auth/login/login.component';
import { RegisterComponent } from './features/auth/register/register.component';
import { DashboardComponent } from './features/dashboard/dashboard.component';
import { QrGeneratorComponent } from './features/qr-generator/qr-generator.component';
import { MyCodesComponent } from './features/my-codes/my-codes.component';
import { CredentialsComponent } from './features/admin/credentials/credentials.component';

export const routes: Routes = [
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: 'login', component: LoginComponent },
  { path: 'register', component: RegisterComponent },
  { path: 'dashboard', component: DashboardComponent, canActivate: [authGuard] },
  { path: 'qr-generator', component: QrGeneratorComponent, canActivate: [authGuard] },
  { path: 'my-codes', component: MyCodesComponent, canActivate: [authGuard] },
  { path: 'admin/credentials', component: CredentialsComponent, canActivate: [authGuard, adminGuard] },
  { path: '**', redirectTo: '/dashboard' }
];
