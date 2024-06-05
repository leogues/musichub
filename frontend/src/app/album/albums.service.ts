import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable, signal } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";
import { SupportedSources } from "@type/providerAuth";

import { Album, AlbumsResponse, ProvidersAlbums } from "./album";
import { filterRemovedItems, filterAddedItems } from "@utils/filter";

@Injectable({
  providedIn: "root",
})
export class AlbumsService {
  private providerAuthService = inject(ProviderAuthService);
  private http = inject(HttpClient);

  private url = "me/platforms/:source/album";

  private authenticatedSources =
    this.providerAuthService.authenticatedProvidersSources;
  private authenticatedSourcesObservable =
    this.providerAuthService.authenticatedSourcesObservable;

  private meAlbumsQuery = new APIQuery<
    ProvidersAlbums,
    Record<SupportedSources, Error>
  >({} as ProvidersAlbums);
  meAlbums = new DataQuery(this.meAlbumsQuery);

  private fetchAlbums(authenticatedSources: SupportedSources[]) {
    if (!authenticatedSources.length) {
      this.meAlbumsQuery.success({} as ProvidersAlbums);
      return;
    }

    this.meAlbumsQuery.fetching();

    const urls = authenticatedSources.map((source) => {
      return {
        url: baseUrl + this.url.replace(":source", source),
        source,
      };
    });

    const requests = urls.reduce(
      (acc, { source, url }) => {
        acc[source] = this.http.get<AlbumsResponse>(url).pipe(
          map((albumsResponse) =>
            albumsResponse.map((album) => {
              return {
                ...album,
                isSelected: signal(false),
              };
            }),
          ),
          catchError((err) => {
            const error = {} as Record<SupportedSources, Error>;
            error[source] = err;
            throw error;
          }),
        );
        return acc;
      },
      {} as Record<SupportedSources, Observable<Album[]>>,
    );

    forkJoin(requests)
      .pipe(
        tap((fetchedAlbums) => {
          const albums = this.meAlbumsQuery.data();
          fetchedAlbums = Object.assign({}, albums, fetchedAlbums);
          this.meAlbumsQuery.success(fetchedAlbums);
        }),

        catchError((error: Record<SupportedSources, Error>) => {
          this.meAlbumsQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  refreshAlbums() {
    const sources = this.authenticatedSources();
    this.fetchAlbums(sources);
  }

  addProviderAlbums(authenticatedSources: SupportedSources[]) {
    this.fetchAlbums(authenticatedSources);
  }

  removeProviderAlbums(removedSources: SupportedSources[]) {
    this.meAlbumsQuery.data.update((albums) => {
      removedSources.forEach((source) => delete albums[source]);
      return albums;
    });
  }

  private sourcesSubscription: Subscription;

  constructor() {
    this.sourcesSubscription = this.authenticatedSourcesObservable.subscribe(
      (authenticatedSources) => {
        filterRemovedItems(
          authenticatedSources,
          this.meAlbumsQuery.data(),
          this.removeProviderAlbums.bind(this),
        );
        filterAddedItems(
          authenticatedSources,
          this.meAlbumsQuery.data(),
          this.addProviderAlbums.bind(this),
        );
      },
    );
  }

  ngOnDestroy() {
    if (this.sourcesSubscription) this.sourcesSubscription.unsubscribe();
  }
}
