import {Component, effect, inject, OnInit} from '@angular/core';
import {Movie} from '../services/models/movie_model';
import {MovieStreamService} from '../services/movie.stream.service';
import { ActivatedRoute } from "@angular/router";
import {debounce, filter, interval, map} from 'rxjs';
import {MovieCard} from '../shared/movie-card/movie-card';
import {InfiniteScrollDirective} from '../shared/infinite.scroll.directive';

@Component({
  selector: 'app-search',
  imports: [
    MovieCard,
    InfiniteScrollDirective
  ],
  templateUrl: './search.html',
  styleUrl: './search.scss',
})
export class Search implements OnInit {
  term: string = "";
  page: number = 0;
  genre: string = "";
  sort: "asc" | "desc" = "desc";
  hasMore: boolean = true;

  movieStreamService: MovieStreamService = inject(MovieStreamService);
  activatedRoute = inject(ActivatedRoute);

  movies: Movie[] = [];
  ngOnInit(): void {
    this.onSearch();
  }

  onNearEndScroll(): void {
    this.page++;
    this.onSearch();
  }

  private onSearch(): void {
    this.activatedRoute.queryParams.pipe(
      filter(params => params['q']),
      debounce(() => interval(300)),
      map(query => ({
        term: String(query['q']),
        page: this.page,
        genre: this.genre,
        sort: this.sort
      })),
    ).subscribe({
      next: query =>
        this.movieStreamService.getMoviesSearch(query.term, query.page, query.genre, query.sort)
    })
  }
  constructor() {
    effect(() => {
      let response = this.movieStreamService.search().value;
      if(response) {
        this.movies = [...this.movies, ...response];
        this.hasMore = response.length > 6;
      }
    });
  }
}
