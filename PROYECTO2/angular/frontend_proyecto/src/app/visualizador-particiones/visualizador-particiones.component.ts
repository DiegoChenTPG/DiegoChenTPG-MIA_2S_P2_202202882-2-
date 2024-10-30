import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from '../servicios/auth.service';
import { CommonModule } from '@angular/common';
@Component({
  selector: 'app-visualizador-paprticiones',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './visualizador-particiones.component.html',
  styleUrl: './visualizador-particiones.component.css'
})
export class VisualizadorParticionesComponent implements OnInit {
  disco: string = ''; // Nombre del disco actual
  particiones: string[] = []; // Lista de particiones del disco

  constructor(private route: ActivatedRoute, private router: Router, private authService: AuthService) {}

  ngOnInit(): void {
    // Obtener el nombre del disco desde el parámetro de la ruta
    this.disco = this.route.snapshot.paramMap.get('disco') || '';

    // Llamar al backend para obtener las particiones de este disco
    this.obtenerParticiones(this.disco);
  }

  obtenerParticiones(disco: string) {
    this.authService.obtenerParticiones(disco).subscribe((data: string[]) => {
      this.particiones = data;
      console.log(this.particiones);
    });

  }

  regresar_discos(){
    this.router.navigate(['/visualizador-discos']);
  }
}
