import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SupportInfoComponent } from './support-info.component';

describe('SupportInfoComponent', () => {
  let component: SupportInfoComponent;
  let fixture: ComponentFixture<SupportInfoComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SupportInfoComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SupportInfoComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
