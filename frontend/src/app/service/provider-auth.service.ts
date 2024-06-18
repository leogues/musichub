import { suportedProviders } from "app/providers";
import axios from "axios";

import { computed, inject, Injectable, signal } from "@angular/core";
import { toObservable } from "@angular/core/rxjs-interop";
import { UserService } from "@services/user.service";
import { ProviderAuthInfo, ProviderAuthResponse } from "@type/providerAuth";

const logoutProviderAuthUrl = "/api/auth/platforms/:provider/logout";

@Injectable({
  providedIn: "root",
})
export class ProviderAuthService {
  private userService = inject(UserService);

  protected user = this.userService.me.data;
  authenticatingProviders = signal<ProviderAuthResponse[]>([]);
  authenticatedProviders = computed(() => {
    if (!this.user()?.provider_auths) return [];
    //@ts-ignore
    return this.user().provider_auths;
  });
  authenticatedProvidersSources = computed(() => {
    return this.authenticatedProviders().map((auth) => auth.source);
  });
  authenticatedSourcesObservable = toObservable(
    this.authenticatedProvidersSources,
  );
  authenticatedWithAuthenticatingProviders = computed(() => {
    return this.authenticatedProviders().concat(this.authenticatingProviders());
  });
  unauthenticatedProviders = computed(() => {
    return Object.values(suportedProviders).filter(
      (providerInfo: ProviderAuthInfo) =>
        this.isProviderAuthenticated(
          providerInfo,
          this.authenticatedWithAuthenticatingProviders(),
        ),
    );
  });

  async logoutProviderAuth(provider: ProviderAuthResponse): Promise<boolean> {
    try {
      const logoutUrl = logoutProviderAuthUrl.replace(
        ":provider",
        provider.source,
      );
      const logoutResponse = await axios.delete(logoutUrl);

      if (logoutResponse.status === 200) {
        this.userService.refreshMe();
        return true;
      }
    } catch (error) {
      console.error(error);
    }

    return false;
  }

  private isProviderAuthenticated(
    provider: ProviderAuthInfo,
    authenticated: ProviderAuthResponse[],
  ): boolean {
    return !authenticated.some((auth) => auth.source === provider.source);
  }

  refreshAuths() {
    this.userService.refreshMe();
  }

  addAuthenticatingProvider(provider: ProviderAuthInfo) {
    this.authenticatingProviders.update((authenticatedProvider) => {
      return [
        ...authenticatedProvider,
        { source: provider.source, isAuthenticating: true },
      ];
    });
  }

  removeAuthenticatingProvider(provider: ProviderAuthInfo) {
    this.authenticatingProviders.update((authenticatedProvider) => {
      return authenticatedProvider.filter(
        (auth) => auth.source !== provider.source,
      );
    });
  }
}
