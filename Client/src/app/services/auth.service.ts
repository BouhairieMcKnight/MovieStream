import {inject, Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment.development';

@Injectable({
  providedIn: 'root'
})

export class AuthService {
  http: HttpClient = inject(HttpClient);

  router: Router = inject(Router);

  signUp(): void {
    this.http.post(environment.apiUrl + '/users/register', {
    })
  }



  signIn(): void {

  }

  signOut(): void {

  }

  constructor() {

  }
}
