import { Component, computed, inject } from '@angular/core';
import { FirstWordPipe } from '@pipe/first-word.pipe';
import { UserService } from '@services/user.service';

@Component({
  selector: "app-profile",
  standalone: true,
  imports: [FirstWordPipe],
  templateUrl: "./profile.component.html",
})
export class ProfileComponent {
  userService = inject(UserService);
  me = this.userService.me;
  user = this.me.data;

  name = computed(() => {
    return this.user()?.name || "";
  });

  logout() {
    this.userService.logout();
  }
}
