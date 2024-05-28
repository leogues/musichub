import { Routes } from '@angular/router';

import { PlaylistComponent } from './pages/playlist/playlist.component';
import { PlaylistsComponent } from './pages/playlists/playlists.component';

export const PLAYLIST_ROUTES: Routes = [
  { path: "", component: PlaylistsComponent },
  { path: ":provider", component: PlaylistsComponent },
  { path: ":provider/:id", component: PlaylistComponent },
];
