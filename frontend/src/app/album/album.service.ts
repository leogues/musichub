import { baseUrl } from "app/baseUrl";
import { catchError, map, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { addPropertyIsSelected } from "@utils/filter";

import { AlbumResponseWithTracks, AlbumWithTracks } from "./album";

@Injectable({
  providedIn: "root",
})
export class AlbumService {
  private url = "platforms/:platform/album/:albumId";
  private http = inject(HttpClient);

  private albumQuery = new APIQuery<AlbumWithTracks | null, Error>(null);
  album = new DataQuery(this.albumQuery);

  fetchAlbum(platform: string, albumId: string) {
    this.albumQuery.fetching();

    const urlRequest = this.url
      .replace(":platform", platform)
      .replace(":albumId", albumId);
    this.http
      .get<AlbumResponseWithTracks>(baseUrl + urlRequest)
      .pipe(
        map((album) => {
          album.tracks = addPropertyIsSelected(album.tracks);
          return album as AlbumWithTracks;
        }),
        tap((album) => {
          this.albumQuery.success(album);
        }),
        catchError((error) => {
          this.albumQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  destroy() {
    this.albumQuery.restoreInitialState();
    console.log("AlbumService destroyed");
  }
}
