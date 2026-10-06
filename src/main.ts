import { bootstrapApplication } from '@angular/platform-browser';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { AppComponent } from './app.component';
import { credentialsInterceptor } from './credentials.interceptor';

bootstrapApplication(AppComponent, {
  providers: [provideHttpClient(withInterceptors([credentialsInterceptor]))],
});

