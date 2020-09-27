import { Component, OnInit } from '@angular/core';
import {
    FormControl,
    FormGroupDirective,
    NgForm,
    Validators,
    FormGroup
} from '@angular/forms';
import { ErrorStateMatcher } from '@angular/material/core';
import {MatSnackBar} from '@angular/material';
import { Router} from '@angular/router';
import {CityService} from 'src/app/services/city.service';
import {City} from 'src/app/models/city';

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
    selector: 'app-add-city',
    templateUrl: './add-city.component.html',
    styleUrls: ['./add-city.component.scss']
})
export class AddCityComponent implements OnInit {

    constructor(
        private router: Router,
        private cityService: CityService,
        private snackBar: MatSnackBar
    ) { }

    newCity: City = {
        englishName: '',
        turkishName: '',
    };

    createCityForm: FormGroup = new FormGroup({});
    public loadingState = false;
    matcher = new MyErrorStateMatcher();

    ngOnInit() {
        this.createCityForm = new FormGroup({
            englishName: new FormControl(this.newCity.englishName, [
                Validators.required,
            ]),
            turkishName: new FormControl(this.newCity.turkishName, [
                Validators.required,
            ]),
        });

    }
    get turkishName() {
        return this.createCityForm.get('turkishName');
    }

    get englishName() {
        return this.createCityForm.get('englishName');
    }
    openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }

    public onSubmitClick(): void {
        this.onSubmitCity(this.createCityForm.value);
    }

    public onSubmitCity(city: City) {
        console.log(city);
        this.loadingState = true;
        this.cityService.createCity(city).subscribe((p) => {
            console.log(p);
            this.openSnackBar('City created successfully');
            this.loadingState = false;
            this.router.navigate(['/admin/cities']);
        }, err => {
            console.error(err);
            this.loadingState = false;
            this.openSnackBar('Error : ' + err.error.message);
        });
    }


}
