import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';


interface CommandResponse {
  output: string;
}


@Injectable({
  providedIn: 'root'
})
export class CommandService {
  private apiUrl = 'http://3.23.105.151:8080/api/execute';
  
  constructor(private http: HttpClient) {}

  executeCommand(command: string): Observable<CommandResponse> {
    return this.http.post<CommandResponse>(this.apiUrl, { command });
  }
}
