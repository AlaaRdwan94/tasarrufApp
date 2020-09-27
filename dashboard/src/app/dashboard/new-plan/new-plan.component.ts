import { Component, OnInit } from '@angular/core';
import {
    FormControl,
    FormGroupDirective,
    NgForm,
    Validators,
    FormGroup
} from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import {PlanService} from 'src/app/services/plan.service';
import {MatSnackBar} from '@angular/material';
import {Plan} from 'src/app/models/plan';
import { Router} from '@angular/router';

export class MyErrorStateMatcher implements ErrorStateMatcher {
    isErrorState(
        control: FormControl | null,
        form: FormGroupDirective | NgForm | null
    ): boolean {
        const isSubmitted = form && form.submitted;
        return !!(
            control &&
                control.invalid &&
                (control.dirty || control.touched || isSubmitted)
        );
    }
}

@Component({
    selector: 'app-new-plan',
    templateUrl: './new-plan.component.html',
    styleUrls: ['./new-plan.component.scss']
})

export class NewPlanComponent implements OnInit {
    public planIcons: Array<string> = [
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/gold.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/silver.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/bronze.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/platinum.png',
        'https://tasarruf-file-repository.s3.amazonaws.com/plan_images/Vector+Smart+Object2.png',
    ];

    constructor(
        private router: Router,
        private plan: PlanService,
        private snackBar: MatSnackBar
    ) { }

    newPlan: Plan = {
        englishName: '',
        engishDescription: '',
        trukishName: '',
        turkishDescription: '',
        price: null,
        countOfOffers: null,
        image: '',
        isDefault: false,
    };

    createPlanForm: FormGroup = new FormGroup({});
    public loadingState = false;
    matcher = new MyErrorStateMatcher();

    get trukishName() {
        return this.createPlanForm.get('trukishName');
    }

    get englishName() {
        return this.createPlanForm.get('englishName');
    }

    get engishDescription() {
        return this.createPlanForm.get('engishDescription');
    }

    get turkishDescription() {
        return this.createPlanForm.get('turkishDescription');
    }

    get price() {
        return this.createPlanForm.get('price');
    }

    get countOfOffers() {
        return this.createPlanForm.get('countOfOffers');
    }

    get image() {
        return this.createPlanForm.get('image');
    }

    get isDefault() {
        return this.createPlanForm.get('isDefault');
    }

    ngOnInit() {
        this.createPlanForm = new FormGroup({
            englishName: new FormControl(this.newPlan.englishName, [
                Validators.required,
            ]),
            trukishName: new FormControl(this.newPlan.trukishName, [
                Validators.required,
            ]),
            engishDescription: new FormControl(this.newPlan.engishDescription, [
                Validators.required,
            ]),
            turkishDescription: new FormControl(this.newPlan.turkishDescription, [
                Validators.required,
            ]),
            price: new FormControl(this.newPlan.price, [
            ]),
            countOfOffers: new FormControl(this.newPlan.countOfOffers, [
            ]),
            image : new FormControl(this.newPlan.image, [
                Validators.required,
            ]),
            isDefault: new FormControl(this.newPlan.isDefault)
        });
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public onSubmitClick(): void {
        this.onSubmitPlan(this.createPlanForm.value);
    }

    public onSubmitPlan(plan: Plan) {
        console.log(plan);
        this.loadingState = true;
        this.plan.createPlan(plan).subscribe((p) => {
            console.log(p);
            this.openSnackBar('Plan created successfully');
            this.loadingState = false;
            this.router.navigate(['/admin/plans']);
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error : ' + err.error.message);
        });
    }
}
