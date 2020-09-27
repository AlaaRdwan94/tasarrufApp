import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map, retry } from 'rxjs/operators';
import {environment} from './../../environments/environment';
import {Offer} from '../models/offer';

interface OffersCountResponse {
  count: number;
}

interface OffersResponse {
    offers: Array<Offer>;
}

@Injectable({
  providedIn: 'root'
})
export class OffersService {

    constructor(private http: HttpClient, private auth: AuthService) {
        this.rootUrl = environment.apiUrl;
    }
    private rootUrl: string;

  /**
   * getTransactiosnCount returns a promise with offers count.
   */
  public getOffersCount(): Observable < number > {
    return this.http
    .get<OffersCountResponse>(`${this.rootUrl}/admin/count/offers`, {
        headers: { Token: this.auth.getToken() }
      }).pipe(
        retry(3),
        map(res => res.count)
      );
  }

  /**
   * getOffers returns a promise with a list of all consumed offers
   */
  public getOffers(): Observable<Array<Offer>> {
      return this.http
      .get<OffersResponse>(`${this.rootUrl}/admin/offers`, {
          headers: {Token: this.auth.getToken()}
      }).pipe(
      retry(3),
      map(res => res.offers)
      );
  }

  /**
   * getOffersOfCustomer returns the offers of the given customer
   */
  public getOffersOfCustomer(customerID: number) {
      return this.http
      .get<OffersResponse>(`${this.rootUrl}/admin/offers/customer/${customerID}`, {
          headers: { Token: this.auth.getToken() }
      })
      .pipe(
          retry(3),
          map(res => res.offers)
      );
  }
}
