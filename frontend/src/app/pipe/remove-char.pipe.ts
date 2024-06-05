import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "removeChar",
  standalone: true,
})
export class RemoveCharPipe implements PipeTransform {
  transform(value: string, removeChar: string): string {
    return value.replace(removeChar, "");
  }
}
