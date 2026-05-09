import {Component, inject, Input} from '@angular/core';
import {Movie} from '../../services/models/movie_model';
import {MovieStreamService} from '../../services/movie.stream.service';

@Component({
  standalone: true,
  selector: 'app-movie-card',
  imports: [],
  templateUrl: './movie-card.html',
  styleUrl: './movie-card.scss',
})
export class MovieCard {
  @Input() movie: Movie | undefined;

  movieStreamService: MovieStreamService = inject(MovieStreamService);
}
