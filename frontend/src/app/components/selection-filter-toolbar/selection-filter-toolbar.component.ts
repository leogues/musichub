import { Component, input, model, output } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { OptionFilterComponent } from './option-filter/option-filter.component';
import { TextFilterComponent } from './text-filter/text-filter.component';
import { ToolbarComponent } from './toolbar/toolbar.component';

@Component({
  selector: "app-selection-filter-toolbar",
  standalone: true,
  imports: [
    ToolbarComponent,
    OptionFilterComponent,
    TextFilterComponent,
    FormsModule,
  ],
  templateUrl: "./selection-filter-toolbar.component.html",
})
export class SelectionFilterToolbarComponent {
  sourcesToFilter = input<string[]>([]);
  isAllSelected = input.required<boolean>();
  textFilter = model.required<string>({});
  sourcesSelected = output<string[]>();

  toggleAll = output();

  toggleAllSelected() {
    this.toggleAll.emit();
  }
}
