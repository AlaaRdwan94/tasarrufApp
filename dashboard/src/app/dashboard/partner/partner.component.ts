import { Component, OnInit } from '@angular/core';
import { User} from 'src/app/models/user';
import { MatSnackBar} from '@angular/material';
import { ActivatedRoute} from '@angular/router';
import { UsersService} from 'src/app/services/users.service';
import { Branch } from 'src/app/models/branch';
import {CategoriesService} from 'src/app/services/categories.service';
import {Category} from 'src/app/models/category';



@Component({
    selector: 'app-partner',
    templateUrl: './partner.component.html',
    styleUrls: ['./partner.component.scss']
})
export class PartnerComponent implements OnInit {

    public partnerID: number;
    public partner: User;
    public loadingState = false;
    public branches: Array<Branch>;
    public category: Category;

    constructor(
        private route: ActivatedRoute,
        private userService: UsersService,
        private snackBar: MatSnackBar,
        private categoriesService: CategoriesService,
    ) { }

    ngOnInit() {
        this.partnerID = parseInt(this.route.snapshot.paramMap.get('id'), 10);
        this.getPartner(this.partnerID);
        this.getBranches(this.partnerID);
    }

    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cance', { duration: 3000});
    }

    /**
     * getPartner gets the partner with the given ID
     */
    public getPartner(partnerID: number) {
        this.loadingState = true;
        this.userService.getPartner(partnerID).subscribe((partner) => {
            this.partner = partner;
            this.loadingState = false;
            console.log('Partner : ', this.partner);
            this.categoriesService.getCategories().subscribe((categories) => {
                console.log('CATEGORIES =', categories);
                categories.forEach((category) => {
                    if (category.ID === this.partner.partnerProfile.categroryID) {
                        this.category = category;
                    }
                });
            });
        });
    }
    public toggleActiveProperty(userID: number) {
        this.userService.toggleActiveProperty(userID)
        .subscribe((partner: User) => {
            this.partner = partner;
        }, err => {
            this.openSnackBar(err);
        });
    }

    public getBranches(userID: number) {
        this.loadingState = true;
        this.userService.getBranchesOfPartner(userID)
        .subscribe((branches: Array<Branch>) => {
            this.branches = branches;
            this.loadingState = false;
            console.log('BRANHCES : ', branches);
        }, err => {
            this.openSnackBar(err);
            this.loadingState = false;
        });
    }
    /**
     * approvePartner approves the given partner
     */
    public approvePartner(partnerID: number) {
        this.loadingState = true;
        this.userService
        .approvePartner(partnerID)
            .subscribe((partner: User) => {
                this.loadingState = false;
                console.log(partner);
                this.partner.partnerProfile.approved = true;
                this.openSnackBar('Partner approved successfully');
            }, err => {
                console.error(err);
                this.loadingState = false;
                this.openSnackBar('Error approving partner, please try again');
            });
    }

    public rejectPartner(partnerID: number) {
        this.loadingState = true;
        this.userService
        .rejectPartner(partnerID)
        .subscribe((partner) => {
            this.loadingState = false;
            console.log(partner);
            this.openSnackBar('Partner rejected');
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error rejecting partner, please try again');
        });
    }


}
