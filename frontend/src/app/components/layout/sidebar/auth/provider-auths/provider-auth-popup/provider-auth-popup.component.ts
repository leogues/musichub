import {
    Component, effect, ElementRef, HostListener, inject, input, model, viewChild
} from "@angular/core";
import { CapitalizeFirstLetterPipe } from "@pipe/capitalize-first-letter.pipe";
import { ProviderAuthService } from "@services/provider-auth.service";
import { ProviderAuthResponse } from "@type/providerAuth";

@Component({
  selector: "app-provider-auth-popup",
  standalone: true,
  imports: [CapitalizeFirstLetterPipe],
  templateUrl: "./provider-auth-popup.component.html",
})
export class ProviderAuthPopupComponent {
  providerAuthService = inject(ProviderAuthService);
  isOpenPopup = model.required<boolean>();
  mouseScrollY = input.required<number>();
  popup = viewChild<ElementRef<HTMLDivElement>>("popup");
  providerAuth = input.required<ProviderAuthResponse | null>();

  handleClosePopover() {
    this.isOpenPopup.update(() => false);
  }

  constructor() {
    effect(() => {
      const popupElement = this.popup()?.nativeElement;
      if (popupElement) {
        popupElement.style.top = `${this.mouseScrollY()}px`;
      }
    });
  }

  async handleDisconnectProvider() {
    const isLogout = await this.providerAuthService.logoutProviderAuth(
      this.providerAuth()!,
    );
    if (isLogout) this.handleClosePopover();
  }

  @HostListener("document:click", ["$event"]) onClick(event: MouseEvent) {
    if (!this.popup()?.nativeElement.contains(event.target as Node)) {
      this.handleClosePopover();
    }
  }
}
