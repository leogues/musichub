import { CommonModule } from "@angular/common";
import {
  Component,
  computed,
  effect,
  input,
  model,
  output,
  viewChild,
} from "@angular/core";
import { CapitalizeFirstLetterPipe } from "@pipe/capitalize-first-letter.pipe";
import { RemoveCharPipe } from "@pipe/remove-char.pipe";
import { TruncatePipe } from "@pipe/truncate.pipe";
import { UnescapeStringPipe } from "@pipe/unescape-string.pipe";
import { SupportedSources } from "@type/providerAuth";
import { getProviderTextColor } from "@utils/getProviderStyle";

@Component({
  selector: "app-detail",
  standalone: true,
  imports: [
    CommonModule,
    CapitalizeFirstLetterPipe,
    TruncatePipe,
    UnescapeStringPipe,
    RemoveCharPipe,
  ],
  templateUrl: "./detail.component.html",
})
export class DetailComponent {
  title = input.required<string>();
  picture = input.required<string>();
  link = input.required<string>();
  author = input.required<string>();
  authorLink = input.required<string>();
  platform = input.required<string>();
  isPublic = input<boolean>();
  description = input<string>();
  tracksCount = input.required<number>();
  followersCount = input<number>();
  releaseDate = input<string>();

  platformColorClass = computed(() => {
    return getProviderTextColor(this.platform() as SupportedSources);
  });
}
