import { Component, input } from '@angular/core';
import { SelectContent } from '@components/music-table/music-table';

@Component({
  selector: "app-table-select",
  standalone: true,
  imports: [],
  templateUrl: "./table-select.component.html",
})
export class TableSelectComponent {
  content = input.required<SelectContent>();

  toggleSelect(event: MouseEvent) {
    this.content().isSelected.update((isSelected) => !isSelected);
    event.stopPropagation();
    event.preventDefault();
  }
}
