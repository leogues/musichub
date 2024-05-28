import { WritableSignal } from '@angular/core';

export type Source = { name: string; isSelected: WritableSignal<boolean> };
