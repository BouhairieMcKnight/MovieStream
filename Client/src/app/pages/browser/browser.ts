import {Component, effect, inject, OnInit, signal, Signal, ViewChild} from '@angular/core';
import {AuthService} from '../../services/auth.service';
import {MovieStreamService} from '../../services/movie.stream.service';
import {Genre, GenresResponse, Movie} from '../../services/models/movie_model';
import {MovieCarousel} from '../../shared/movie-carousel/movie-carousel';
import {FaIconComponent, FaIconLibrary} from '@fortawesome/angular-fontawesome';
import {MovieCard} from '../../shared/movie-card/movie-card';
import {NgbModal, NgbModalModule} from '@ng-bootstrap/ng-bootstrap';
import {MovieInfo} from '../../movie-info/movie-info';
import {Navbar} from '../../navbar/navbar';

@Component({
  selector: 'app-browser',
  standalone: true,
  imports: [
    MovieCard,
    MovieCarousel,
    FaIconComponent,
    MovieInfo,
    Navbar
  ],
  templateUrl: './browser.html',
  styleUrl: './browser.scss',
})
export class Browser implements OnInit {
  movieStreamService: MovieStreamService = inject(MovieStreamService);
  isModalOpen = signal(false);
  movie: Movie | undefined;

  openModal(movie: Movie): void {
    this.isModalOpen.set(true);
    this.movie = movie;
  }

  closeModal(): void {
    this.isModalOpen.set(false);
  }

  constructor() {
  }

  ngOnInit(): void {
    this.fetchMovieGenres();
  }

  fetchMovieGenres(): void {
    this.movieStreamService.getAllGenres();
  }
}
