import { Component } from "@angular/core";

import { NavigationItemComponent } from "./navigation-item/navigation-item.component";
import { ProfileComponent } from "./profile/profile.component";

@Component({
  selector: "app-navigation",
  standalone: true,
  imports: [ProfileComponent, NavigationItemComponent],
  templateUrl: "./navigation.component.html",
})
export class NavigationComponent {}
