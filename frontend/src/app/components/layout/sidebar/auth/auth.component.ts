import { Component, inject, signal } from "@angular/core";
import { ProviderAuthService } from "@services/provider-auth.service";
import { ProviderAuthInfo } from "@type/providerAuth";

import { ProviderAuthsComponent } from "./provider-auths/provider-auths.component";
import { UnauthenticatedAuthsPopoverComponent } from "./unauthenticated-auths-popover/unauthenticated-auths-popover.component";

@Component({
  selector: "app-auth",
  standalone: true,
  imports: [ProviderAuthsComponent, UnauthenticatedAuthsPopoverComponent],
  templateUrl: "./auth.component.html",
})
export class AuthComponent {
  isPopoverOpen = signal<boolean>(false);
  private providerAuthService: ProviderAuthService =
    inject(ProviderAuthService);

  openAuthPopup(provider: ProviderAuthInfo) {
    this.providerAuthService.addAuthenticatingProvider(provider);

    const authWindow = window.open(
      "/api/auth/" + provider.source,
      "_blank",
      "width=480,height=700",
    );

    if (!authWindow) {
      return alert("Please disable your popup blocker and try again.");
    }

    const interval = setInterval(() => {
      if (authWindow.closed) {
        clearInterval(interval);
        this.providerAuthService.removeAuthenticatingProvider(provider);
        this.providerAuthService.refreshAuths();
      }
    }, 200);
  }
}
