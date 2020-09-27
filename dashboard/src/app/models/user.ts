import {SafeUrl} from '@angular/platform-browser';
import { PartnerProfile } from './partner-profile';
export class User {
  ID: number;
  CreatedAt: Date;
  UpdatedAt: Date;
  DeletedAt: Date;
  email: string;
  firstName: string;
  lastName: string;
  mobile: string;
  phone: string;
  hashedPassword: string;
  accountType: string;
  dateOfBirth: string;
  profileImageURL: string | SafeUrl;
  profileImageKey: string;
  partnerProfile ?: PartnerProfile;
  Subscription?: {
      planID: number;
      expireDate: Date;
  };
  active: boolean;
}
