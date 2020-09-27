import { Component, OnInit } from '@angular/core';
import { MatSnackBar} from '@angular/material';
import { ActivatedRoute } from '@angular/router';
import { UsersService } from 'src/app/services/users.service';
import { User } from 'src/app/models/user';
import { OffersService } from 'src/app/services/offers.service';
import {Offer} from 'src/app/models/offer';


@Component({
    selector: 'app-customer',
    templateUrl: './customer.component.html',
    styleUrls: ['./customer.component.scss']
})
export class CustomerComponent implements OnInit {
    public customer: User;
    public loadingState = false;
    public offers: Array<Offer>;
    constructor(
        private route: ActivatedRoute,
        private userService: UsersService,
        private offersService: OffersService,
        private snackBar: MatSnackBar
    ) { }


    ngOnInit() {
        const id: number = parseInt(this.route.snapshot.paramMap.get('id'), 10);
        this.getCustomer(id);
        this.getCustomerOffers(id);
    }

    /**
     * getCustomer returns the customer with the given ID
     */
    public getCustomer(id: number) {
        this.loadingState = true;
        this.userService.getCustomer(id).subscribe((customer) => {
            this.customer = customer;
            this.loadingState = false;
            console.log('CUSTOMER = ', this.customer);
        });
    }
    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }
    public getCustomerOffers(id: number) {
        this.loadingState = true;
        this.offersService.getOffersOfCustomer(id)
        .subscribe((offers) => {
            this.loadingState = false;
            this.offers = offers;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('error loading offers, please try again');
        });

    }
}
