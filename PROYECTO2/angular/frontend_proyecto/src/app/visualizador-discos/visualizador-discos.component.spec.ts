import { ComponentFixture, TestBed } from '@angular/core/testing';

import { VisualizadorDiscosComponent } from './visualizador-discos.component';

describe('VisualizadorDiscosComponent', () => {
  let component: VisualizadorDiscosComponent;
  let fixture: ComponentFixture<VisualizadorDiscosComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [VisualizadorDiscosComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(VisualizadorDiscosComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
