import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "humanIndex",
  standalone: true,
})
export class HumanIndexPipe implements PipeTransform {
  transform(index: number): string {
    index += 1;

    const formatedIndex = this.formatNumber(index);

    return formatedIndex;
  }

  private formatNumber(value: number): string {
    return value < 10 ? `0${value}` : `${value}`;
  }
}
