import { Component, computed, inject, input } from '@angular/core';
import { TrackContent } from '@components/music-table/music-table';
import { MusicPlayerService } from '@services/music-player.service';

@Component({
  selector: "app-table-track",
  standalone: true,
  imports: [],
  templateUrl: "./table-track.component.html",
})
export class TableTrackComponent {
  track = input.required<TrackContent>();
  trackId = computed(() => this.track().id);

  private musicPlayerService = inject(MusicPlayerService);
  private playingTrackSignal = this.musicPlayerService.getPlayingTrack();

  isPlaying = computed(() => {
    const playingTrack = this.playingTrackSignal();
    return playingTrack === this.trackId();
  });

  togglePreview() {
    if (!this.track().trackurl) {
      return;
    }

    if (this.isPlaying()) {
      this.musicPlayerService.pause();
    } else {
      this.musicPlayerService.play(this.trackId());
    }
  }
}
