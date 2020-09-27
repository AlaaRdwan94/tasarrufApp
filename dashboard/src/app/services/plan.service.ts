import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';
import {Plan} from '../models/plan';
import {Category} from '../models/category';

interface PlansResponse {
    plans: Array<Plan>;
}

interface CategoriesResponse {
    categories: Array<Category>;
}


interface PlanResponse {
    plan: Plan;
}

interface SuccessResponse {
    success: string;
    subscription: {
        ID: number;
        plan: Plan;
    };
}

@Injectable({
    providedIn: 'root'
})
export class PlanService {

    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }
    private rootUrl: string;

    /**
     * getPlans returns all the plans
     */
    public getPlans(): Observable< Array<Plan>> {
        return this.http
        .get<PlansResponse>(`${this.rootUrl}/plans`, {
            headers: { Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.plans)
        );
    }

    /**
     * deletePlan deletes the plan with the given id
     */
    public deletePlan(id: number): Observable<Plan> {
        return this.http
        .delete<PlanResponse>(`${this.rootUrl}/plans/${id}`, {
            headers: { Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.plan)
        );
    }

    /**
     * createPlan creates a new plan
     */
    public createPlan(plan: Plan): Observable<Plan> {
        return this.http
        .post<PlanResponse>(`${this.rootUrl}/plans`, {
            englishName: plan.englishName,
            turkishName: plan.trukishName,
            englishDescription: plan.engishDescription,
            turkishDescription: plan.turkishDescription,
            price: plan.price,
            countOfOffers: plan.countOfOffers,
            image: plan.image,
            isDefault: plan.isDefault,
        }, {
            headers : {
                Token : this.auth.getToken()
            }
        }).pipe(
        retry(3),
        map(res => res.plan)
        );
    }

    public updatePlan(plan: Plan): Observable<Plan> {
        return this.http
        .put<PlanResponse>(`${this.rootUrl}/plans/${plan.ID}`, {
            englishName: plan.englishName,
            turkishName: plan.trukishName,
            englishDescription: plan.engishDescription,
            turkishDescription: plan.turkishDescription,
            price: plan.price,
            countOfOffers: plan.countOfOffers,
            image: plan.image,
            isDefault: plan.isDefault,
        }, {
            headers: {
                Token: this.auth.getToken()
            }
        } )
        .pipe(
            retry(3),
            map(res => res.plan)
        );
    }

    public adminUpgradePlan(planID: number, userID: number): Observable<Plan> {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/admin/upgrade-plan?userID=${userID}&planID=${planID}`, {}, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.subscription.plan)
        );
    }

    public associatePlanWithCategory(planID: number, categoryID: number): Observable<string> {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/admin/associate-plan-category?planID=${planID}&categoryID=${categoryID}`, {}, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.success)
        );
    }

    public removePlanCategoryAssociation(planID: number, categoryID: number): Observable<string> {
        return this.http
        .delete<SuccessResponse>(`${this.rootUrl}/admin/associate-plan-category?planID=${planID}&categoryID=${categoryID}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.success)
        );
    }

    public getCategoriesOfPlan(planID: number): Observable<Category[]> {
        return this.http
        .get<CategoriesResponse>(`${this.rootUrl}/admin/categories?planID=${planID}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.categories)
        );
    }
}
