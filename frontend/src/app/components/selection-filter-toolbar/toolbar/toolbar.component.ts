import { Component, effect, input, output } from "@angular/core";
import { FormsModule } from "@angular/forms";

@Component({
  selector: "app-toolbar",
  standalone: true,
  imports: [FormsModule],
  templateUrl: "./toolbar.component.html",
})
export class ToolbarComponent {
  isAllSelected = input.required<boolean>();
  toggleAll = output();

  toggleAllSelected() {
    this.toggleAll.emit();
  }
}
