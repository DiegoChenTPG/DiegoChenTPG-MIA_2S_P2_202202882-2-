
import { bootstrapApplication } from '@angular/platform-browser';
import { AppComponent } from './app/app.component';
import { appConfig } from './app/app.config';
import { provideRouter } from '@angular/router';
import { routes } from './app/app.routes'; // Asegúrate de importar tus rutas

bootstrapApplication(AppComponent, {
    providers: [
        appConfig.providers,
        provideRouter(routes)
    ]
});
