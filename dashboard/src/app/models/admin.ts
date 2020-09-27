import { SafeUrl } from '@angular/platform-browser';

/*
|*Admin class represents a profile of the logged in Admin.
 */
export class Admin {
  id: number;
  dateOfBirth: Date;
  email: string;
  firstName: string;
  lastName: string;
  mobile: string;
  profileImageURL: string;
  safeProfileImage?: SafeUrl;
  accountType: string;
  verified: boolean;
}
