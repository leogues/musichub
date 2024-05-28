import {
    Component, computed, ElementRef, HostListener, input, model, viewChild
} from '@angular/core';
import { CapitalizeFirstLetterPipe } from '@pipe/capitalize-first-letter.pipe';

import { Source } from '../source';

@Component({
  selector: "app-popup-select-option",
  standalone: true,
  imports: [CapitalizeFirstLetterPipe],
  templateUrl: "./popup-select-option.component.html",
})
export class PopupSelectOptionComponent {
  selectOption = viewChild<ElementRef<HTMLDivElement>>("selectOption");
  isOpen = model.required<boolean>();
  sourcesToFilter = input.required<Source[]>();

  isAllSelected = computed(() => {
    return this.sourcesToFilter().every((source) => source.isSelected());
  });

  handleCloseSidebar() {
    this.isOpen.set(false);
  }

  handleSelect(source: Source) {
    source.isSelected.set(!source.isSelected());
  }

  handleSelectAll() {
    const isAllSelected = this.isAllSelected();

    this.sourcesToFilter().forEach((source) => {
      return source.isSelected.set(!isAllSelected);
    });
  }

  @HostListener("document:click", ["$event"]) onClick(event: MouseEvent) {
    if (!this.selectOption()?.nativeElement.contains(event.target as Node)) {
      this.handleCloseSidebar();
    }
  }
}
