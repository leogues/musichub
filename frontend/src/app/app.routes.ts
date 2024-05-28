import { Routes } from "@angular/router";

import { AuthClosedAPIComponent } from "./auth/pages/auth-closed-api/auth-closed-api.component";
import { LoginComponent } from "./auth/pages/login/login.component";
import { authGuard } from "./guard/auth.guard";

export const routes: Routes = [
  { path: "login", component: LoginComponent },
  { path: "authClosedAPI", component: AuthClosedAPIComponent },
  {
    path: "",
    redirectTo: "playlist",
    pathMatch: "full",
  },
  {
    path: "track",
    loadChildren: () =>
      import("./track/track.routes").then((m) => m.TRACK_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "playlist",
    loadChildren: () =>
      import("./playlist/playlist.routes").then((m) => m.PLAYLIST_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "album",
    loadChildren: () =>
      import("./album/album.routes").then((m) => m.ALBUM_ROUTES),
    canActivate: [authGuard],
  },
  {
    path: "artist",
    loadChildren: () =>
      import("./artist/artist.routes").then((m) => m.ARTIST_ROUTES),
    canActivate: [authGuard],
  },
];
