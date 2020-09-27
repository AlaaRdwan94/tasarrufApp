import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CititesComponent } from './citites.component';

describe('CititesComponent', () => {
  let component: CititesComponent;
  let fixture: ComponentFixture<CititesComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CititesComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CititesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
