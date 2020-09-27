import { Component, OnInit } from '@angular/core';
import { FormControl, FormBuilder, FormGroup } from '@angular/forms';
import { UsersService } from 'src/app/services/users.service';
import { MatSnackBar, PageEvent } from '@angular/material';
import {debounceTime} from 'rxjs/operators';

import { User } from 'src/app/models/user';



@Component({
  selector: 'app-partners',
  templateUrl: './partners.component.html',
  styleUrls: ['./partners.component.scss']
})
export class PartnersComponent implements OnInit {
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
        'SetIsSharable',
        'SetExclusive',
        'ViewUser',
        'DeleteUser'
    ];
    usersData: Array<User> = [];
    exclusivePartners: Array<User> = [];
    searchForm: FormGroup = this.formBuilder.group({
        emailSearch: new FormControl(),
    });
    constructor(
        private users: UsersService,
        private snackBar: MatSnackBar,
        private formBuilder: FormBuilder,
    ) { }
    get emailSearch() {
        return this.searchForm.get('emailSearch');
    }
    ngOnInit() {
        this.getPartnersCount();
        this.getPartners();
        this.getExclusivePartners();
        this.emailSearch.valueChanges
        .pipe(
            debounceTime(1000),
        ).subscribe((x: string) => {
            this.searchUsers(x);
        });

    }
    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    paginatorCount(count: number): number {
        return Math.ceil(count / 10);
    }

    isPartnerExclusive(id: number): boolean {
        let isExclusive = false;
        this.exclusivePartners.forEach(user => {
            if (user.ID === id) {
                isExclusive = true;
            }
        });
        return isExclusive;
    }
    searchUsers(x: string) {
        console.log('SEARCHING FOR ' , x);
        this.users.searchUsers(x)
        .subscribe(users => {
            const partners: User[] = [];
            if (users) {
                users.forEach(user => {
                    if (user.accountType === 'partner') {
                        partners.push(user);
                    }
                });
            }
            this.usersData = partners;
        });
    }


    removePartnerExclusive(id: number) {
        this.loadingState = true;
        this.users.removePartnerExclusive(id)
        .subscribe(() => {
            this.getPartners();
            this.getExclusivePartners();
            this.loadingState = false;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error removing partner from exclusive partners, please try again');
        });
    }

    setPartnerExclusive(id: number) {
        this.loadingState = true;
        this.users.setPartnerExclusive(id)
        .subscribe(() => {
            this.getPartners();
            this.getExclusivePartners();
            this.loadingState = false;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error setting partner from exclusive partners, please try again');
        });
    }

    togglePartnerIsSharable(id: number) {
        this.loadingState = true;
        this.users.togglePartnerIsSharable(id)
        .subscribe(() => {
            this.getPartners();
            this.getExclusivePartners();
            this.loadingState = false;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error setting partner is sharable property, please try again');
        });
    }

    getNewUsers(event: PageEvent) {
        console.log(event);
    }

    /**
     * getUsersCount calls the users service to get the count of users.
     */
    getPartnersCount() {
        this.loadingState = true;
        this.users
            .getCustomersCount()
            .subscribe(count => {
                this.loadingState = false;
                this.usersCount = count;
                this.paginatorPagesCount = this.paginatorCount(count);
            }, err => {
                console.error(err);
                this.loadingState = false;
                this.openSnackBar('Error Getting Users Count');
            });
    }

    /**
     * getUsers returns the latest users with the given offset.
     * @param offset represents the number of items to skip.
     */
    getPartners() {
        this.loadingState = true;
        this.users
            .getPartners()
            .subscribe((returnedUsers: Array<User>) => {
                console.log(returnedUsers);
                this.loadingState = false;
                const partners: Array<User> = [];
                returnedUsers.forEach((user) => {
                    if (user.partnerProfile.approved ) {
                        partners.push(user);
                    }
                });
                this.usersData = partners;
            }, err => {
                console.error(err);
                this.openSnackBar('Error Loading Users');
            });
    }
    openCreateUserDialog() {
   }

   getExclusivePartners() {
       this.loadingState = true;
       this.users.getExclusivePartners()
       .subscribe((exclusiveUsers: Array<User>) => {
           console.log('EXCLUSIVE PARTNERS : ', exclusiveUsers);
           this.loadingState = false;
           this.exclusivePartners = exclusiveUsers;
       }, err => {
           console.error(err);
           this.openSnackBar('Error loading exclusive partners');
       });
   }

   public deletePartner(id: number) {
       this.loadingState = true;
       this.users.deleteUser(id)
       .subscribe((user: User) => {
           console.log('DELETED', user);
           this.loadingState = false;
           this.getPartners();
       }, err => {
           console.error(err);
           this.openSnackBar('Error deleting merchant, please try again');
       });
   }
}
