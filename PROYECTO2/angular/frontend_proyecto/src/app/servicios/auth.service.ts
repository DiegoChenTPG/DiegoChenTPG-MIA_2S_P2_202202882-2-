// auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private apiUrl = 'http://localhost:8080/api/login'; // Cambia la URL según tu API
    private loggedInUser: string = ""

    constructor(private http: HttpClient) { }

    login(credentials: { ID_particion: string; nombre_usuario: string; password: string }): Observable<any> {
        return this.http.post(`${this.apiUrl}`, credentials);
    }

    setLoggedInUser(nombre_usuario: string) {
        this.loggedInUser = nombre_usuario;
        localStorage.setItem('loggedInUser', nombre_usuario)
    }

    getLoggedInUser(): string {
        if (!this.loggedInUser) {
            this.loggedInUser = localStorage.getItem('loggedInUser') || ''; // Recuperar de localStorage
        }
        return this.loggedInUser;
        //return this.loggedInUser;
    }

    // Método para cerrar sesión
    logout() {
        return this.http.post('http://localhost:8080/api/logout', {}).subscribe(
            () => {
                this.loggedInUser = '';
                localStorage.removeItem('loggedInUser'); // Limpiar el localStorage
                console.log("Sesión cerrada exitosamente.");
            },
            error => console.error("Error cerrando sesión:", error)
        );
    }

    obtenerDiscos(): Observable<string[]> {
        return this.http.get<string[]>('http://localhost:8080/api/discos');
    }

    obtenerParticiones(disco: string): Observable<string[]> {
        return this.http.get<string[]>(`http://localhost:8080/api/particiones?disco=${disco}`);
    }

}
