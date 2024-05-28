import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: "firstWord",
  standalone: true,
})
export class FirstWordPipe implements PipeTransform {
  transform(value: string): string {
    return value.split(" ")[0];
  }
}
