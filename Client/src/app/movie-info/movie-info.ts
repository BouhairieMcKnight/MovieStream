import {Component, effect, EventEmitter, inject, Input, OnDestroy, OnInit, Output, ViewChild} from '@angular/core';
import {Movie} from '../services/models/movie_model';
import {MovieStreamService} from '../services/movie.stream.service';
import {NgOptimizedImage} from '@angular/common';
import {CdkPortal} from '@angular/cdk/portal';
import {Overlay, OverlayConfig} from '@angular/cdk/overlay';

@Component({
  selector: 'app-movie-info',
  standalone: true,
  imports: [
    NgOptimizedImage,
    CdkPortal
  ],
  templateUrl: './movie-info.html',
  styleUrl: './movie-info.scss',
})
export class MovieInfo implements OnInit, OnDestroy {
  @ViewChild(CdkPortal) portal: CdkPortal | undefined;
  @Output() closeModal = new EventEmitter<void>();

  overlay = inject(Overlay)
  overlayConfig = new OverlayConfig({
    hasBackdrop: true,
    positionStrategy: this.overlay
      .position()
      .global()
      .centerVertically()
      .centerHorizontally(),
    scrollStrategy: this.overlay.scrollStrategies.block(),
    minWidth: 300,
  });
  overlayRef = this.overlay.create(this.overlayConfig);




  ngAfterViewInit() {
    this.overlayRef.attach(this.portal);
  }

  ngOnInit(): void {
    const overlayConfig = new OverlayConfig({
      hasBackdrop: true,
      positionStrategy: this.overlay
        .position()
        .global()
        .centerVertically()
        .centerHorizontally(),
      scrollStrategy: this.overlay.scrollStrategies.close(),
      minWidth: 500,
    });

    this.overlayRef = this.overlay.create(overlayConfig);

    this.overlayRef.backdropClick().subscribe(() => {
      this.closeModal.emit();
    });
  }

  ngOnDestroy(): void {
    this.overlayRef.detach();
    this.overlayRef.dispose();
  }
}
