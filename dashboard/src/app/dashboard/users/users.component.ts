import { Component,  OnInit } from '@angular/core';
import { FormControl, FormBuilder, FormGroup } from '@angular/forms';
import { UsersService } from 'src/app/services/users.service';
import { MatSnackBar, PageEvent } from '@angular/material';

import { User } from 'src/app/models/user';
import {debounceTime} from 'rxjs/operators';

@Component({
    selector: 'app-users',
    templateUrl: './users.component.html',
    styleUrls: ['./users.component.scss']
})
export class UsersComponent implements OnInit {
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
        'ViewUser',
        'DeleteUser',
    ];
    usersData: Array<User> = [];
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
        this.getUsersCount();
        this.getCustomers();
        this.emailSearch.valueChanges
        .pipe(
            debounceTime(1000),
        ).subscribe((x: string) => {
            this.searchUsers(x);
        });
    }

    searchUsers(x: string) {
        console.log('SEARCHING FOR ' , x);
        this.users.searchUsers(x)
        .subscribe(users => {
            const customers: User[] = [];
            if (users) {
                users.forEach(user => {
                    if (user.accountType === 'user') {
                        customers.push(user);
                    }
                });
            }
            this.usersData = customers;
        });
    }
    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    paginatorCount(count: number): number {
        return Math.ceil(count / 10);
    }

    getNewUsers(event: PageEvent) {
        console.log(event);
    }

    /**
     * getUsersCount calls the users service to get the count of users.
     */
    getUsersCount() {
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
    getCustomers() {
        this.loadingState = true;
        this.users
        .getCustomers()
        .subscribe((returnedUsers: Array<User>) => {
            console.log(returnedUsers);
            this.loadingState = false;
            this.usersData = returnedUsers;
        }, err => {
            console.error(err);
            this.openSnackBar('Error Loading Users');
        });
    }

    public deleteCustomer(id: number) {
        this.loadingState = true;
        this.users
        .deleteUser(id)
        .subscribe((customer: User) => {
            console.log('DELETED', customer);
            this.loadingState = false;
            this.openSnackBar('Customer deleted successfully');
            this.getCustomers();
        }, err => {
            console.error(err);
            this.openSnackBar('Error deleting user, please try again');
        });
    }

    openCreateUserDialog() {
    }
    openDialog(userID: number, firstName: string, lastName: string): void {
        console.log(userID, firstName, lastName);
    }

}
