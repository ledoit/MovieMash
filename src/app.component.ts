import { Component } from '@angular/core';
import { BannerComponent } from './components/banner/banner.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [BannerComponent], // Import components you want to use
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {}
