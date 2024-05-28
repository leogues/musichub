import { Routes } from '@angular/router';

import { ArtistsComponent } from './page/artists/artists.component';

export const ARTIST_ROUTES: Routes = [
  { path: "", component: ArtistsComponent },
  { path: ":provider", component: ArtistsComponent },
];
