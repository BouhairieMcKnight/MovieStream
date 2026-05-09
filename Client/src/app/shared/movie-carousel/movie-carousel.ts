import {
  AfterViewInit,
  Component,
  DestroyRef,
  effect,
  ElementRef, EventEmitter,
  inject,
  Input,
  OnInit, Output, signal,
  ViewChild, WritableSignal
} from '@angular/core';
import Swiper from 'swiper';
import {Movie} from '../../services/models/movie_model';
import {MovieStreamService} from '../../services/movie.stream.service';
import {HttpErrorResponse} from '@angular/common/http';
import {State} from '../../services/models/state.model';
import {takeUntilDestroyed} from '@angular/core/rxjs-interop';
import {NgbModalModule} from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-movie-carousel',
  standalone: true,
  imports: [],
  templateUrl: './movie-carousel.html',
  styleUrl: './movie-carousel.scss',
})

export class MovieCarousel implements OnInit, AfterViewInit {
  @Input() title!: string;
  @ViewChild('swiperContainer') swiperContainer!: ElementRef;
  @Output() openModal = new EventEmitter<Movie>();

  movieStreamService: MovieStreamService = inject(MovieStreamService);
  selectedContent: string | null = null;

  movie$: WritableSignal<State<Movie[], HttpErrorResponse>>
    = signal(State.Builder<Movie[], HttpErrorResponse>().forInit().build());

  private readonly DestroyRef = inject(DestroyRef);

  onClick(movie: Movie) {
    this.openModal.emit(movie);
  }

  constructor() {
  }

  ngAfterViewInit(): void {
    this.initSwiper();
  }

  ngOnInit() {
    this.movieStreamService.getMoviesByGenre(this.title)
      .pipe(takeUntilDestroyed(this.DestroyRef))
      .subscribe({
        next: response =>
          this.movie$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forSuccess(response).build()),
        error: err => {
          this.movie$.set(
            State.Builder<Movie[], HttpErrorResponse>()
              .forError(err).build())
        }
      });
  }

  private initSwiper() {
    return new Swiper(this.swiperContainer.nativeElement, {
      slidesPerView: 3,
      slidesPerGroup: 3,
      centeredSlides: true,
      loop: false,
      slidesPerGroupAuto: true,
      edgeSwipeThreshold: 0,
      spaceBetween: 30,
      pagination: {
        el: '.swiper-pagination',
        clickable: true,
      },
      watchOverflow: true,
      breakpoints: {
        600: {
          slidesPerView: 2,
          slidesPerGroup: 2,
          spaceBetween: 5,
          centeredSlides: true,
        },
        900: {
          slidesPerView: 3,
          slidesPerGroup: 3,
          spaceBetween: 5,
          centeredSlides: true,
        },
        1200: {
          slidesPerView: 4,
          slidesPerGroup: 4,
          spaceBetween: 5,
          centeredSlides: false,
        },
        1500: {
          slidesPerView: 5,
          slidesPerGroup: 5,
          spaceBetween: 5,
          centeredSlides: false,
        },
        1800: {
          slidesPerView: 5,
          slidesPerGroup: 6,
          spaceBetween: 5,
          centeredSlides: false,
        }
      }
    })
  }

  setHoverMovie(movie: Movie) {
    this.selectedContent = movie.title;
  }

  clearHoverMovie() {
    this.selectedContent = null;
  }
}
