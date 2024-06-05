import { CommonModule } from "@angular/common";
import { Component, computed, ElementRef, HostListener, signal, viewChild } from "@angular/core";
import { cn } from "@utils/cn";

import { AuthComponent } from "./auth/auth.component";
import { NavigationComponent } from "./navigation/navigation.component";

@Component({
  selector: "app-sidebar",
  standalone: true,
  imports: [NavigationComponent, AuthComponent, CommonModule],
  templateUrl: "./sidebar.component.html",
})
export class SiderbarComponent {
  menu = viewChild<ElementRef<HTMLDivElement>>("menu");
  isOpen = signal(false);

  protected sidebarMenuClass = computed(() =>
    cn("absolute z-10 ml-3.5 hidden h-full pl-2 2xl:hidden 2xl:hidden", {
      block: !this.isOpen(),
    }),
  );

  protected sidebarContentClass = computed(() =>
    cn(
      "absolute z-20 hidden h-full w-[370px] shadow-menu 2xl:static 2xl:flex 2xl:w-80",
      {
        flex: this.isOpen(),
      },
    ),
  );

  toggleSidebar(mouseEvent: MouseEvent) {
    mouseEvent.stopPropagation();
    this.isOpen.update((isOpen) => !isOpen);
    this.isOpen.set(false);
  }

  handleCloseSidebar() {
    this.isOpen.set(false);
  }

  @HostListener("document:click", ["$event"]) onClick(event: MouseEvent) {
    if (!this.menu()?.nativeElement.contains(event.target as Node)) {
      this.handleCloseSidebar();
    }
  }
}
