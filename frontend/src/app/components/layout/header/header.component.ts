import { Component, input } from '@angular/core';

@Component({
  selector: "app-header",
  standalone: true,
  imports: [],
  templateUrl: "./header.component.html",
})
export class HeaderComponent {
  title = input.required<string>();
}
