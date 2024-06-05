import { CommonModule } from "@angular/common";
import { Component, input } from "@angular/core";
import { RouterLink } from "@angular/router";

@Component({
  selector: "app-link-container",
  standalone: true,
  imports: [CommonModule, RouterLink],
  templateUrl: "./link-container.component.html",
})
export class LinkContainerComponent {
  link = input<string>();
}
