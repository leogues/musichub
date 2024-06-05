import { Component, input } from "@angular/core";
import { DurationContent } from "@components/music-table/music-table";
import { DurationInMinutesPipe } from "@pipe/duration-in-minutes.pipe";

@Component({
  selector: "app-table-duration",
  standalone: true,
  imports: [DurationInMinutesPipe],
  templateUrl: "./table-duration.component.html",
})
export class TableDurationComponent {
  data = input.required<DurationContent>();
}
