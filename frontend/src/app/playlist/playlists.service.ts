import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable, signal } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";

import { SupportedSources } from "../../types/providerAuth";
import { Playlist, PlaylistsResponse, ProvidersPlaylists } from "./playlist";
import { filterAddedItems, filterRemovedItems } from "@utils/filter";

@Injectable({
  providedIn: "root",
})
export class PlaylistsService {
  private providerAuthService = inject(ProviderAuthService);
  private http = inject(HttpClient);

  private url = "me/platforms/:source/playlist";

  private authenticatedSources =
    this.providerAuthService.authenticatedProvidersSources;
  private authenticatedSourcesObservable =
    this.providerAuthService.authenticatedSourcesObservable;

  private mePlaylistsQuery = new APIQuery<ProvidersPlaylists, Error[]>(
    {} as ProvidersPlaylists,
  );
  mePlaylists = new DataQuery(this.mePlaylistsQuery);

  private fetchPlaylists(authenticatedSources: SupportedSources[]) {
    if (!authenticatedSources.length) {
      this.mePlaylistsQuery.success({} as ProvidersPlaylists);
      return;
    }
    this.mePlaylistsQuery.fetching();

    const urls = authenticatedSources.map((source) => {
      return {
        url: baseUrl + this.url.replace(":source", source),
        source,
      };
    });

    const requests = urls.reduce(
      (acc, { source, url }) => {
        acc[source] = this.http.get<PlaylistsResponse>(url).pipe(
          map((playlistsResponse) =>
            playlistsResponse.map((playlist) => {
              return {
                ...playlist,
                isSelected: signal(false),
              };
            }),
          ),
          catchError((error) => {
            this.mePlaylistsQuery.fail(error);
            return [];
          }),
        );
        return acc;
      },
      {} as Record<SupportedSources, Observable<Playlist[]>>,
    );
    forkJoin(requests)
      .pipe(
        tap((fetchedPlaylists) => {
          const playlists = this.mePlaylistsQuery.data();
          fetchedPlaylists = Object.assign({}, playlists, fetchedPlaylists);
          this.mePlaylistsQuery.success(fetchedPlaylists);
        }),
        catchError((error: Error[]) => {
          this.mePlaylistsQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  refreshPlaylists() {
    const sources = this.authenticatedSources();
    this.fetchPlaylists(sources);
  }

  addProviderPlaylists(authenticatedSources: SupportedSources[]) {
    this.fetchPlaylists(authenticatedSources);
  }

  removeProviderPlaylists(removedSources: SupportedSources[]) {
    this.mePlaylistsQuery.data.update((playlists) => {
      removedSources.forEach((source) => {
        delete playlists[source];
      });
      return playlists;
    });
  }

  private sourcesSubscription: Subscription;

  constructor() {
    this.sourcesSubscription = this.authenticatedSourcesObservable.subscribe(
      (authenticatedSources) => {
        filterRemovedItems(
          authenticatedSources,
          this.mePlaylistsQuery.data(),
          this.removeProviderPlaylists.bind(this),
        );
        filterAddedItems(
          authenticatedSources,
          this.mePlaylistsQuery.data(),
          this.addProviderPlaylists.bind(this),
        );
      },
    );
  }

  ngOnDestroy() {
    if (this.sourcesSubscription) this.sourcesSubscription.unsubscribe();
  }
}
