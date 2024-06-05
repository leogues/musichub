import { Component, input, WritableSignal } from "@angular/core";

@Component({
  selector: "app-table-select",
  standalone: true,
  imports: [],
  templateUrl: "./table-select.component.html",
})
export class TableSelectComponent {
  content = input.required<
    { isSelected: WritableSignal<boolean> },
    WritableSignal<boolean>
  >({
    transform: (isSelected) => {
      return {
        isSelected,
      };
    },
  });
  toggleSelect(event: MouseEvent) {
    this.content().isSelected.update((isSelected) => !isSelected);
    event.stopPropagation();
    event.preventDefault();
  }
}
