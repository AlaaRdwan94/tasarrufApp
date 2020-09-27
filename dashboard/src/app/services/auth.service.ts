import { Injectable } from "@angular/core";
import { HttpClient, HttpErrorResponse } from "@angular/common/http";
import { environment } from "./../../environments/environment";
import { map, retry } from "rxjs/operators";
import { Admin } from "../models/admin";

/**
 * LoginResponse is the expected response for admin login request.
 */
interface LoginResponse {
  success: string;
  message: string;
  token: string;
  user: Admin;
}

@Injectable({
  providedIn: "root"
})

/**
 * The admin authentication service.
 */
export class AuthService {
  constructor(private http: HttpClient) {
    this.rootUrl = environment.apiUrl;
  }

  public loggedIn = false;
  private currentAdmin: Admin = null;
  private token: string = null;
  private rootUrl: string;

  /**
   * Send a login request to the backend.
   * Only admins can login with this backend request.
   * @param email [string]--- email of the user.
   * @param password [string]--- password of the user.
   */
  public logIn(email: string, password: string) {
    return new Promise<boolean>((accept, reject) => {
      this.http
        .post(`${this.rootUrl}/public/email-login`, {
          email,
          password
        })
        .toPromise()
        .then((res: LoginResponse) => {
          if (res.user.accountType === "admin") {
            console.log("Logged in", res);
            this.loggedIn = true;
            this.setToken(res.token);
            this.currentAdmin = res.user;
            accept(true);
          } else {
            reject(false);
          }
        })
        .catch((err: HttpErrorResponse) => {
          console.error("Error Logging User In : ", err.error);
          reject(false);
        });
    });
  }

  public recoverPassword(mobile: string) {
    return this.http
      .post(`${this.rootUrl}/public/admin/forget-password`, {
        mobile
      })
      .pipe(retry(3));
  }

  /**
   * isLoggedIn returns true if the current user is logged in.
   */
  public isLoggedIn(): boolean {
    const token: string = this.getToken();
    return !(token == null);
  }

  private setUser(): Promise<Admin> {
    return new Promise<Admin>((accept, reject) => {
      this.http
        .get(`${this.rootUrl}/user`, { headers: { Token: this.getToken() } })
        .toPromise()
        .then((res: LoginResponse) => {
          this.currentAdmin = res.user;
          accept(res.user);
        })
        .catch((err: HttpErrorResponse) => {
          console.error("Error Getting User: ", err.error);
          reject(err);
        });
    });
  }

  /**
   * getToken returns the token of the current logged in admin.
   */
  public getToken(): string {
    if (this.token == null) {
      return window.localStorage.getItem("token");
    } else {
      return this.token;
    }
  }

  /**
   * setToken sets the token value in local storage and auth service.
   * @param token [string]--- authentication token.
   */
  private setToken(token: string): void {
    this.token = token;
    window.localStorage.setItem("token", token);
  }
  /**
   * returns the profile of the current logged in admin.
   */
  public getUser(): Promise<Admin> {
    return new Promise<Admin>((accept, reject) => {
      if (this.currentAdmin) {
        accept(this.currentAdmin);
      } else {
        this.setUser()
          .then(admin => {
            accept(admin);
          })
          .catch(err => {
            reject(err);
          });
      }
    });
  }
}
