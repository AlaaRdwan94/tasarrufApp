import { Component, OnInit } from '@angular/core';
import { MatSnackBar} from '@angular/material';
import {Plan} from 'src/app/models/plan';
import {PlanService} from 'src/app/services/plan.service';

@Component({
    selector: 'app-plans',
    templateUrl: './plans.component.html',
    styleUrls: ['./plans.component.scss']
})
export class PlansComponent implements OnInit {

    constructor(private plan: PlanService, private snackBar: MatSnackBar) { }
    loadingState = false;
    categoriesCount = 0;
    paginatorPagesCount = 0;
    paginatorPageIndex = 0;
    displayedColumns: string[] = [
        'ID',
        'default',
        'englishName',
        'turkishName',
        'price',
        'countOfOffers',
        'edit',
        'delete',
    ];

    public plans: Array<Plan> = [];

    ngOnInit() {
        this.getPlans();
    }

    /**
     * getPlans gets all plans
     */
    public getPlans() {
        this.plan.getPlans()
        .subscribe((plan: Array<Plan>) => {
            this.plans = plan;
        });
    }

    /**
     * deletePlan deletes the plan with the given ID
     */
    public deletePlan(id: number) {
        this.plan.deletePlan(id)
        .subscribe((plan: Plan) => {
            console.log('DELETED', plan);
            this.openSnackBar('Plan deleted successfully');
            this.getPlans();
        }, err => {
            console.error(err);
            this.openSnackBar('Cannot delete default plan');
        });
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }
}
