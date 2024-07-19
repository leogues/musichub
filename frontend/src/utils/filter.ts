import { signal, WritableSignal } from "@angular/core";

export const filterRemovedItems = <T extends string>(
  array: T[],
  obj: Record<T, any>,
  callback: (item: T[]) => void,
) => {
  const removedItems = Object.keys(obj).filter(
    (item) => !array.includes(item as T),
  );
  if (removedItems.length === 0) return;
  callback(removedItems as T[]);
};

export const filterAddedItems = <T extends string>(
  array: T[],
  obj: Record<T, any>,
  callback: (item: T[]) => void,
) => {
  const addedItems = array.filter((item) => !Object.keys(obj).includes(item));
  if (addedItems.length === 0) return;
  callback(addedItems as T[]);
};

export const addPropertyIsSelected = <T extends object>(
  items: T[],
): (T & { isSelected: WritableSignal<boolean> })[] => {
  return items.map((item) => ({
    ...item,
    isSelected: signal(false),
  }));
};
