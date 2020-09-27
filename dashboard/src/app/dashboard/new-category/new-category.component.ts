import { Component, OnInit } from '@angular/core';
import {
    FormControl,
    FormGroupDirective,
    NgForm,
    Validators,
    FormGroup
} from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import {CategoriesService} from 'src/app/services/categories.service';
import {MatSnackBar} from '@angular/material';
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
export class NewCategory {
    englishName: string;
    turkishName: string;
}

@Component({
    selector: 'app-new-category',
    templateUrl: './new-category.component.html',
    styleUrls: ['./new-category.component.scss']
})

export class NewCategoryComponent implements OnInit {

    constructor(private router: Router, private category: CategoriesService, private snackBar: MatSnackBar) { }

    newCategory: NewCategory = {
        englishName: '',
        turkishName: '',
    };
    public loadingState = false;

    createCategoryForm: FormGroup = new FormGroup({});
    matcher = new MyErrorStateMatcher();

    get turkishName() {
        return this.createCategoryForm.get('turkishName');
    }

    get englishName() {
        return this.createCategoryForm.get('englishName');
    }

    ngOnInit() {
        this.createCategoryForm = new FormGroup({
            englishName: new FormControl(this.newCategory.englishName, [
                Validators.required,
            ]),
            turkishName: new FormControl(this.newCategory.turkishName, [
                Validators.required,
            ])
        });
    }

    public onSubmitClick(): void {
        this.submitCategory(this.createCategoryForm.value);
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    private submitCategory(category: NewCategory) {
        console.log(category);
        this.loadingState = true;
        this.category.createCategory(category.englishName, category.turkishName).subscribe((c) => {
            console.log(c);
            this.openSnackBar('category created successfully');
            this.loadingState = false;
            this.router.navigate(['/admin/categories']);
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error creating category , please try again');
        });
    }
}
