import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <div class="container">
      <header>
        <h1>MovieMash</h1>
        <p>Compare top 4s and discover your next favorite movie</p>
      </header>
      <nav class="tabs">
        <a routerLink="/" routerLinkActive="active" [routerLinkActiveOptions]="{exact: true}">Gladiators</a>
        <a routerLink="/recommendations" routerLinkActive="active">Recommendations</a>
        <a routerLink="/leaderboard" routerLinkActive="active">Leaderboard</a>
      </nav>
      <main>
        <router-outlet></router-outlet>
      </main>
    </div>
  `,
  styles: [`
    :host {
      display: block;
      min-height: 100vh;
    }
    .container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 10px 20px;
      min-height: 100vh;
      display: flex;
      flex-direction: column;
    }
    header {
      text-align: center;
      margin-bottom: 10px;
      flex-shrink: 0;
    }
    h1 {
      font-size: 2rem;
      margin-bottom: 5px;
    }
    .tabs {
      display: flex;
      justify-content: center;
      gap: 20px;
      margin-bottom: 10px;
      border-bottom: 2px solid #ddd;
      flex-shrink: 0;
    }
    .tabs a {
      padding: 12px 24px;
      text-decoration: none;
      color: #666;
      font-weight: 500;
      border-bottom: 3px solid transparent;
      margin-bottom: -2px;
      transition: all 0.3s;
    }
    .tabs a:hover {
      color: #2196F3;
    }
    .tabs a.active {
      color: #2196F3;
      border-bottom-color: #2196F3;
    }
    main {
      flex: 1;
      min-height: 0;
    }
  `]
})
export class AppComponent {
  title = 'moviemash';
}

