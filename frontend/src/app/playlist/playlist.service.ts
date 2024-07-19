import { baseUrl } from "app/baseUrl";
import { catchError, map, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { addPropertyIsSelected } from "@utils/filter";

import { PlaylistResponseWithTracks, PlaylistWithTracks } from "./playlist";

@Injectable({
  providedIn: "root",
})
export class PlaylistService {
  private url = "platforms/:platform/playlist/:id";
  private http = inject(HttpClient);

  private playlistQuery = new APIQuery<PlaylistWithTracks | null, Error>(null);
  playlist = new DataQuery(this.playlistQuery);

  fetchPlaylist(platform: string, id: string) {
    this.playlistQuery.fetching();

    const urlRequest = this.url
      .replace(":platform", platform)
      .replace(":id", id);

    this.http
      .get<PlaylistResponseWithTracks>(baseUrl + urlRequest)
      .pipe(
        map((playlist) => {
          playlist.tracks = addPropertyIsSelected(playlist.tracks);
          return playlist as PlaylistWithTracks;
        }),
        tap((playlist) => {
          this.playlistQuery.success(playlist);
        }),
        catchError((error) => {
          this.playlistQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  destroy() {
    this.playlistQuery.restoreInitialState();
  }
}
