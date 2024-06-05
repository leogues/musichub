import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "unescapeString",
  standalone: true,
})
export class UnescapeStringPipe implements PipeTransform {
  transform(value: string): string {
    const parser = new DOMParser();
    const dom = parser.parseFromString(value, "text/html");

    return dom.body.textContent || "";
  }
}
