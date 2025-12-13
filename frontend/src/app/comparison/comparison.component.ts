import { Component, OnInit, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';

interface Movie {
  id: number;
  title: string;
  year: number;
  director: string;
  poster?: string;
}

interface Top4Set {
  id: number;
  movies: Movie[];
}

interface Comparison {
  id: number;
  set_a: Top4Set;
  set_b: Top4Set;
  votes_a: number;
  votes_b: number;
}

@Component({
  selector: 'app-comparison',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="comparison-container" *ngIf="comparison">
      <div class="sets-grid">
        <div 
          class="set-card left-side" 
          (click)="selectCard('a'); vote(comparison.set_a.id)"
          [class.clickable]="!voting"
          [class.selected]="selectedCard === 'a'">
          <div class="movies-grid">
            <div class="movie-card" *ngFor="let movie of comparison.set_a.movies">
              <div class="movie-image-wrapper">
                <img [src]="movie.poster || '/assets/placeholder.png'" [alt]="movie.title" />
                <div class="movie-overlay" [style.background-color]="getRandomColor(movie.id)">
                  <h3>{{ movie.title }}</h3>
                  <p>{{ movie.year }} • {{ movie.director }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="vs-divider">VS</div>

        <div 
          class="set-card right-side" 
          (click)="selectCard('b'); vote(comparison.set_b.id)"
          [class.clickable]="!voting"
          [class.selected]="selectedCard === 'b'">
          <div class="movies-grid">
            <div class="movie-card" *ngFor="let movie of comparison.set_b.movies">
              <div class="movie-image-wrapper">
                <img [src]="movie.poster || '/assets/placeholder.png'" [alt]="movie.title" />
                <div class="movie-overlay" [style.background-color]="getRandomColor(movie.id)">
                  <h3>{{ movie.title }}</h3>
                  <p>{{ movie.year }} • {{ movie.director }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div *ngIf="loading">Loading comparison...</div>
    <div *ngIf="error">{{ error }}</div>
  `,
  styles: [`
    :host {
      display: block;
      height: calc(100vh - 200px);
      overflow: hidden;
    }
    .comparison-container {
      padding: 10px;
      height: 100%;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
    .sets-grid {
      display: grid;
      grid-template-columns: 1fr auto 1fr;
      gap: 20px;
      align-items: stretch;
      flex: 1;
      min-height: 0;
    }
    .set-card {
      border: 2px solid #ddd;
      border-radius: 12px;
      padding: 15px;
      background: white;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      transition: all 0.3s;
      display: flex;
      flex-direction: column;
      min-height: 100%;
      overflow: hidden;
    }
    .set-card.clickable {
      cursor: pointer;
    }
    .set-card.clickable.selected {
      border-color: #2196F3;
      box-shadow: 0 8px 16px rgba(33, 150, 243, 0.3);
      transform: translateY(-2px);
    }
    .movies-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 10px;
      flex: 1;
      min-height: 0;
    }
    .movie-card {
      text-align: center;
      display: flex;
      flex-direction: column;
      min-height: 0;
    }
    .movie-image-wrapper {
      position: relative;
      width: 100%;
      aspect-ratio: 2/3;
      overflow: hidden;
      border-radius: 8px;
    }
    .movie-card img {
      width: 100%;
      height: 100%;
      object-fit: cover;
      border-radius: 8px;
      display: block;
    }
    .movie-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      color: white;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      opacity: 0;
      transition: opacity 0.3s ease;
      border-radius: 8px;
      padding: 10px;
      text-align: center;
    }
    .movie-image-wrapper:hover .movie-overlay {
      opacity: 1;
    }
    .movie-overlay h3 {
      margin: 0 0 8px 0;
      font-size: 1rem;
      font-weight: bold;
      line-height: 1.3;
    }
    .movie-overlay p {
      margin: 0;
      font-size: 0.85rem;
      color: rgba(255, 255, 255, 0.9);
    }
    .vs-divider {
      font-size: 1.5rem;
      font-weight: bold;
      color: #666;
      display: flex;
      align-items: center;
    }
  `]
})
export class ComparisonComponent implements OnInit {
  comparison: Comparison | null = null;
  loading = true;
  error: string | null = null;
  voting = false;
  selectedCard: 'a' | 'b' | null = null;
  private apiUrl = 'http://localhost:8080/api/v1';
  private retryCount = 0;
  private maxRetries = 3;
  private isLoading = false; // Prevent concurrent requests

  constructor(private http: HttpClient) {}

  @HostListener('window:keydown', ['$event'])
  handleKeyDown(event: KeyboardEvent) {
    if (this.voting || !this.comparison) return;
    
    if (event.key === 'ArrowLeft') {
      event.preventDefault();
      this.selectCard('a');
      this.vote(this.comparison.set_a.id);
    } else if (event.key === 'ArrowRight') {
      event.preventDefault();
      this.selectCard('b');
      this.vote(this.comparison.set_b.id);
    }
  }

  selectCard(card: 'a' | 'b') {
    this.selectedCard = card;
  }

  getRandomColor(movieId: number): string {
    // Use movie ID as seed for consistent colors per movie
    // Golden angle approximation for good color distribution
    const hue = (movieId * 137.508) % 360;
    // Return hsla with 0.6 opacity - opacity will be controlled by CSS hover
    return `hsla(${hue}, 70%, 30%, 0.6)`;
  }

  ngOnInit() {
    this.loadComparison();
  }

  loadComparison() {
    // Prevent concurrent requests
    if (this.isLoading || this.voting) {
      return;
    }
    
    if (this.retryCount >= this.maxRetries) {
      this.error = 'Failed to load comparison after multiple attempts. Please refresh the page.';
      this.loading = false;
      this.isLoading = false;
      return;
    }
    
    this.isLoading = true;
    this.loading = true;
    this.error = null;
    this.retryCount++;
    
    this.http.get<Comparison>(`${this.apiUrl}/comparison`).subscribe({
      next: (data) => {
        this.isLoading = false;
        // Validate response has required fields
        if (data && data.set_a && data.set_b && 
            Array.isArray(data.set_a.movies) && Array.isArray(data.set_b.movies) &&
            data.set_a.movies.length > 0 && data.set_b.movies.length > 0) {
          this.comparison = data;
          this.selectedCard = null; // Reset selection on new comparison
          this.loading = false;
          this.retryCount = 0; // Reset on success
        } else {
          this.error = 'Invalid comparison data received';
          this.loading = false;
          console.error('Invalid data structure:', data);
        }
      },
      error: (err) => {
        this.isLoading = false;
        this.error = 'Failed to load comparison';
        this.loading = false;
        console.error('API Error:', err);
        // Don't retry automatically to prevent loops
      }
    });
  }

  vote(winnerSetId: number) {
    if (!this.comparison || this.voting || this.isLoading) return;

    this.voting = true;
    this.http.post(`${this.apiUrl}/votes`, {
      comparison_id: this.comparison.id,
      winner_set_id: winnerSetId
    }).subscribe({
      next: () => {
        this.voting = false;
        // Small delay before reloading to prevent rapid requests
        setTimeout(() => {
          this.loadComparison();
        }, 300);
      },
      error: (err) => {
        console.error('Failed to vote', err);
        alert('Failed to submit vote. Please try again.');
        this.voting = false;
      }
    });
  }
}

