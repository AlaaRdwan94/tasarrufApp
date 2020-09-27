import { Component, Input, OnInit } from '@angular/core';
import { MatSnackBar, PageEvent } from '@angular/material';
import { MatDialog } from '@angular/material/dialog';
import { AuthService } from 'src/app/services/auth.service';
import { Admin } from 'src/app/models/admin';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss']
})
export class MainComponent {

  constructor(
    private auth: AuthService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog
  ) {}

  @Input('title') title: string;
  @Input('loading') loading = false;

  public admin: Admin = null;

    dialogRef;
  loadingBarMode = (): string => {
    return this.loading ? 'indeterminate' : '';
  }

  ngOnInit() {
    this.auth
      .getUser()
      .then(admin => {
        this.admin = admin;
      })
      .catch(err => {
        console.error(err);
        this.openSnackBar('Error loading current admin user please try again');
      });
  }

  openSnackBar(message: string) {
    this.snackBar.open(message, 'cancel', { duration: 3000 });
  }
  openDialog(id: number, firstName: string, lastName: string) {
  }
}
