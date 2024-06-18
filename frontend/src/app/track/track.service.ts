import { baseUrl } from "app/baseUrl";
import { catchError, forkJoin, map, Observable, Subscription, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable, signal } from "@angular/core";
import { APIQuery, DataQuery } from "@services/APIQuery";
import { ProviderAuthService } from "@services/provider-auth.service";
import { SupportedSources } from "@type/providerAuth";
import { Track, TracksResponse } from "@type/track";
import { filterAddedItems, filterRemovedItems } from "@utils/filter";

import { ProvidersTracks } from "./track";

@Injectable({
  providedIn: "root",
})
export class TrackService {
  private providerAuthService = inject(ProviderAuthService);
  private http = inject(HttpClient);

  private url = "me/platforms/:source/track";

  private authenticatedSources =
    this.providerAuthService.authenticatedProvidersSources;
  private authenticatedSourcesObservable =
    this.providerAuthService.authenticatedSourcesObservable;

  private meTracksQuery = new APIQuery<ProvidersTracks, Error[]>(
    {} as ProvidersTracks,
  );
  meTracks = new DataQuery(this.meTracksQuery);

  private fetchMeTracks(authenticatedSources: SupportedSources[]) {
    if (!authenticatedSources.length) {
      this.meTracksQuery.success({} as ProvidersTracks);
      return;
    }
    this.meTracksQuery.fetching();
    const urls = authenticatedSources.map((source) => {
      return {
        url: baseUrl + this.url.replace(":source", source),
        source,
      };
    });
    const requests = urls.reduce(
      (acc, { url, source }) => {
        acc[source] = this.http.get<TracksResponse>(url).pipe(
          map((TracksReponse) =>
            TracksReponse.map((track) => {
              return {
                ...track,
                isSelected: signal(false),
              };
            }),
          ),
        );
        return acc;
      },
      {} as Record<SupportedSources, Observable<Track[]>>,
    );
    forkJoin(requests)
      .pipe(
        tap((fetchedTracks) => {
          const tracks = this.meTracksQuery.data();
          fetchedTracks = Object.assign({}, tracks, fetchedTracks);
          this.meTracksQuery.success(fetchedTracks);
        }),
        catchError((error) => {
          this.meTracksQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  refreshTracks() {
    const sources = this.authenticatedSources();
    this.fetchMeTracks(sources);
  }

  addProviderTracks(sources: SupportedSources[]) {
    this.fetchMeTracks(sources);
  }

  removeProviderTracks(sources: SupportedSources[]) {
    this.meTracksQuery.data.update((tracks) => {
      sources.forEach((source) => delete tracks[source]);
      return { ...tracks };
    });
  }

  private sourcesSubscription: Subscription;

  constructor() {
    this.sourcesSubscription = this.authenticatedSourcesObservable.subscribe(
      (authenticatedSources) => {
        filterRemovedItems(
          authenticatedSources,
          this.meTracksQuery.data(),
          this.removeProviderTracks.bind(this),
        );
        filterAddedItems(
          authenticatedSources,
          this.meTracksQuery.data(),
          this.addProviderTracks.bind(this),
        );
      },
    );
  }

  ngOnDestroy() {
    if (this.sourcesSubscription) this.sourcesSubscription.unsubscribe();
  }
}
