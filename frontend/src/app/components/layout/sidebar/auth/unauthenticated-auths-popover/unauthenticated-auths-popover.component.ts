import {
    Component, computed, ElementRef, HostListener, inject, model, output, signal, viewChild
} from "@angular/core";
import { FormsModule } from "@angular/forms";
import { ProviderAuthService } from "@services/provider-auth.service";
import { ProviderAuthInfo } from "@type/providerAuth";

@Component({
  selector: "app-unauthenticated-auths-popover",
  standalone: true,
  imports: [FormsModule],
  templateUrl: "./unauthenticated-auths-popover.component.html",
})
export class UnauthenticatedAuthsPopoverComponent {
  popup = viewChild<ElementRef<HTMLDivElement>>("popup");
  providerAuthService: ProviderAuthService = inject(ProviderAuthService);
  unauthenticatedProviderAuths =
    this.providerAuthService.unauthenticatedProviders;

  isOpen = model.required<boolean>();

  openAuthPopup = output<ProviderAuthInfo>();
  textFilter = signal<string>("");

  unauthenticatedProviderAuthsFiltered = computed(() => {
    return this.unauthenticatedProviderAuths().filter((provider) =>
      provider.source.toLowerCase().includes(this.textFilter().toLowerCase()),
    );
  });

  protected handleClosePopover() {
    this.isOpen.update(() => false);
  }

  protected handleOpenAuthPopup(provider: ProviderAuthInfo) {
    this.openAuthPopup.emit(provider);
  }

  @HostListener("document:click", ["$event"]) onClick(event: MouseEvent) {
    if (!this.popup()?.nativeElement.contains(event.target as Node)) {
      this.handleClosePopover();
    }
  }
}
