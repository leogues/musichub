import { Component, input } from "@angular/core";
import { TitleContent } from "@components/music-table/music-table";

@Component({
  selector: "app-table-title",
  standalone: true,
  imports: [],
  templateUrl: "./table-title.component.html",
})
export class TableTitleComponent {
  data = input.required<TitleContent>();
}
