import { Component, OnInit } from '@angular/core';
import { UsersService } from 'src/app/services/users.service';
import { MatSnackBar} from '@angular/material';
import { User } from 'src/app/models/user';
import {error} from '@angular/compiler/src/util';



@Component({
  selector: 'app-not-approved',
  templateUrl: './not-approved.component.html',
  styleUrls: ['./not-approved.component.scss']
})
export class NotApprovedComponent implements OnInit {
  loadingState = false;
    usersCount = 0;
    paginatorPagesCount = 0;
    paginatorPageIndex = 0;
    displayedColumns: string[] = [
        'ID',
        'Name',
        'Email',
        'Mobile',
        'AccountType',
        'TradeLicese',
        'ViewUser',
    ];
    usersData: Array<User> = [];
   constructor(
        private users: UsersService,
        private snackBar: MatSnackBar,
    ) { }

    ngOnInit() {
        this.getPartners();
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    /**
     * getPartners calls the users service to get the not approved users.
     */
    getPartners() {
        this.loadingState = true;
        this.users
            .getNotApprovedPartners()
            .subscribe(users => {
                console.log('USERS: ' , users);
                this.loadingState = false;
                this.usersData = users;
                users.forEach((user) => {
                    console.log(user.partnerProfile);
                });
            }, err => {
                console.error(err);
                this.loadingState = false;
                this.openSnackBar('Error Getting Users');
            });
    }

    /**
     * approvePartner approves the given partner
     */
    public approvePartner(partnerID: number) {
        this.loadingState = true;
        this.users
        .approvePartner(partnerID)
            .subscribe((partner: User) => {
                this.loadingState = false;
                console.log(partner);
                this.openSnackBar('Partner approved successfully');
                this.getPartners();
            }, err => {
                console.error(err);
                this.loadingState = false;
                this.openSnackBar('Error approving partner, please try again');
            });
    }

    public rejectPartner(partnerID: number) {
        this.loadingState = true;
        this.users
        .rejectPartner(partnerID)
        .subscribe((partner) => {
            this.loadingState = false;
            console.log(partner);
            this.openSnackBar('Partner rejected');
            this.getPartners();
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error rejecting partner, please try again');
        });
    }
}
