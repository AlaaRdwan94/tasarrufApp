import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material';
import {City} from 'src/app/models/city';
import {CityService} from 'src/app/services/city.service';

@Component({
    selector: 'app-citites',
    templateUrl: './citites.component.html',
    styleUrls: ['./citites.component.scss']
})
export class CititesComponent implements OnInit {

    constructor(private city: CityService, private snackBar: MatSnackBar) { }
    loadingState = false;
    categoriesCount = 0;
    paginatorPagesCount = 0;
    paginatorPageIndex = 0;
    displayedColumns: string[] = [
        'ID',
        'englishName',
        'turkishName',
        'edit',
        'delete',
    ];

    public cities: Array<City> = [];

    ngOnInit() {
        this.getCities();
    }

    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public getCities() {
        this.loadingState = true;
        this.city.getCities().subscribe((cities) => {
            this.cities = cities.sort((a, b) => {
                const aname = a.englishName.toLowerCase();
                const bname = b.englishName.toLowerCase();
                if (aname < bname) { return -1; }
                if (aname > bname) { return 1; }
                return 0;
            });
            this.loadingState = false;
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar(err);
        });
    }

    public deleteCity(id: number) {
        this.loadingState = true;
        this.city.deleteCity(id).subscribe((city) => {
            this.loadingState = false;
            console.log('deleted city', city);
            this.openSnackBar('City Deleted');
            this.getCities();
        }, err => {
            this.loadingState = false;
            console.error(err);
            this.openSnackBar('Deleting this city would cause data inconsistency');
        });
    }


}
