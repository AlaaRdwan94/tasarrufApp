import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { FormGroup, FormControl } from '@angular/forms';
import {MatSnackBar} from '@angular/material';
import { Category } from 'src/app/models/category';
import {CategoriesService} from 'src/app/services/categories.service';



@Component({
    selector: 'app-edit-category',
    templateUrl: './edit-category.component.html',
    styleUrls: ['./edit-category.component.scss']
})
export class EditCategoryComponent implements OnInit {
    public categoryID: number;
    public toEditCategory: Category;
    public loadingState = false;
    editingCategory = new FormGroup({
        englishName: new FormControl(''),
        turkishName: new FormControl(''),
    });
    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private categoriesService: CategoriesService,
        private snackBar: MatSnackBar,
    ) {}

    ngOnInit() {
        this.categoryID = parseInt(this.route.snapshot.paramMap.get('id'), 10);
        this.getCategory(this.categoryID);
    }

    public getCategory(id: number) {
        this.loadingState = true;
        this.categoriesService.getCategories()
        .subscribe((categories: Array<Category>) => {
            categories.forEach((c) => {
                if (c.ID === id) {
                    this.toEditCategory = c;
                    this.loadingState = false;
                }
            });
            this.editingCategory = new FormGroup({
                englishName: new FormControl(this.toEditCategory.englishName),
                turkishName: new FormControl(this.toEditCategory.turkishName),
            });
        }, err => {
            console.error(err);
            this.openSnackBar('Error getting category, please try again');
        });
    }
    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public updateCategory() {
        this.categoriesService.editCategory(
            this.categoryID,
            this.editingCategory.value.englishName,
            this.editingCategory.value.turkishName
        )
        .subscribe((c) => {
            console.log('EDITED ' , c);
            this.openSnackBar('Updated successfully');
            this.router.navigate(['admin/categories']);
        }, err => {
            console.error(err);
            this.openSnackBar('Error editing category, please try again');
        });
    }

}
