import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OffersTableComponent } from './offers-table.component';

describe('OffersTableComponent', () => {
  let component: OffersTableComponent;
  let fixture: ComponentFixture<OffersTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OffersTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OffersTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
