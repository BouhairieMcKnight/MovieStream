import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MovieCarousel } from './movie-carousel';

describe('MovieCarousel', () => {
  let component: MovieCarousel;
  let fixture: ComponentFixture<MovieCarousel>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MovieCarousel],
    }).compileComponents();

    fixture = TestBed.createComponent(MovieCarousel);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
