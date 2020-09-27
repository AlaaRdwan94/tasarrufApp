import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import {OffersService} from 'src/app/services/offers.service';
import {UsersService} from 'src/app/services/users.service';


export interface Link {
  position: number;
  title: string;
  description: string;
  value: string;
}

const LINKS: Link[] = [
  {
    position: 1,
    title: 'Twillio',
    description: 'The service responsible for sending SMS messages to users',
    value: 'https://www.twilio.com/console/usage'
  },
  {
    position: 2,
    title: 'Iyzipay',
    description: 'The payment gateway',
    value: 'https://sandbox-merchant.iyzipay.com/auth/'
  },
  {
    position: 3,
    title: 'Amazon AWS',
    description: 'S3 bucket holding static assets ,uploaded images and trade licenses',
    value: 'https://s3.console.aws.amazon.com/s3/buckets/tasarruf-file-repository/?region=us-east-1&tab=overview#'
  },
  {
    position: 4,
    title: 'Amazon AWS',
    description: 'EC2 instances manager',
    value: 'https://console.aws.amazon.com/ec2/v2/home?region=us-east-1#Instances:sort=instanceId'
  }
];

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  usersCount = 0; // the count of users in lendish
  partnersCount = 0; // the count of circles in lendish
  offersCount = 0; // the count of transactions in lendish
  pageTitle = 'Home Page';
  loadingState = false;
  cardTitles: Array<string> = ['Customers', 'Merchants', 'Offers'];
  cardSubtitles: Array<string> = [
    'The count of customers on the Tasarruf platform',
    'The count of merchants on the Tasarruf platform',
    'The count of consumed offers (successful scans)'
  ];
  displayedColumns: string[] = ['position', 'title', 'description', 'value'];
  dataSource = LINKS;

  constructor(
    private users: UsersService,
    private offers: OffersService,
    private snackBar: MatSnackBar
  ) { }
  openSnackBar(message: string) {
    this.snackBar.open(message, 'cancel', { duration: 3000 });
  }
  ngOnInit() {
    this.getCustomersCount();
    this.getPartnersCount();
    this.getOffersCount();
  }

  /**
   * getCustomersCount gets the count of users from the users service.
   */
  getCustomersCount() {
    this.loadingState = true;
    this.users
      .getCustomersCount()
      .subscribe((count: number) => {
        this.usersCount = count;
        this.loadingState = false;
      }, err => {
        this.openSnackBar('Error Getting Users Count');
        console.error(err);
        this.loadingState = false;
      });
  }

  /**
   * getCustomersCount gets the count of users from the users service.
   */
  getPartnersCount() {
    this.loadingState = true;
    this.users
      .getPartnersCount()
      .subscribe((count: number) => {
        this.partnersCount = count;
        this.loadingState = false;
      }, err => {
        this.openSnackBar('Error Getting Users Count');
        console.error(err);
        this.loadingState = false;
      });
  }


  /**
   * getTransactionsCount gets the count of transactions made on lendish.
   */
  getOffersCount() {
    this.loadingState = true;
    this.offers
      .getOffersCount()
      .subscribe((count: number) => {
        this.offersCount = count;
        setTimeout(() => {
          this.loadingState = false;
        }, 3000);
      }, err => {
          console.error(err);
          this.openSnackBar('Error Getting Transactions Count');
          this.loadingState = false;
      });
  }
}
