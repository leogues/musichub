import { PlaylistWithTracks } from 'app/playlist/playlist';

import { Component, input } from '@angular/core';
import { DetailComponent } from '@components/detail/detail.component';

@Component({
  selector: "app-playlist-detail",
  standalone: true,
  imports: [DetailComponent],
  templateUrl: "./playlist-detail.component.html",
})
export class PlaylistDetailComponent {
  playlist = input.required<PlaylistWithTracks>();
  playlistTracksCount = input.required<number>();
}
