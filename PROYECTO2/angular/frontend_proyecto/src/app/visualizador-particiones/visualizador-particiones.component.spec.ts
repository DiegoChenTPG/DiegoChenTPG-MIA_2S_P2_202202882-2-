import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VisualizadorParticionesComponent } from './visualizador-particiones.component';

describe('VisualizadorParticionesComponent', () => {
  let component: VisualizadorParticionesComponent;
  let fixture: ComponentFixture<VisualizadorParticionesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VisualizadorParticionesComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VisualizadorParticionesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
