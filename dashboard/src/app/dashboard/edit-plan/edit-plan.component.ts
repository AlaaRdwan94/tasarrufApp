import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import {Plan} from 'src/app/models/plan';
import {PlanService} from 'src/app/services/plan.service';
import { FormGroup, FormControl } from '@angular/forms';
import {MatSnackBar} from '@angular/material';
import { Category } from 'src/app/models/category';
import {CategoriesService} from 'src/app/services/categories.service';

@Component({
    selector: 'app-edit-plan',
    templateUrl: './edit-plan.component.html',
    styleUrls: ['./edit-plan.component.scss']
})
export class EditPlanComponent implements OnInit {
    public planID: number;
    public toEditPlan: Plan;
    public loadingState = false;
    public planIcons: Array<string> = [
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/gold.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/silver.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/bronze.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/platinum.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/Vector+Smart+Object2.png',
    ];
    public currentCategories: Array<Category> = [];

    public notAssociatedCategories: Array<Category>;

    editingPlan = new FormGroup({
        englishName: new FormControl(''),
        engishDescription: new FormControl(''),
        trukishName: new FormControl(''),
        turkishDescription: new FormControl(''),
        price: new FormControl(''),
        icon: new FormControl(''),
        newCategoryID: new FormControl(0),
        isDefault: new FormControl(false),
    });
    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private planService: PlanService,
        private categoriesService: CategoriesService,
        private snackBar: MatSnackBar
    ) { }

    ngOnInit() {
        this.planID = parseInt(this.route.snapshot.paramMap.get('id'), 10);
        this.getPlans();
        this.getCurrentCategories();
        this.getNotAssociatedCategories();
    }

    /**
     * getPlans gets the all plans
     */
    public getPlans() {
        this.loadingState = true;
        this.planService.getPlans().subscribe((plans) => {
            plans.forEach((plan) => {
                if (plan.ID === this.planID) {
                    this.toEditPlan = plan;
                    console.log('PLAN = ', plan );
                }
            });
            this.editingPlan = new FormGroup({
                englishName: new FormControl(this.toEditPlan.englishName),
                turkishName: new FormControl(this.toEditPlan.trukishName),
                englishDescription: new FormControl(this.toEditPlan.engishDescription),
                turkishDescription: new FormControl(this.toEditPlan.turkishDescription),
                countOfOffers: new FormControl(this.toEditPlan.countOfOffers),
                price: new FormControl(this.toEditPlan.price),
                icon : new FormControl(this.toEditPlan.image),
                newCategoryID : new FormControl([]),
                isDefault: new FormControl(this.toEditPlan.isDefault)
            });
            this.loadingState = false;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('error getting plans, please try again');
        });
    }

    public getNotAssociatedCategories() {
        this.loadingState = true;
        this.categoriesService.getCategories()
            .subscribe((categories) => {
                this.loadingState = false;
                this.notAssociatedCategories = categories.filter((c) => {
                    let isNotAssociated = true;
                    this.currentCategories.forEach((x) => {
                        if (x.ID === c.ID) {
                            isNotAssociated = false;
                        }
                    });
                    return isNotAssociated;
                });
            });
    }

    public getCurrentCategories() {
        this.loadingState = true;
        this.planService.getCategoriesOfPlan(this.planID)
            .subscribe((categories) => {
                console.log(categories);
                this.loadingState = false;
                this.currentCategories = categories;
            }, err => {
                this.loadingState = false;
                console.error(err);
                this.openSnackBar('error getting categories, please try again');
            });
    }

    public associatePlanWithCategory(categoryID: number) {
        this.loadingState = true;
        this.planService.associatePlanWithCategory(this.planID, categoryID)
            .subscribe(() => {
                this.loadingState = false;
                this.openSnackBar('Associated Succesfully');
            }, err => {
                this.loadingState = false;
                console.error(err);
            });
    }

    public removePlanCategoryAssociation(categoryID: number) {
        this.loadingState = true;
        this.planService.removePlanCategoryAssociation(this.planID, categoryID)
            .subscribe(() => {
                this.loadingState = false;
                this.openSnackBar('Association removed successfully');
                this.currentCategories = this.currentCategories.filter((c) => {
                    return c.ID !== categoryID;
                });
                this.getNotAssociatedCategories();
            }, err => {
                this.loadingState = false;
                console.error(err);
                this.openSnackBar('error removing association, please try again');
            });
    }

    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public updatePlan() {
        this.loadingState = true;
        const plan: Plan = {
            ID: this.toEditPlan.ID,
            englishName: this.editingPlan.value.englishName,
            trukishName: this.editingPlan.value.turkishName,
            engishDescription: this.editingPlan.value.englishDescription,
            turkishDescription: this.editingPlan.value.turkishDescription,
            price: this.editingPlan.value.price,
            countOfOffers: this.editingPlan.value.countOfOffers,
            image: this.editingPlan.value.icon,
            isDefault: this.editingPlan.value.isDefault,
        };
        this.editingPlan.value.newCategoryID.forEach((id: number) => {
            this.planService.associatePlanWithCategory(this.planID, id).subscribe((success) => {
                console.log(success);
            }, err => {
                this.openSnackBar('There has been an error processing your request, please try again');
                console.error(err);
                return;
            });
        });
        this.planService.updatePlan(plan).subscribe((newPlan) => {
            console.log('PLAN = ', newPlan);
            this.openSnackBar('Plan Editied successfully');
            this.router.navigate(['/admin/plans']);
        });
    }
}
