import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable, signal } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";
import { SupportedSources } from "@type/providerAuth";

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
        acc[source] = this.http.get<ArtistsResponse>(baseUrl + url).pipe(
          map((artistsResponse) =>
            artistsResponse.map((artist) => {
              return {
                ...artist,
                isSelected: signal(false),
              };
            }),
          ),
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
      return artists;
    });
  }

  private removedSourcesSubscription: Subscription;
  private addedSourcesSubscription: Subscription;

  constructor() {
    this.removedSourcesSubscription =
      this.authenticatedSourcesObservable.subscribe((authenticatedSources) => {
        const removedSources = Object.keys(this.meArtistsQuery.data()).filter(
          (source) =>
            !authenticatedSources.includes(source as SupportedSources),
        );

        if (removedSources.length === 0) return;
        this.removeProviderAlbums(removedSources as SupportedSources[]);
      });

    this.addedSourcesSubscription =
      this.authenticatedSourcesObservable.subscribe((authenticatedSources) => {
        const addedSources = authenticatedSources.filter(
          (source) => !Object.keys(this.meArtistsQuery.data()).includes(source),
        );
        if (addedSources.length === 0) return;
        this.addProviderAlbums(addedSources);
      });
  }

  ngOnDestrou() {
    if (this.removedSourcesSubscription)
      this.removedSourcesSubscription.unsubscribe();
    if (this.addedSourcesSubscription)
      this.addedSourcesSubscription.unsubscribe();
  }
}
