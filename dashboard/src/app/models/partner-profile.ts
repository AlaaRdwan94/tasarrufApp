import { City } from './city';

export class PartnerProfile {
    public ID: number;
    public CreatedAt: Date;
    public partnerID: number;
    public approved: boolean;
    public discountValue: number;
    public categroryID: number;
    public mainBranchAddress: string;
    public phone: string;
    public city: City;
    public country: string;
    public brandName: string;
    public offerDiscription: string;
    public licenceURL?: string;
}
