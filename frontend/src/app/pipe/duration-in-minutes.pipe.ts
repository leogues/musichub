import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
  name: "durationInMinutes",
  standalone: true,
})
export class DurationInMinutesPipe implements PipeTransform {
  transform(durationMs: number): string {
    const minutes: number = Math.floor(durationMs / 60000);
    const seconds: number = Math.floor((durationMs % 60000) / 1000);

    const formattedMinutes = this.formatNumber(minutes);
    const formattedSeconds = this.formatNumber(seconds);

    return `${formattedMinutes}:${formattedSeconds}`;
  }

  private formatNumber(value: number): string {
    return value < 10 ? `0${value}` : `${value}`;
  }
}
