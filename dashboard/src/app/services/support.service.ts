import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { AuthService } from 'src/app/services/auth.service';
import { Support } from 'src/app/models/support';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';

interface SupportResponse {
    info: Support;
}

interface SupportSuccess {
    success: string;
    info: Support;
}


@Injectable({
    providedIn: 'root'
})
export class SupportService {

    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }
    private rootUrl: string;

    /**
     * returns support info from the backend
     */
    public getSupportInfo(): Observable<Support> {
        return this.http
        .get<SupportResponse>(`${this.rootUrl}/public/support-info`, {
            headers: { Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.info)
        );
    }

    /**
     * updates support info values
     */
    public updateSupportInfo(email: string , mobile: string) {
        return this.http
        .post<SupportResponse>(`${this.rootUrl}/support`, {
            mobile, email,
        }, {
            headers: { Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.info)
        );
    }
}
