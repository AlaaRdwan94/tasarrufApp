import { Component, OnInit } from '@angular/core';
import { MatSnackBar, PageEvent } from '@angular/material';
import { Offer } from 'src/app/models/offer';
import {OffersService} from 'src/app/services/offers.service';
import {FormGroup, FormBuilder, FormControl} from '@angular/forms';



@Component({
    selector: 'app-offers',
    templateUrl: './offers.component.html',
    styleUrls: ['./offers.component.scss']
})
export class OffersComponent implements OnInit {

    loadingState = false;
    offersCount = 0;
    paginatorPagesCount = 0;
    paginatorPageIndex = 0;
    displayedColumns: string[] = [
        'ID',
        'CreatedAt',
        'CustomerID',
        'PartnerID',
        'Amount',
        'Discount',
        'Total',
    ];
    offersData: Array<Offer> = [];
    filterForm: FormGroup = this.formBuilder.group({
        customerID: new FormControl(),
        partnerID: new FormControl(),
    });
    constructor(
        private offers: OffersService,
        private snackBar: MatSnackBar,
        private formBuilder: FormBuilder,
    ) { }

    get customerIDFilter() {
        return this.filterForm.get('customerID');
    }

    get partnerIDFilter() {
        return this.filterForm.get('partnerID');
    }

    ngOnInit() {
        this.getOffersCount();
        this.getOffers();
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
    getOffersCount() {
        this.loadingState = true;
        this.offers
        .getOffersCount()
        .subscribe(count => {
            this.loadingState = false;
            this.offersCount = count;
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
    getOffers() {
        this.loadingState = true;
        this.offers
        .getOffers()
        .subscribe((returnedOffers: Array<Offer>) => {
            console.log(returnedOffers);
            this.loadingState = false;
            this.offersData = returnedOffers;
        }, err => {
            console.error(err);
            this.openSnackBar('Error Loading Offers');
        });
    }

    filterOffers() {
        console.log('FILTER OFFERS');
        const filters = this.filterForm.value;
        if (filters.customerID > 0) {
        this.offersData = this.offersData.filter((offer) => {
            return offer.customerID === filters.customerID;
        });
        }
        if (filters.partnerID > 0) {
        this.offersData = this.offersData.filter((offer) => {
            return offer.partnerID === filters.partnerID;
        });
        }
        if (filters.customerID === null && filters.partnerID === null) {
            this.getOffers();
        }
    }
}
