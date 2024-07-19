import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";
import { SupportedSources } from "@type/providerAuth";
import { addPropertyIsSelected, filterAddedItems, filterRemovedItems } from "@utils/filter";

import { Artist, ArtistsResponse, ProvidersArtists } from "./artist";

@Injectable({
  providedIn: "root",
})
export class ArtistService {
  private providerAuthService = inject(ProviderAuthService);
  private http = inject(HttpClient);

  private url = "me/platforms/:source/artist";

  private authenticatedSources =
    this.providerAuthService.authenticatedProvidersSources;
  private authenticatedSourcesObservable =
    this.providerAuthService.authenticatedSourcesObservable;

  private meArtistsQuery = new APIQuery<ProvidersArtists, Error[]>(
    {} as ProvidersArtists,
  );
  meArtists = new DataQuery(this.meArtistsQuery);

  private fetchArtists(authenticatedSources: SupportedSources[]) {
    if (!authenticatedSources.length) {
      this.meArtistsQuery.success({} as ProvidersArtists);
      return;
    }

    this.meArtistsQuery.fetching();

    const urls = authenticatedSources.map((source) => {
      return {
        url: this.url.replace(":source", source),
        source,
      };
    });

    const requests = urls.reduce(
      (acc, { source, url }) => {
        acc[source] = this.http
          .get<ArtistsResponse>(baseUrl + url)
          .pipe(
            map((artistsResponse) => addPropertyIsSelected(artistsResponse)),
          );
        return acc;
      },
      {} as Record<SupportedSources, Observable<Artist[]>>,
    );

    forkJoin(requests)
      .pipe(
        tap((fetchedArtists) => {
          const artists = this.meArtistsQuery.data();
          fetchedArtists = Object.assign({}, artists, fetchedArtists);
          this.meArtistsQuery.success(fetchedArtists);
        }),
        catchError((error: Error[]) => {
          this.meArtistsQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  refreshArtists() {
    const sources = this.authenticatedSources();
    this.fetchArtists(sources);
  }

  addProviderAlbums(authenticatedSources: SupportedSources[]) {
    this.fetchArtists(authenticatedSources);
  }

  removeProviderAlbums(removedSources: SupportedSources[]) {
    this.meArtistsQuery.data.update((artists) => {
      removedSources.forEach((source) => delete artists[source]);
      return { ...artists };
    });
  }

  private sourcesSubscription: Subscription;

  constructor() {
    this.sourcesSubscription = this.authenticatedSourcesObservable.subscribe(
      (authenticatedSources) => {
        filterRemovedItems(
          authenticatedSources,
          this.meArtistsQuery.data(),
          this.removeProviderAlbums.bind(this),
        );
        filterAddedItems(
          authenticatedSources,
          this.meArtistsQuery.data(),
          this.addProviderAlbums.bind(this),
        );
      },
    );
  }

  ngOnDestroy() {
    if (this.sourcesSubscription) this.sourcesSubscription.unsubscribe();
  }
}
