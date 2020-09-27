import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { AuthGuard } from '../guards/auth.guard';
import { HomeComponent } from './home/home.component';
import { UsersComponent } from './users/users.component';
import {PartnersComponent} from 'src/app/dashboard/partners/partners.component';
import { OffersComponent } from './offers/offers.component';
import {NotApprovedComponent} from 'src/app/dashboard/not-approved/not-approved.component';
import { CategoriesComponent } from './categories/categories.component';
import { NewCategoryComponent } from './new-category/new-category.component';
import { PlansComponent } from './plans/plans.component';
import {NewPlanComponent} from 'src/app/dashboard/new-plan/new-plan.component';
import {EditPlanComponent} from 'src/app/dashboard/edit-plan/edit-plan.component';
import {EditCategoryComponent} from 'src/app/dashboard/edit-category/edit-category.component';
import {ExclusiveOffersComponent} from 'src/app/dashboard/exclusive-offers/exclusive-offers.component';
import {PartnerComponent} from 'src/app/dashboard/partner/partner.component';
import {CustomerComponent} from 'src/app/dashboard/customer/customer.component';
import {CititesComponent} from 'src/app/dashboard/citites/citites.component';
import {AddCityComponent} from 'src/app/dashboard/add-city/add-city.component';
import {EditCityComponent} from 'src/app/dashboard/edit-city/edit-city.component';
import {SupportInfoComponent} from 'src/app/dashboard/support-info/support-info.component';

const routes: Routes = [
    {
        path: '',
        component: HomeComponent,
        canActivate: [AuthGuard]
    },
    {
        path: 'login',
        component: LoginComponent
    },
    {
        path: 'customers',
        component: UsersComponent
    },
    {
        path: 'partners',
        component: PartnersComponent
    },
    {
        path: 'offers',
        component: OffersComponent
    },
    {
        path: 'approve-partner',
        component: NotApprovedComponent
    },
    {
        path: 'categories',
        component: CategoriesComponent
    },
    {
        path: 'categories/new-category',
        component: NewCategoryComponent
    },
    {
        path: 'categories/edit-category/:id',
        component: EditCategoryComponent
    },
    {
        path: 'plans',
        component: PlansComponent
    },
    {
        path: 'plans/new-plan',
        component: NewPlanComponent
    },
    {
        path: 'plans/edit-plan/:id',
        component: EditPlanComponent
    },
    {
        path: 'exclusive',
        component: ExclusiveOffersComponent
    },
    {
        path: 'partner/:id',
        component: PartnerComponent
    },
    {
        path: 'customer/:id',
        component: CustomerComponent
    },
    {
        path: 'cities',
        component: CititesComponent
    },
    {
        path: 'cities/add-city',
        component: AddCityComponent
    },
    {
        path: 'cities/edit-city/:id',
        component: EditCityComponent
    },
    {
        path: 'support',
        component: SupportInfoComponent
    },
    {
        path: 'exclusive-offers',
        component: ExclusiveOffersComponent
    }
];

@NgModule({
    imports: [RouterModule.forChild(routes)],
    exports: [RouterModule]
})
export class DashboardRoutingModule {}
