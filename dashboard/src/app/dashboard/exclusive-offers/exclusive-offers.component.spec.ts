import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ExclusiveOffersComponent } from './exclusive-offers.component';

describe('ExclusiveOffersComponent', () => {
  let component: ExclusiveOffersComponent;
  let fixture: ComponentFixture<ExclusiveOffersComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ExclusiveOffersComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ExclusiveOffersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
