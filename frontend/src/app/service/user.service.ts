import { baseUrl } from "app/baseUrl";
import { catchError, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { toObservable } from "@angular/core/rxjs-interop";
import { Router } from "@angular/router";
import { User } from "@type/user";

import { APIQuery, DataQuery } from "./APIQuery";

@Injectable({ providedIn: "root" })
export class UserService {
  private url = "me";
  private http = inject(HttpClient);

  private meQuery = new APIQuery<User | null, Error>(null);
  me = new DataQuery(this.meQuery);
  isFetchingObservable = toObservable(this.me.isFetching);

  fetchMe() {
    this.meQuery.fetching();
    this.http
      .get<User>(baseUrl + this.url)
      .pipe(
        tap((user) => {
          this.meQuery.success(user);
        }),
        catchError((error) => {
          this.meQuery.fail(error);
          return [];
        }),
      )
      .subscribe();
  }

  async isUserLoggedIn() {
    if (this.me.data() !== null) {
      return true;
    }

    return new Promise<boolean>((resolve) => {
      this.isFetchingObservable.subscribe((isFetching) => {
        if (!isFetching) {
          resolve(this.me.data() !== null);
        }
      });
    });
  }

  async logout() {
    const logoutResponse = await fetch("api/auth/logout", {
      method: "POST",
      credentials: "include",
    });

    if (logoutResponse.ok) {
      window.location.href = "/login";
    }
  }

  refreshMe() {
    this.fetchMe();
  }

  constructor() {
    this.fetchMe();
  }
}
