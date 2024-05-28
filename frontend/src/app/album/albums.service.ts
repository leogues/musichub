import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable, signal } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";
import { SupportedSources } from "@type/providerAuth";

import { Album, AlbumsResponse, ProvidersAlbums } from "./album";

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

  private meAlbumsQuery = new APIQuery<ProvidersAlbums, Error[]>(
    {} as ProvidersAlbums,
  );
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

        catchError((error) => {
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

  private removedSourcesSubscription: Subscription;
  private addedSourcesSubscription: Subscription;

  constructor() {
    this.removedSourcesSubscription =
      this.authenticatedSourcesObservable.subscribe((authenticatedSources) => {
        const removedSources = Object.keys(this.meAlbumsQuery.data()).filter(
          (source) =>
            !authenticatedSources.includes(source as SupportedSources),
        );

        if (removedSources.length === 0) return;
        this.removeProviderAlbums(removedSources as SupportedSources[]);
      });

    this.addedSourcesSubscription =
      this.authenticatedSourcesObservable.subscribe((authenticatedSources) => {
        const addedSources = authenticatedSources.filter(
          (source) => !Object.keys(this.meAlbumsQuery.data()).includes(source),
        );
        if (addedSources.length === 0) return;
        this.addProviderAlbums(addedSources);
      });
  }

  ngOnDestroy() {
    if (this.removedSourcesSubscription)
      this.removedSourcesSubscription.unsubscribe();
    if (this.addedSourcesSubscription)
      this.addedSourcesSubscription.unsubscribe();
  }
}
