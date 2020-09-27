import { Component, OnInit } from '@angular/core';
import { MatSnackBar} from '@angular/material';
import {CategoriesService} from 'src/app/services/categories.service';
import {Category} from 'src/app/models/category';

@Component({
    selector: 'app-categories',
    templateUrl: './categories.component.html',
    styleUrls: ['./categories.component.scss']
})
export class CategoriesComponent implements OnInit {

    constructor(
        private category: CategoriesService,
        private snackBar: MatSnackBar,
    ) { }
    loadingState = false;
    categoriesCount = 0;
    paginatorPagesCount = 0;
    paginatorPageIndex = 0;
    displayedColumns: string[] = [
        'ID',
        'CreatedAt',
        'englishName',
        'turkishName',
        'edit',
        'delete',
    ];
    public categories: Array<Category> = [];

    ngOnInit() {
        this.getCategories();
    }

    /*
     * getCategories gets all categories
     */
    public getCategories() {
        this.category.getCategories()
        .subscribe((categories: Array<Category>) => {
             this.categories = categories.sort((a, b) => {
                const aname = a.englishName.toLowerCase();
                const bname = b.englishName.toLowerCase();
                if (aname < bname) { return -1; }
                if (aname > bname) { return 1; }
                return 0;
            });
        });
    }

    /**
     * deleteCategory deletes the plan with the given id
     */
    public deleteCategory(id: number) {
        this.category.deleteCategoy(id)
        .subscribe((category: Category) => {
            console.log('DELETED', category);
            this.openSnackBar('Category deleted successfully');
        });
    }

    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    paginatorCount(count: number): number {
        return Math.ceil(count / 10);
    }

}
