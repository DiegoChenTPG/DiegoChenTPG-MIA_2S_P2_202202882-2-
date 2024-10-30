import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConsolaLogeadaComponent } from './consola-logeada.component';

describe('ConsolaLogeadaComponent', () => {
  let component: ConsolaLogeadaComponent;
  let fixture: ComponentFixture<ConsolaLogeadaComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConsolaLogeadaComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ConsolaLogeadaComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
