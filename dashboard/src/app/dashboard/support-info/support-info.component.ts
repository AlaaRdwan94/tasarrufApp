import { Component, OnInit } from '@angular/core';
import {Support} from 'src/app/models/support';
import { FormGroup, FormControl } from '@angular/forms';
import {MatSnackBar} from '@angular/material';
import {SupportService} from 'src/app/services/support.service';

@Component({
    selector: 'app-support-info',
    templateUrl: './support-info.component.html',
    styleUrls: ['./support-info.component.scss']
})
export class SupportInfoComponent implements OnInit {

    public supportInfo: Support;
    public loadingState = false;
    constructor(private support: SupportService, private snackBar: MatSnackBar) { }
    editingSupportInfo = new FormGroup({
        mobile: new FormControl(''),
        email: new FormControl(''),
    });
    ngOnInit() {
        this.getSupportInfo();
    }

    /**
     * getSupportInfo gets the current support info
     */
    public getSupportInfo() {
        this.loadingState = true;
        this.support.getSupportInfo().subscribe((info) => {
            console.log('INFO', info);
            this.supportInfo = info;
            this.loadingState = false;
            this.editingSupportInfo = new FormGroup({
                mobile: new FormControl(this.supportInfo.mobile),
                email: new FormControl(this.supportInfo.email),
            });
        }, err => {
            this.loadingState = false;
            console.error(err);
            this.openSnackBar('Error getting support info , please try again');
        });

    }

    public updateSupportInfo() {
        this.loadingState = true;
        const supportInfo: Support = {
            ID: this.supportInfo.ID,
            CreatedAt: this.supportInfo.CreatedAt,
            UpdatedAt: this.supportInfo.UpdatedAt,
            DeletedAt: this.supportInfo.DeletedAt,
            mobile: this.editingSupportInfo.value.mobile,
            email: this.editingSupportInfo.value.email,
        };
        this.support.updateSupportInfo(supportInfo.email, supportInfo.mobile).subscribe((info) => {
            this.loadingState = false;
            console.log('UPDATED', info);
            this.getSupportInfo();
        }, err => {
            this.loadingState = false;
            console.error(err);
            this.openSnackBar('Ivalid email address or phone number');
        });
    }

    public openSnackBar(message: string) {
        this.snackBar.open(message, 'cancel', { duration: 3000 });
    }
}
