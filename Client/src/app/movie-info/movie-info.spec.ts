import { ComponentFixture, TestBed } from '@angular/core/testing';

import { MovieInfo } from './movie-info';

describe('MovieInfo', () => {
  let component: MovieInfo;
  let fixture: ComponentFixture<MovieInfo>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MovieInfo],
    }).compileComponents();

    fixture = TestBed.createComponent(MovieInfo);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
