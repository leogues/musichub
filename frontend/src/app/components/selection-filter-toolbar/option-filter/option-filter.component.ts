import { Component, computed, effect, input, output, signal } from "@angular/core";

import { PopupSelectOptionComponent } from "./popup-select-option/popup-select-option.component";
import { Source } from "./source";

@Component({
  selector: "app-option-filter",
  standalone: true,
  imports: [PopupSelectOptionComponent],
  templateUrl: "./option-filter.component.html",
})
export class OptionFilterComponent {
  sourcesToFilter = input.required<Source[], string[]>({
    transform: (sources) =>
      sources.map(
        (source) => ({ name: source, isSelected: signal(true) }) as Source,
      ),
  });
  sourcesSelectedOutput = output<string[]>();

  isOpenPopup = signal(false);

  sourcesSelectedCount = computed(
    () => this.sourcesToFilter().filter((source) => source.isSelected()).length,
  );

  togglePopup(mouseEvent: MouseEvent) {
    mouseEvent.stopPropagation();
    this.isOpenPopup.update((isOpen) => !isOpen);
  }
  constructor() {
    effect(() => {
      const selectedSources = this.sourcesToFilter()
        .filter((source) => source.isSelected())
        .map((source) => source.name);
      this.sourcesSelectedOutput.emit(selectedSources);
    });
  }
}
