import {
  Component,
  ElementRef,
  HostListener,
  effect,
  input,
  model,
  signal,
  viewChild,
} from "@angular/core";
import { CapitalizeFirstLetterPipe } from "@pipe/capitalize-first-letter.pipe";
import { ProviderAuthResponse } from "@type/providerAuth";

@Component({
  selector: "app-provider-auth-popup",
  standalone: true,
  imports: [CapitalizeFirstLetterPipe],
  templateUrl: "./provider-auth-popup.component.html",
})
export class ProviderAuthPopupComponent {
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

  @HostListener("document:click", ["$event"]) onClick(event: MouseEvent) {
    if (!this.popup()?.nativeElement.contains(event.target as Node)) {
      this.handleClosePopover();
    }
  }
}
