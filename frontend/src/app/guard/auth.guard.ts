import { inject } from "@angular/core";
import { CanActivateFn, Router } from "@angular/router";
import { UserService } from "@services/user.service";

export const authGuard: CanActivateFn = async (route, state) => {
  const router = inject(Router);
  const userService = inject(UserService);

  const isUserLoggedIn = await userService.isUserLoggedIn();

  if (isUserLoggedIn) {
    return true;
  } else {
    router.navigate(["login"], {
      queryParams: { redirect: state.url },
    });
    return false;
  }
};
