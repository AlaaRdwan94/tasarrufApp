import { Component, OnInit, Input } from '@angular/core';
import {Branch} from 'src/app/models/branch';

@Component({
    selector: 'app-branches-table',
    templateUrl: './branches-table.component.html',
    styleUrls: ['./branches-table.component.scss']
})
export class BranchesTableComponent implements OnInit {

    @Input() branchesData: Array<Branch>;
    displayedColumns: string[] = [
        'ID',
        'cityName',
        'phone',
        'mobile',
        'address',
    ];
    constructor() { }

    ngOnInit() {
    }

}
