import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';
import {Category} from '../models/category';


interface CategoriesResponse {
    categories: Array<Category>;
}

interface CategoryResponse {
    category: Category;
}

@Injectable({
    providedIn: 'root'
})
export class CategoriesService {
    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }
    private rootUrl: string;

    /**
     * getCategories returns an array of all categories
     */
    public getCategories(): Observable < Array < Category >> {
        return this.http
        .get<CategoriesResponse>(`${this.rootUrl}/category`, {
            headers: {Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.categories)
        );
    }

    /**
     * deleteCategoy deletes the category with the given id
     */
    public deleteCategoy(id: number): Observable<Category> {
        return this.http
        .delete<CategoryResponse>(`${this.rootUrl}/category/${id}`, {
            headers: {Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.category)
        );
    }

    /**
     * editCategory sends a request to edit the category with the given ID
     */
    public editCategory(id: number, englishName: string, turkishName: string): Observable<Category> {
        return this.http
        .put<CategoryResponse>(`${this.rootUrl}/category/${id}`, {
            englishName,
            turkishName,
        }, {
            headers : {
                Token: this.auth.getToken()
            }
        }).pipe(
        retry(3),
        map(res => res.category)
        );
    }

    /**
     * createCategory creates a new category
     */
    public createCategory(englishName: string, turkishName: string): Observable<Category> {
        return this.http
        .post<CategoryResponse>(`${this.rootUrl}/category`, {
            englishName,
            turkishName,
        }, {
            headers : {Token: this.auth.getToken()}
        }).pipe(
        retry(3),
        map(res => res.category)
        );
    }
}
