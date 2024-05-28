import { suportedProviders } from "app/providers";

import { CommonModule } from "@angular/common";
import { Component, Input } from "@angular/core";
import { ProviderAuth, ProviderAuthResponse } from "@type/providerAuth";

@Component({
  selector: "app-provider-auth",
  standalone: true,
  imports: [CommonModule],
  templateUrl: "./provider-auth.component.html",
})
export class ProviderAuthComponent {
  @Input({
    required: true,
    transform: (authenticatedAuth: ProviderAuthResponse) => {
      const providerInfo = suportedProviders[authenticatedAuth.source];
      return {
        ...authenticatedAuth,
        picture: providerInfo.picture,
        label: providerInfo.label,
      };
    },
  })
  authenticatedAuth!: ProviderAuth;
}
