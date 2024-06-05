import { CommonModule } from "@angular/common";
import { Component, inject, input, output } from "@angular/core";
import { UserService } from "@services/user.service";

import { HeaderComponent } from "./header/header.component";
import { SiderbarComponent } from "./sidebar/sidebar.component";

@Component({
  selector: "app-layout",
  standalone: true,
  imports: [HeaderComponent, SiderbarComponent, CommonModule],
  templateUrl: "./layout.component.html",
})
export class LayoutComponent {
  title = input.required<string>();
  userService = inject(UserService);
  refreshData = output<void>();

  handleRefreshData() {
    this.refreshData.emit();
  }
}
