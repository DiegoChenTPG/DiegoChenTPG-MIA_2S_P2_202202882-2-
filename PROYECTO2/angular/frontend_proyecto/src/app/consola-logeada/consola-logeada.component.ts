import { Component } from '@angular/core';
import { AuthService } from '../servicios/auth.service';
import { CommandService } from '../command.service';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { RouterOutlet, RouterModule, Router} from '@angular/router';

@Component({
  selector: 'app-consola-logeada',
  standalone: true,
  imports: [RouterOutlet, RouterModule, FormsModule, CommonModule],
  templateUrl: './consola-logeada.component.html',
  styleUrl: './consola-logeada.component.css'
})
export class ConsolaLogeadaComponent {
  nombre_usuario: string = ""

  commandInput: string = '';
  commandOutput: string = '';
  nombreUsuario: string = '';

  constructor(private authService: AuthService, private commandService: CommandService, private router: Router) {}
  

  ngOnInit(){
    this.nombre_usuario = this.authService.getLoggedInUser()
  }

  executeCommand() {
    this.commandService.executeCommand(this.commandInput).subscribe(response => {
      this.commandOutput = response.output;
    });
  }

  clearCommand() {
    this.commandInput = '';
    this.commandOutput = '';
  }

  loadFile(event: any) {
    const file = event.target.files[0];
    const reader = new FileReader();

    reader.onload = (e: any) => {
      this.commandInput = e.target.result;
    };

    reader.readAsText(file);
  }

  logout(){
    this.authService.logout()
    this.router.navigate(['/login']); 
  }


}
