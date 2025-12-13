import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';

interface Movie {
  id: number;
  title: string;
  year: number;
  director: string;
  poster?: string;
  score?: number;
}

@Component({
  selector: 'app-recommendations',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="recommendations-container">
      <h2>Your Movie Recommendations</h2>
      <div class="movies-grid" *ngIf="recommendations.length > 0">
        <div class="movie-card" *ngFor="let movie of recommendations">
          <img [src]="movie.poster || '/assets/placeholder.png'" [alt]="movie.title" />
          <h3>{{ movie.title }}</h3>
          <p>{{ movie.year }} • {{ movie.director }}</p>
          <div class="score" *ngIf="movie.score">Score: {{ movie.score | number:'1.2-2' }}</div>
        </div>
      </div>
      <div *ngIf="loading">Loading recommendations...</div>
      <div *ngIf="error">{{ error }}</div>
    </div>
  `,
  styles: [`
    .recommendations-container {
      padding: 20px;
    }
    .movies-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
      gap: 20px;
      margin-top: 20px;
    }
    .movie-card {
      border: 1px solid #ddd;
      border-radius: 8px;
      padding: 15px;
      text-align: center;
    }
    .movie-card img {
      width: 100%;
      height: 300px;
      object-fit: cover;
      border-radius: 4px;
    }
    .score {
      margin-top: 10px;
      font-weight: bold;
      color: #4CAF50;
    }
  `]
})
export class RecommendationsComponent implements OnInit {
  recommendations: Movie[] = [];
  loading = true;
  error: string | null = null;
  userId: number = 1; // TODO: Get from auth
  private apiUrl = 'http://localhost:8080/api/v1';

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute
  ) {}

  ngOnInit() {
    this.route.params.subscribe(params => {
      if (params['user_id']) {
        this.userId = +params['user_id'];
      }
      this.loadRecommendations();
    });
  }

  loadRecommendations() {
    this.loading = true;
    this.http.get<{ recommendations: Movie[] }>(`${this.apiUrl}/recommendations/${this.userId}`)
      .subscribe({
        next: (data) => {
          this.recommendations = data.recommendations;
          this.loading = false;
        },
        error: (err) => {
          this.error = 'Failed to load recommendations';
          this.loading = false;
          console.error(err);
        }
      });
  }
}

