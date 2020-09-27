import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { AuthService } from 'src/app/services/auth.service';
import { User } from 'src/app/models/user';
import { Branch } from 'src/app/models/branch';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';

interface CountResponse {
    count: number;
}

interface UsersResponse {
    users: Array<User>;
}

interface BranchesResponse {
    branches: Array<Branch>;
}

interface SuccessResponse {
    success: string;
    user: User;
}

@Injectable({
    providedIn: 'root'
})

/**
 * UsersService handles all backend requests regarding the users.
 */
export class UsersService {
    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }

    private rootUrl: string;

    /**
     * GetCustomersCount returns the count of all users in lendish.
     */
    getCustomersCount(): Observable<number> {
        const token: string = this.auth.getToken();
        return this.http
        .get<CountResponse>(`${this.rootUrl}/admin/count/customers`, {
            headers: { Token: token }
        })
        .pipe(
            retry(3),
            map(res => res.count)
        );
    }
    /**
     * GetPartnersCount returns the count of all users in lendish.
     */
    getPartnersCount(): Observable<number> {
        const token: string = this.auth.getToken();
        return this.http
        .get<CountResponse>(`${this.rootUrl}/admin/count/partners`, {
            headers: { Token: token }
        })
        .pipe(
            retry(3),
            map(res => res.count)
        );
    }

    /**
     * getCustomers returns an array of users with the given offset.
     * @param offset number is skipped
     */
    getCustomers(): Observable<Array<User>> {
        return this.http
        .get<UsersResponse>(`${this.rootUrl}/admin/customers` , {
            headers: { Token: this.auth.getToken() }
        })
        .pipe(
            retry(3),
            map(res => res.users)
        );
    }

    getPartners(): Observable<Array<User>> {
        return this.http
        .get<UsersResponse>(`${this.rootUrl}/admin/partners` , {
            headers: { Token: this.auth.getToken() }
        })
        .pipe(
            retry(3),
            map(res => res.users)
        );
    }

    getNotApprovedPartners(): Observable<Array<User>> {
        return this.http
        .get<UsersResponse>(`${this.rootUrl}/admin/partners/not-approved`, {
            headers: { Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.users)
        );
    }

    public approvePartner(partnerID: number) {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/admin/approve/${partnerID}`, null, {
            headers : {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.user)
        );
    }
    public rejectPartner(partnerID: number) {
        return this.http
        .delete<SuccessResponse>(`${this.rootUrl}/admin/user/${partnerID}`, {
            headers : {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.user)
        );
    }

    public getCustomer(customerID: number) {
        return this.http
        .get<SuccessResponse>(`${this.rootUrl}/admin/customer/${customerID}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public getPartner(partnerID: number) {
        return this.http
        .get<SuccessResponse>(`${this.rootUrl}/admin/partner/${partnerID}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public getExclusivePartners() {
        return this.http
        .get<UsersResponse>(`${this.rootUrl}/exclusive`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.users;
            })
        );
    }

    public removePartnerExclusive(id: number) {
        return this.http
        .delete<SuccessResponse>(`${this.rootUrl}/exclusive/${id}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public setPartnerExclusive(id: number) {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/exclusive/${id}`, {}, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public togglePartnerIsSharable(id: number) {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/admin/is-sharable/${id}`, {}, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public deleteUser(id: number) {
        return this.http
        .delete<SuccessResponse>(`${this.rootUrl}/admin/user/${id}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => {
                console.log(res);
                return res.user;
            })
        );
    }

    public searchUsers(searchTerm: string) {
        return this.http
        .get<UsersResponse>(`${this.rootUrl}/admin/users?q=${searchTerm}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe (
            retry(3),
            map(res => res.users)
        );
    }

    public toggleActiveProperty(userID: number) {
        return this.http
        .post<SuccessResponse>(`${this.rootUrl}/admin/activate-user/${userID}`, {}, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe (
            retry(3),
            map(res => res.user)
        );

    }
    public getBranchesOfPartner(partnerID: number) {
        return this.http
        .get<BranchesResponse>(`${this.rootUrl}/branches-by-owner/${partnerID}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.branches)
        );
    }
}


