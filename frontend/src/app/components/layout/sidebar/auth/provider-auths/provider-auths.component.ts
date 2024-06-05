import { Component, inject, model, signal } from "@angular/core";
import { RouterLink } from "@angular/router";
import { ProviderAuthService } from "@services/provider-auth.service";
import { ProviderAuthResponse } from "@type/providerAuth";

import { ProviderAuthPopupComponent } from "./provider-auth-popup/provider-auth-popup.component";
import { ProviderAuthComponent } from "./provider-auth/provider-auth.component";

@Component({
  selector: "app-provider-auths",
  standalone: true,
  imports: [ProviderAuthComponent, RouterLink, ProviderAuthPopupComponent],
  templateUrl: "./provider-auths.component.html",
})
export class ProviderAuthsComponent {
  private providerAuthService: ProviderAuthService =
    inject(ProviderAuthService);

  isOpenPopover = model.required<boolean>();
  isOpenAuthPopup = signal(false);
  mouseScrollY = signal<number>(0);
  providerAuthSelected = signal<ProviderAuthResponse | null>(null);

  authenticatedWithAuthenticatingProviders =
    this.providerAuthService.authenticatedWithAuthenticatingProviders;

  unauthenticatedProviders = this.providerAuthService.unauthenticatedProviders;

  togglePopover(event: MouseEvent) {
    event.stopPropagation();
    this.isOpenPopover.update((isOpen) => !isOpen);
  }

  toggleAuthPopup(event: MouseEvent, providerAuth: ProviderAuthResponse) {
    if (providerAuth.isAuthenticating) {
      return;
    }

    const buttonElement = event.target as HTMLButtonElement;
    const rect = buttonElement
      .closest(".provider-auth")!
      .getBoundingClientRect();
    event.stopPropagation();
    this.isOpenAuthPopup.update((isOpen) => !isOpen);
    this.mouseScrollY.set(rect.top);
    this.providerAuthSelected.set(providerAuth);
  }
}
