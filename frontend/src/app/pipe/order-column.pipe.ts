import { Pipe, PipeTransform } from "@angular/core";
import { TData, THeaderWithOrder } from "@components/music-table/music-table";

type ColumnType = {
  type: string;
  text?: string;
  index?: number;
  title?: string;
  duration?: number;
};

@Pipe({
  name: "orderColumn",
  standalone: true,
})
export class OrderColumnPipe implements PipeTransform {
  private extractValue<T>(item: T, index: number): string | number | undefined {
    const column = (item as any).content[index] as ColumnType;
    switch (column.type) {
      case "text":
        return column.text?.toUpperCase();
      case "index":
        return column.index;
      case "title":
        return column.title?.toUpperCase();
      case "duration":
        return column.duration;
      default:
        return undefined;
    }
  }

  private compareStrings(a: string, b: string, order: "asc" | "desc"): number {
    return order === "asc" ? a.localeCompare(b) : b.localeCompare(a);
  }

  transform(
    data: TData[],
    columnWillOrdered: THeaderWithOrder | null,
  ): TData[] {
    if (!columnWillOrdered || !columnWillOrdered.order) {
      return data;
    }

    const sortedData = [...data];

    return sortedData.sort((a, b) => {
      const valueA = this.extractValue(a, columnWillOrdered.index);
      const valueB = this.extractValue(b, columnWillOrdered.index);

      if (typeof valueA === "string" && typeof valueB === "string") {
        return this.compareStrings(
          valueA,
          valueB,
          columnWillOrdered.order as "asc" | "desc",
        );
      } else if (typeof valueA === "number" && typeof valueB === "number") {
        return columnWillOrdered.order === "asc"
          ? valueA - valueB
          : valueB - valueA;
      }

      return 0;
    });
  }
}
