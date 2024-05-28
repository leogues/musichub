import { Routes } from '@angular/router';

import { AlbumComponent } from './page/album/album.component';
import { AlbumsComponent } from './page/albums/albums.component';

export const ALBUM_ROUTES: Routes = [
  { path: "", component: AlbumsComponent },
  { path: ":provider", component: AlbumsComponent },
  {
    path: ":provider/:id",
    component: AlbumComponent,
  },
];
