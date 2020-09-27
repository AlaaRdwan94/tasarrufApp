import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FormGroup, FormControl } from '@angular/forms';
import {MatSnackBar} from '@angular/material';
import { Router} from '@angular/router';
import {CityService} from 'src/app/services/city.service';
import {City} from 'src/app/models/city';



@Component({
    selector: 'app-edit-city',
    templateUrl: './edit-city.component.html',
    styleUrls: ['./edit-city.component.scss']
})
export class EditCityComponent implements OnInit {
    public cityID: number;
    public loadingState = false;
    public toEditCity: City;
    public editingForm: FormGroup = new FormGroup({});
    constructor(
        private route: ActivatedRoute,
        private router: Router,
        private cityService: CityService,
        private snackBar: MatSnackBar
    ) { }

    ngOnInit() {
        this.cityID = parseInt(this.route.snapshot.paramMap.get('id'), 10);
        this.getCity(this.cityID);
    }

    /**
     * getCity returns the city with the given ID
     */
    public getCity(id: number) {
        this.loadingState = true;
        this.cityService.getCityByID(id).subscribe((city) => {
            this.toEditCity = city;
            console.log(city);
            this.setFormGroup();
            this.loadingState = false;
        });
    }

    public setFormGroup() {
        this.editingForm = new FormGroup({
            englishName: new FormControl(this.toEditCity.englishName),
            turkishName: new FormControl(this.toEditCity.turkishName)
        });
    }

    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public updateCity() {
        this.loadingState = true;
        const city: City =  {
            englishName : this.editingForm.value.englishName,
            turkishName: this.editingForm.value.turkishName
        };
        this.cityService.editCity(this.cityID, city).subscribe((updated) => {
            console.log('UPDATED CITY = ', updated);
            this.loadingState = false;
            this.openSnackBar('City updated successfully');
            this.router.navigate(['/admin/cities']);
        } , err => {
            this.loadingState = false;
            this.openSnackBar('Error updating city, please try again');
            console.error('ERROR UPDATING CITY : ', err);
        });
    }
}
