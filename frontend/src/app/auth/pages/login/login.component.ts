import { Component, inject } from "@angular/core";
import { ActivatedRoute } from "@angular/router";

const baseAuthUrl = "/api/auth/:provider?redirect=:redirect";

@Component({
  selector: "app-login",
  standalone: true,
  imports: [],
  templateUrl: "./login.component.html",
})
export class LoginComponent {
  activatedRoute = inject(ActivatedRoute);

  login(provider: string) {
    let redirectUrl = this.activatedRoute.snapshot.queryParams["redirect"];
    if (!redirectUrl) redirectUrl = "/";

    const authUrl = baseAuthUrl
      .replace(":provider", provider)
      .replace(":redirect", redirectUrl);

    window.location.href = authUrl;
  }
}
