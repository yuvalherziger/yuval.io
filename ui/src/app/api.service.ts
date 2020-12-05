import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../environments/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  baseUrl: string;
  host: string;
  apiVersion = 'v1beta1';
  requestOpts = {};

  constructor(private http: HttpClient) {
    if (!environment.production) {
      const scheme = environment.apiScheme;
      const host = environment.apiHost;
      const port = environment.apiPort;
      this.host = `${scheme}://${host}:${port}`;
      this.baseUrl = `${this.host}/api/${this.apiVersion}`;
    } else {
      this.host = '';
      this.baseUrl = `/api/${this.apiVersion}`;
    }
  }

  public executeCommand(input: string, width: number): Observable<{ input: string, output: string }> {
    const url = `${this.baseUrl}/cmd`;
    return this.http.post<any>(url, { input, width });
  }
}
