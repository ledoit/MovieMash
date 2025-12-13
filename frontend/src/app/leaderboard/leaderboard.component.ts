import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';

interface MovieRanking {
  rank: number;
  movie: {
    id: number;
    title: string;
    year: number;
    director: string;
    poster?: string;
  };
  wins: number;
  appearances: number;
  win_rate: number;
}

@Component({
  selector: 'app-leaderboard',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="leaderboard-container">
      <h2>Global Movie Rankings</h2>
      <p class="subtitle">Movies ranked by their performance across all comparisons</p>
      
      <div class="leaderboard-table" *ngIf="rankings.length > 0">
        <div class="table-header">
          <div class="col-rank">Rank</div>
          <div class="col-movie">Movie</div>
          <div class="col-stats">Wins</div>
          <div class="col-stats">Appearances</div>
          <div class="col-stats">Win Rate</div>
        </div>
        
        <div class="table-row" *ngFor="let ranking of rankings" [class.top-three]="ranking.rank <= 3">
          <div class="col-rank">
            <span class="medal" *ngIf="ranking.rank === 1">🥇</span>
            <span class="medal" *ngIf="ranking.rank === 2">🥈</span>
            <span class="medal" *ngIf="ranking.rank === 3">🥉</span>
            <span class="rank-number" *ngIf="ranking.rank > 3">{{ ranking.rank }}</span>
          </div>
          <div class="col-movie">
            <div class="movie-info">
              <img 
                [src]="ranking.movie.poster || '/assets/placeholder.png'" 
                [alt]="ranking.movie.title"
                class="movie-poster" />
              <div class="movie-details">
                <h3>{{ ranking.movie.title }}</h3>
                <p>{{ ranking.movie.year }} • {{ ranking.movie.director }}</p>
              </div>
            </div>
          </div>
          <div class="col-stats">{{ ranking.wins }}</div>
          <div class="col-stats">{{ ranking.appearances }}</div>
          <div class="col-stats">
            <span class="win-rate">{{ (ranking.win_rate * 100).toFixed(1) }}%</span>
          </div>
        </div>
      </div>
      
      <div *ngIf="loading" class="loading">Loading leaderboard...</div>
      <div *ngIf="error" class="error">{{ error }}</div>
    </div>
  `,
  styles: [`
    .leaderboard-container {
      padding: 20px;
      max-width: 1000px;
      margin: 0 auto;
    }
    h2 {
      text-align: center;
      font-size: 2rem;
      margin-bottom: 10px;
    }
    .subtitle {
      text-align: center;
      color: #666;
      margin-bottom: 30px;
    }
    .leaderboard-table {
      background: white;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0,0,0,0.1);
    }
    .table-header {
      display: grid;
      grid-template-columns: 80px 1fr 100px 120px 120px;
      gap: 20px;
      padding: 15px 20px;
      background: #f5f5f5;
      font-weight: bold;
      border-bottom: 2px solid #ddd;
    }
    .table-row {
      display: grid;
      grid-template-columns: 80px 1fr 100px 120px 120px;
      gap: 20px;
      padding: 15px 20px;
      border-bottom: 1px solid #eee;
      transition: background 0.2s;
    }
    .table-row:hover {
      background: #f9f9f9;
    }
    .table-row.top-three {
      background: linear-gradient(90deg, rgba(255,215,0,0.1) 0%, transparent 100%);
    }
    .col-rank {
      display: flex;
      align-items: center;
      gap: 8px;
    }
    .rank-number {
      font-weight: bold;
      font-size: 1.2rem;
    }
    .medal {
      font-size: 1.5rem;
    }
    .col-movie {
      display: flex;
      align-items: center;
    }
    .movie-info {
      display: flex;
      align-items: center;
      gap: 15px;
    }
    .movie-poster {
      width: 60px;
      height: 90px;
      object-fit: cover;
      border-radius: 4px;
    }
    .movie-details h3 {
      margin: 0 0 5px 0;
      font-size: 1rem;
    }
    .movie-details p {
      margin: 0;
      color: #666;
      font-size: 0.9rem;
    }
    .col-stats {
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 500;
    }
    .win-rate {
      color: #4CAF50;
      font-weight: bold;
    }
    .loading, .error {
      text-align: center;
      padding: 40px;
      color: #666;
    }
    .error {
      color: #f44336;
    }
    @media (max-width: 768px) {
      .table-header, .table-row {
        grid-template-columns: 60px 1fr;
        gap: 10px;
      }
      .col-stats {
        display: none;
      }
      .table-header .col-stats {
        display: none;
      }
    }
  `]
})
export class LeaderboardComponent implements OnInit {
  rankings: MovieRanking[] = [];
  loading = true;
  error: string | null = null;
  private apiUrl = 'http://localhost:8080/api/v1';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadLeaderboard();
  }

  loadLeaderboard() {
    this.loading = true;
    this.http.get<{ rankings: MovieRanking[] }>(`${this.apiUrl}/leaderboard`)
      .subscribe({
        next: (data) => {
          this.rankings = data.rankings;
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to load leaderboard';
          this.loading = false;
          console.error(err);
        }
      });
  }
}

