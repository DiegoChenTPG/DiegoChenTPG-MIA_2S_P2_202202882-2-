// auth.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private apiUrl = 'http://3.23.105.151:8080/api/login'; // Cambia la URL según tu API, regresar a localhost cuando se deje de trabajar con AWS y asi con las demas URL
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
        return this.http.post('http://3.23.105.151:8080/api/logout', {}).subscribe( //regresar a localhost cuando se deje de trabajar con AWS
            () => {
                this.loggedInUser = '';
                localStorage.removeItem('loggedInUser'); // Limpiar el localStorage
                console.log("Sesión cerrada exitosamente.");
            },
            error => console.error("Error cerrando sesión:", error)
        );
    }

    obtenerDiscos(): Observable<string[]> {
        return this.http.get<string[]>('http://3.23.105.151:8080/api/discos'); //regresar a localhost cuando se deje de trabajar con AWS
    }

    obtenerParticiones(disco: string): Observable<string[]> {
        return this.http.get<string[]>(`http://3.23.105.151:8080/api/particiones?disco=${disco}`); //regresar a localhost cuando se deje de trabajar con AWS
    }

}
