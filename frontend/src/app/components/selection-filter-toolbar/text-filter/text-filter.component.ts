import { Component, model } from "@angular/core";
import { FormsModule } from "@angular/forms";

@Component({
  selector: "app-text-filter",
  standalone: true,
  imports: [FormsModule],
  templateUrl: "./text-filter.component.html",
})
export class TextFilterComponent {
  filter = model.required<string>({
    alias: "filterText",
  });
}
