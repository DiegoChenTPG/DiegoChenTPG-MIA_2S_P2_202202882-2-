import { Routes } from '@angular/router';
import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';
import { VisualizadorDiscosComponent } from './visualizador-discos/visualizador-discos.component';
import { ConsolaLogeadaComponent } from './consola-logeada/consola-logeada.component';
import { VisualizadorParticionesComponent } from './visualizador-particiones/visualizador-particiones.component';


export const routes: Routes = [
    //{ path: '', component: AppComponent }, //Comentamos esto para evitar el duplicado del inicio
    { path: 'login', component: LoginComponent },
    { path: 'visualizador-discos', component: VisualizadorDiscosComponent},
    { path: 'consola-logeada', component: ConsolaLogeadaComponent},
    { path: 'visualizador-particiones/:disco', component: VisualizadorParticionesComponent },
    { path: '**', redirectTo: '' }
];