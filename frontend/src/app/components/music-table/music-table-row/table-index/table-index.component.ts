import { Component, input } from '@angular/core';
import { IndexContent } from '@components/music-table/music-table';
import { HumanIndexPipe } from '@pipe/human-index.pipe';

@Component({
  selector: "app-table-index",
  standalone: true,
  imports: [HumanIndexPipe],
  templateUrl: "./table-index.component.html",
})
export class TableIndexComponent {
  content = input.required<IndexContent>();
}
