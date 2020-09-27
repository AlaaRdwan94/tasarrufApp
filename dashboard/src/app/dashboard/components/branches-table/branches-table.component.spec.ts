import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { BranchesTableComponent } from './branches-table.component';

describe('BranchesTableComponent', () => {
  let component: BranchesTableComponent;
  let fixture: ComponentFixture<BranchesTableComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ BranchesTableComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(BranchesTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
