import { Component } from '@angular/core';
import { RouterOutlet, RouterModule, Router} from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CommandService } from './command.service';
import { CommonModule } from '@angular/common';



@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterModule, FormsModule, CommonModule],

  templateUrl: './app.component.html',
  /* para las url es ./nombre_archivo.html
    y lo mismo ./nombre_archivo.css
  */
  styleUrl: './app.component.css'
})

export class AppComponent {
  tittle = "frontend_proyecto"
  commandInput: string = '';
  commandOutput: string = '';

  constructor(private commandService: CommandService, public router:Router) { }
  // Metodo que se ejecuta cuando el usuario hace click en "Ejecutar Comando"
  executeCommand() {
    this.commandService.executeCommand(this.commandInput)
      .subscribe(response => {
        this.commandOutput = response.output;
      }, error => {
        console.error('Error ejecutando el comando:', error);
        this.commandOutput = 'Error al ejecutar el comando.';
      });
  }

  loadFile(event: any) {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e: any) => {
      this.commandInput = e.target.result;
    };

    reader.readAsText(file);
  }

  /*
  Creo que aqui se agregan comandos simples sin necesidad
  de programar gran cosa
  */

  clearCommand() {
    // Comando para limpiar la entrada
    this.commandInput = " ";
    this.commandOutput = " ";
  }

}
