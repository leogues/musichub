import { Routes } from "@angular/router";

import { TrackComponent } from "./pages/track/track.component";

export const TRACK_ROUTES: Routes = [
  { path: "", component: TrackComponent },
  { path: ":provider", component: TrackComponent },
];
