import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';
import {City} from '../models/city';

interface CitiesResponse {
    cities: Array<City>;
}

interface CityReponse {
    city: City;
}

@Injectable({
    providedIn: 'root'
})
export class CityService {

    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }
    private rootUrl: string;


    /**
     * getCities returns all the cities
     */
    public getCities(): Observable<Array<City>> {
        return this.http
        .get<CitiesResponse>(`${this.rootUrl}/public/cities`, {
            headers: {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.cities)
        );
    }


    /**
     * editCity updates the city with the given details
     */
    public editCity(id: number, updated: City): Observable<City> {
        return this.http
        .put<CityReponse>(`${this.rootUrl}/city/${id}`, {
            englishName: updated.englishName,
            turkishName: updated.turkishName,
        }, {
            headers: {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.city)
        );
    }

    /**
     * getCityByID returns the city with the given ID
     */
    public getCityByID(id: number): Observable<City> {
        return this.http
        .get<CityReponse>(`${this.rootUrl}/city/${id}`, {
            headers: {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.city)
        );
    }

    public createCity(newCity: City ): Observable<City> {
        return this.http
        .post<CityReponse>(`${this.rootUrl}/city`, {
            englishName: newCity.englishName,
            turkishName: newCity.turkishName,
        }, {
            headers: {Token: this.auth.getToken()}
        })
        .pipe(
            retry(3),
            map(res => res.city)
        );
    }

    public deleteCity(id: number): Observable<City> {
        return this.http
        .delete<CityReponse>(`${this.rootUrl}/city/${id}`, {
            headers: {
                Token: this.auth.getToken()
            }
        })
        .pipe(
            retry(3),
            map(res => res.city)
        );
    }
}
