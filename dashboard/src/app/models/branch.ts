import {City} from 'src/app/models/city';
import {User} from 'src/app/models/user';
import {Category} from 'src/app/models/category';

export class Branch {
    country: string;
    cityID: string;
    city: City;
    address: string;
    phone: string;
    mobile: string;
    ownerID: number;
    owner: User;
    categoryID: number;
    category: Category;
}
