import { Component, OnInit, Input } from '@angular/core';
import {Offer} from 'src/app/models/offer';

@Component({
    selector: 'app-offers-table',
    templateUrl: './offers-table.component.html',
    styleUrls: ['./offers-table.component.scss']
})
export class OffersTableComponent implements OnInit {

    @Input() offersData: Offer[];
    displayedColumns: string[] = [
        'ID',
        'CreatedAt',
        'CustomerID',
        'PartnerID',
        'Amount',
        'Discount',
        'Total',
    ];


    constructor() { }

    ngOnInit() {
    }

}
