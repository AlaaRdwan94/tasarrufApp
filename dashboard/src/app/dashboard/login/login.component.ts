import { Component, OnInit } from "@angular/core";
import { NgForm } from "@angular/forms";
import { AuthService } from "src/app/services/auth.service";
import { MatSnackBar } from "@angular/material/snack-bar";
import { Router } from "@angular/router";

@Component({
  selector: "app-login",
  templateUrl: "./login.component.html",
  styleUrls: ["./login.component.scss"]
})
export class LoginComponent implements OnInit {
  loading = false;
  recoverPasswordState = false;
  lodingIndicatorColor = "primary";
  credentials: any = {};
  recoveryPhone = "";
  constructor(
    private auth: AuthService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit() {}

  openSnackBar(message: string) {
    this.snackBar.open(message, "cancel", { duration: 3000 });
  }

  login(form: NgForm) {
    const loginInfo = form.value;
    this.loading = true;
    this.auth
      .logIn(loginInfo.email, loginInfo.password)
      .then((ok: boolean) => {
        console.log("Login :", ok);
        this.loading = false;
        this.openSnackBar("Login Successfull");
        console.log("admin : ", this.auth.getUser());
        console.log("token : ", this.auth.getToken());
        this.router.navigate(["/admin"]);
      })
      .catch((ok: boolean) => {
        this.loading = false;
        this.openSnackBar("Wrong Credentials");
        console.log("Login : ", ok);
      });
  }
  forgetPassword() {
    this.recoverPasswordState = true;
  }

  sendRecoverPassword(form: NgForm) {
    const { mobile } = form.value;
    this.auth.recoverPassword(mobile).subscribe(
      (res: any) => {
        console.log(res);
        this.openSnackBar(res.success);
      },
      err => {
        this.openSnackBar(
          "Phone number is not associated with any admin account"
        );
      }
    );
  }
}
