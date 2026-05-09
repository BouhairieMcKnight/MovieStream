import {computed, inject, Injectable, OnInit, signal, WritableSignal} from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpParams} from '@angular/common/http';
import {Genre, Movie} from './models/movie_model';
import {environment} from '../../environments/environment.development';
import {State} from './models/state.model';
import {NgbModal} from '@ng-bootstrap/ng-bootstrap';
import {DomSanitizer, SafeResourceUrl} from '@angular/platform-browser';
import {map, startWith} from 'rxjs';
import {toSignal} from '@angular/core/rxjs-interop';

@Injectable({
  providedIn: 'root'
})

export class MovieStreamService {
  http: HttpClient = inject(HttpClient);
  private sanitizer = inject(DomSanitizer)
  modalService: NgbModal = inject(NgbModal)

  private fetchTrendMovie$: WritableSignal<State<Movie[], HttpErrorResponse>>
    = signal(State.Builder<Movie[], HttpErrorResponse>().forInit().build());
  fetchTrendMovie = computed(() => this.fetchTrendMovie$());

  private genres$: WritableSignal<State<Genre[], HttpErrorResponse>>
    = signal(State.Builder<Genre[], HttpErrorResponse>().forInit().build());
  genres = computed(() => this.genres$());

  private movieById$: WritableSignal<State<Movie, HttpErrorResponse>>
    = signal(State.Builder<Movie, HttpErrorResponse>().forInit().build());
  movieById = computed(() => this.movieById$());

  private search$: WritableSignal<State<Movie[], HttpErrorResponse>>
    = signal(State.Builder<Movie[], HttpErrorResponse>().forInit().build());
  search = computed(() => this.search$());

  getImageURL(url: string, size: 'original' | 'w500' | 'w200'): SafeResourceUrl {
    const id = url.split("/").pop()?.replace(".jpg", "");
    return this.sanitizer.bypassSecurityTrustResourceUrl(`https://image.tmdb.org/t/p/${size}/${id}.jpg`);
  }

  getAllGenres(): void {
    this.http.get<Genre[]>(environment.apiUrl + '/movies/genres', {headers: {'Content-Type': 'application/json'}})
      .subscribe({
        next: response =>
          this.genres$.set(
            State.Builder<Genre[], HttpErrorResponse>()
              .forSuccess(response).build()),
        error: err => {
          this.genres$.set(
            State.Builder<Genre[], HttpErrorResponse>()
              .forError(err).build())
        }
      });
  }

  getMoviesByGenre(genre: string) {
    let queryParams = new HttpParams();
    queryParams = queryParams.append('genre', genre);
    return this.http.get<Movie[]>(environment.apiUrl + '/movies/search', {params: queryParams})
  }

  getMoviesSearch(term: string, page: number, genre: string, sort: "asc" | "desc" = "desc"): void {
    let queryParams = new HttpParams();
    queryParams = queryParams.append('term', term);
    queryParams = queryParams.append('page', page);
    queryParams = queryParams.append('sort', sort);
    queryParams = queryParams.append('genre', genre);
    this.http.get<Movie[]>(environment.apiUrl + '/movies/search', {params: queryParams})
      .subscribe({
        next: response =>
          this.search$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forSuccess(response).build()),
        error: err => {
          this.search$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forError(err).build())
        }
      });
  }

  getMovieById(movieId: string): void {
    this.http.get<Movie>(environment.apiUrl + `/movie/${movieId}`)
      .subscribe({
        next: response =>
          this.movieById$.set(
            State.Builder<Movie, HttpErrorResponse>()
              .forSuccess(response).build()),
        error: err => {
          this.movieById$.set(
            State.Builder<Movie, HttpErrorResponse>()
              .forError(err).build())
        }
      });
  }

  getTrendMovies(): void {
    this.http.get<Movie[]>(environment.apiUrl + '/movies/trending')
      .subscribe({
        next: response =>
          this.fetchTrendMovie$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forSuccess(response).build()),
        error: err => {
          this.fetchTrendMovie$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forError(err).build())
        }
      });
  }

  clearGetMovieById() {
    this.movieById$.set(State.Builder<Movie, HttpErrorResponse>().forInit().build());
  }

  constructor() {}
}
