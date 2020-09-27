import { Component, OnInit, Input } from '@angular/core';
import { MatSnackBar} from '@angular/material';
import { User } from 'src/app/models/user';
import {PlanService} from 'src/app/services/plan.service';
import {Plan} from 'src/app/models/plan';
import {Offer} from 'src/app/models/offer';
import { UsersService } from 'src/app/services/users.service';

@Component({
    selector: 'app-user',
    templateUrl: './user.component.html',
    styleUrls: ['./user.component.scss']
})
export class UserComponent implements OnInit {

    @Input() customer: User;
    currentPlan: Plan;
    allPlans: Array<Plan>;
    changedPlanID: number;
    offers: Array<Offer>;

    constructor(private plansService: PlanService, private snackBar: MatSnackBar, private userService: UsersService) { }

    ngOnInit() {
        this.getCurrentPlan();
    }

    getCurrentPlan() {
        this.plansService.getPlans()
        .subscribe((plans) => {
            this.allPlans = plans;
            plans.forEach(p => {
                if (p.ID === this.customer.Subscription.planID) {
                    this.currentPlan = p;
                    this.changedPlanID = p.ID;
                }
            });
        }, err => {
            console.error(err);
        });
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    upgradePlan() {
       this.plansService.adminUpgradePlan(this.changedPlanID, this.customer.ID)
       .subscribe((newPlan: Plan) => {
           console.log('NEW PLAN', newPlan);
           this.openSnackBar('Subscription Upgraded !');
           this.currentPlan = newPlan;
       }, err => {
           console.error(err);
           this.openSnackBar('Error updating subscription, please try again');
       });
    }

    public toggleActiveProperty(userID: number) {
        this.userService.toggleActiveProperty(userID)
        .subscribe((customer: User) => {
            this.customer = customer;
        }, err => {
            this.openSnackBar(err);
        });
    }

}
