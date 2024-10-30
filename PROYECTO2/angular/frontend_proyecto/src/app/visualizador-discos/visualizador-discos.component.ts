import { Component, OnInit } from '@angular/core';
import { AuthService } from '../servicios/auth.service';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-visualizador-discos',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './visualizador-discos.component.html',
  styleUrl: './visualizador-discos.component.css'
})
export class VisualizadorDiscosComponent implements OnInit {
  discos: string[] = [];

  constructor(private authService: AuthService, private router: Router) {}

  ngOnInit(): void {
    this.authService.obtenerDiscos().subscribe((data: string[]) => {
      if (Array.isArray(data)) {
        this.discos = data;
        console.log("Discos recibidos:", this.discos); // Aquí deberías ver un arreglo
      } else {
        console.error("Formato inesperado de datos:", data);
      }
    });
  }

  irAParticiones(disco: string): void {
    this.router.navigate(['/visualizador-particiones', disco]);
  }

  regresar_consola(){
    this.router.navigate(['/consola-logeada']);
  }

}
