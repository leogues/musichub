import { CommonModule } from "@angular/common";
import { Component, input, model } from "@angular/core";
import { RouterLink, RouterLinkActive } from "@angular/router";

@Component({
  selector: "app-navigation-item",
  standalone: true,
  imports: [RouterLink, RouterLinkActive, CommonModule],
  templateUrl: "./navigation-item.component.html",
})
export class NavigationItemComponent {
  icon = input<string>("");
  name = input.required<string>({ alias: "routeName" });
  routerLink = input<string>();
  afterIcon = input<string>("");
  isOpen = model<boolean>(false);
  isDisabled = input<boolean>(false);

  handleClick() {
    if (this.isDisabled()) return;
    this.isOpen.update((value) => !value);
  }
}
