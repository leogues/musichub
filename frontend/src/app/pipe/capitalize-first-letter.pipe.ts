import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "capitalizeFirstLetter",
  standalone: true,
})
export class CapitalizeFirstLetterPipe implements PipeTransform {
  transform(value: string | null | undefined) {
    if (!value) {
      return value;
    }
    return value.charAt(0).toUpperCase() + value.slice(1);
  }
}
