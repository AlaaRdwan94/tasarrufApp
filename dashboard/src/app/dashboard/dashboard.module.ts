import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login/login.component';
import { MainComponent } from './main/main.component';
import { DashboardRoutingModule } from './dashboard-routing.module';
import { MatCardModule} from '@angular/material';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatMenuModule } from '@angular/material/menu';
import { MatIconModule } from '@angular/material/icon';
import { LayoutModule } from '@angular/cdk/layout';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatListModule } from '@angular/material/list';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { HomeComponent } from './home/home.component';
import { HomeCardComponent } from './components/home-card/home-card.component';
import { SkeletonComponent } from './components/skeleton/skeleton.component';
import { MatTableModule } from '@angular/material/table';
import { UsersComponent } from './users/users.component';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatDialogModule } from '@angular/material/dialog';
import { MatTooltipModule} from '@angular/material/tooltip';
import {MatCheckboxModule} from '@angular/material/checkbox';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material';
import { MatSelectModule } from '@angular/material/select';
import {MatExpansionModule} from '@angular/material/expansion';
import {MatAutocompleteModule} from '@angular/material/autocomplete';
import { PartnersComponent } from './partners/partners.component';
import { OffersComponent } from './offers/offers.component';
import { NotApprovedComponent } from './not-approved/not-approved.component';
import { CategoriesComponent } from './categories/categories.component';
import { NewCategoryComponent } from './new-category/new-category.component';
import { PlansComponent } from './plans/plans.component';
import { NewPlanComponent } from './new-plan/new-plan.component';
import { UserComponent } from './components/user/user.component';
import { PartnerComponent } from './partner/partner.component';
import { CustomerComponent } from './customer/customer.component';
import { EditPlanComponent } from './edit-plan/edit-plan.component';
import { SupportInfoComponent } from './support-info/support-info.component';
import { ExclusiveOffersComponent } from './exclusive-offers/exclusive-offers.component';
import { EditCategoryComponent } from './edit-category/edit-category.component';
import { CititesComponent } from './citites/citites.component';
import { AddCityComponent } from './add-city/add-city.component';
import { EditCityComponent } from './edit-city/edit-city.component';
import { OffersTableComponent } from './components/offers-table/offers-table.component';
import { BranchesTableComponent } from './components/branches-table/branches-table.component';

@NgModule({
    declarations: [
        LoginComponent,
        MainComponent,
        HomeComponent,
        HomeCardComponent,
        SkeletonComponent,
        UsersComponent,
        PartnersComponent,
        OffersComponent,
        NotApprovedComponent,
        CategoriesComponent,
        NewCategoryComponent,
        PlansComponent,
        NewPlanComponent,
        UserComponent,
        PartnerComponent,
        CustomerComponent,
        EditPlanComponent,
        SupportInfoComponent,
        ExclusiveOffersComponent,
        EditCategoryComponent,
        CititesComponent,
        AddCityComponent,
        EditCityComponent,
        OffersTableComponent,
        BranchesTableComponent,
    ],
    entryComponents: [
    ],
    imports: [
        CommonModule,
        DashboardRoutingModule,
        FormsModule,
        ReactiveFormsModule,
        MatCardModule,
        MatFormFieldModule,
        MatInputModule,
        MatButtonModule,
        MatProgressSpinnerModule,
        MatSnackBarModule,
        MatGridListModule,
        MatMenuModule,
        MatIconModule,
        LayoutModule,
        MatToolbarModule,
        MatSidenavModule,
        MatListModule,
        MatProgressBarModule,
        MatTableModule,
        MatCheckboxModule,
        MatPaginatorModule,
        MatDialogModule,
        MatTooltipModule,
        MatDatepickerModule,
        MatNativeDateModule,
        MatExpansionModule,
        MatSelectModule,
        MatAutocompleteModule
    ]
})
export class DashboardModule { }
