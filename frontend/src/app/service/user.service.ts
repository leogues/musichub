import { baseUrl } from "app/baseUrl";
import axios from "axios";
import { catchError, tap } from "rxjs";

import { HttpClient } from "@angular/common/http";
import { inject, Injectable } from "@angular/core";
import { toObservable } from "@angular/core/rxjs-interop";
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
    try {
      const logoutResponse = await axios.post("api/auth/logout");

      if (logoutResponse.status === 200) {
        window.location.href = "/login";
      }
    } catch (error) {
      console.error(error);
    }
  }

  refreshMe() {
    this.fetchMe();
  }

  constructor() {
    this.fetchMe();
  }
}
