import { Component } from '@angular/core';

@Component({
  selector: "app-auth-closed-api",
  standalone: true,
  imports: [],
  templateUrl: "./auth-closed-api.component.html",
})
export class AuthClosedAPIComponent {
  constructor() {
    window.close();
  }
}
