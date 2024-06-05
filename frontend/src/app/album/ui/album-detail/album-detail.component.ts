import { AlbumWithTracks } from "app/album/album";

import { Component, input } from "@angular/core";
import { DetailComponent } from "@components/detail/detail.component";

@Component({
  selector: "app-album-detail",
  standalone: true,
  imports: [DetailComponent],
  templateUrl: "./album-detail.component.html",
})
export class AlbumDetailComponent {
  album = input.required<AlbumWithTracks>();
  albumTracksCount = input.required<number>();
}
