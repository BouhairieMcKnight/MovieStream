import { Routes } from '@angular/router';
import {Browser} from './pages/browser/browser';

export const routes: Routes = [
  // { path: '', loadComponent: () => import('./pages/login/login').then(a => a.Login)},
  { path: '', component: Browser }
];
