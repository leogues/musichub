import { CommonModule } from "@angular/common";
import { Component, input } from "@angular/core";
import { TextContent } from "@components/music-table/music-table";
import { CapitalizeFirstLetterPipe } from "@pipe/capitalize-first-letter.pipe";

@Component({
  selector: "app-table-text",
  standalone: true,
  imports: [CapitalizeFirstLetterPipe, CommonModule],
  templateUrl: "./table-text.component.html",
})
export class TableTextComponent {
  data = input.required<TextContent>();
}
