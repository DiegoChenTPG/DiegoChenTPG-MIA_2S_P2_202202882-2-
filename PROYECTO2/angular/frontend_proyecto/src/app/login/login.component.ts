// login.component.ts
import { Component } from '@angular/core';
import { AuthService } from '../servicios/auth.service'; // Ajusta la ruta según tu estructura
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'; // Importar CommonModule
@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrl: './login.component.css',
    standalone: true,
    imports: [FormsModule, CommonModule]
})
export class LoginComponent {
    ID_particion: string = '';
    nombre_usuario: string = '';
    password: string = '';
    MensajeLogin: string = '';

    constructor(private authService: AuthService, private router: Router) { }

    login() {
        const credentials = {
            ID_particion: this.ID_particion,
            nombre_usuario: this.nombre_usuario,
            password: this.password
        };

        this.authService.login(credentials).subscribe({
            next: (response) => {
                if (response.success) {
                    // Redirigir si el login es exitoso
                    this.authService.setLoggedInUser(this.nombre_usuario);
                    console.log('Inicio de sesión exitoso:', response);
                    this.MensajeLogin = "Sesion para el usuario: " + this.nombre_usuario + "iniciada correctamente"
                    this.router.navigate(['/consola-logeada']); // Redirige a la página principal
                } else {
                    // Mostrar mensaje de error
                    this.MensajeLogin = response.message;
                }
            },
            error: (err) => {
                console.error('Error en el inicio de sesión:', err);
                this.MensajeLogin = 'Error de conexión con el servidor.';
            }
        });
    }

    goToConsolaPrincipal() {
        this.router.navigate(['/']); // Redirige a la consola principal
    }
}
