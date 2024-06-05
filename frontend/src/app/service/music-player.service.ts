import { Injectable, signal } from "@angular/core";

@Injectable({
  providedIn: "root",
})
export class MusicPlayerService {
  private playingTrack = signal<string | null>(null);

  play(trackId: string) {
    this.playingTrack.set(trackId);
  }

  ended() {
    this.playingTrack.set(null);
  }

  pause() {
    this.playingTrack.set(null);
  }

  getPlayingTrack() {
    return this.playingTrack.asReadonly();
  }
}
