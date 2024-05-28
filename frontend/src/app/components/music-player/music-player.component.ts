import { Component, effect, ElementRef, inject, Input, ViewChild } from '@angular/core';

import { Track } from '../../../types/track';
import { MusicPlayerService } from '../../service/music-player.service';

type MusicTrack = Record<string, Track>;

@Component({
  selector: "app-music-player",
  standalone: true,
  imports: [],
  templateUrl: "./music-player.component.html",
})
export class MusicPlayerComponent {
  @ViewChild("audioPlayer") audioPlayerRef!: ElementRef<HTMLAudioElement>;
  @Input({
    transform: (tracks: Track[]) => {
      return tracks.reduce((acc, track) => {
        acc[track.id] = track;
        return acc;
      }, {} as MusicTrack);
    },
  })
  tracks: MusicTrack = {};

  private musicPlayerService = inject(MusicPlayerService);
  constructor() {
    const playingTrackSignal = this.musicPlayerService.getPlayingTrack();

    effect(() => {
      const playTrackId = playingTrackSignal();
      if (playTrackId) {
        this.playPreview(playTrackId);
      } else {
        this.audioPlayerRef.nativeElement.pause();
      }
    });
  }

  playPreview(trackId: string) {
    const track = this.tracks[trackId];
    if (!track || !track.preview) {
      return;
    }

    if (this.audioPlayerRef.nativeElement.src === track.preview) {
      this.audioPlayerRef.nativeElement.play();
      return;
    }

    const audioPlayer = this.audioPlayerRef.nativeElement;
    audioPlayer.src = track.preview;
    audioPlayer.play();
  }

  endedPreview() {
    this.musicPlayerService.ended();
  }

  ngOnDestroy() {
    this.musicPlayerService.ended();
  }
}
