import { Component, OnInit, OnDestroy, HostListener } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface Movie {
  id: number;
  title: string;
  year: number;
  director: string;
  poster: string;
}

interface Top4Set {
  id: number;
  movies: Movie[];
}

interface Comparison {
  id: number;
  set_a: Top4Set;
  set_b: Top4Set;
}

@Component({
  selector: 'app-comparison',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './comparison.component.html',
  styleUrl: './comparison.component.scss',
})
export class ComparisonComponent implements OnInit, OnDestroy {
  comparison: Comparison | null = null;
  isLoading = false;
  voting = false;
  selectedCard: 'left' | 'right' | null = null;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    // Only load comparison if we don't have one already
    if (!this.comparison) {
      this.loadComparison();
    }
  }

  ngOnDestroy() {
    // Cleanup if needed
  }

  @HostListener('window:keydown', ['$event'])
  handleKeyDown(event: KeyboardEvent) {
    if (this.voting || this.isLoading) return;

    if (event.key === 'ArrowLeft') {
      this.vote('left');
    } else if (event.key === 'ArrowRight') {
      this.vote('right');
    }
  }

  loadComparison() {
    if (this.isLoading) return;
    
    this.isLoading = true;
    this.selectedCard = null;

    this.http.get<Comparison>(`${environment.apiUrl}/v1/comparison`).subscribe({
      next: (data) => {
        this.comparison = data;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading comparison:', error);
        this.isLoading = false;
      },
    });
  }

  selectCard(side: 'left' | 'right') {
    if (this.voting || this.isLoading) return;
    this.selectedCard = side;
  }

  vote(side: 'left' | 'right') {
    if (this.voting || this.isLoading || !this.comparison) return;

    this.voting = true;
    this.selectedCard = side;

    const winnerSetId = side === 'left' 
      ? this.comparison.set_a.id 
      : this.comparison.set_b.id;

    this.http.post(`${environment.apiUrl}/v1/votes`, {
      comparison_id: this.comparison.id,
      winner_set_id: winnerSetId,
    }).subscribe({
      next: () => {
        // Small delay before loading next comparison
        setTimeout(() => {
          this.loadComparison();
          this.voting = false;
          this.selectedCard = null;
        }, 300);
      },
      error: (error) => {
        console.error('Error submitting vote:', error);
        this.voting = false;
        this.selectedCard = null;
      },
    });
  }

  getRandomColor(movieId: number): string {
    // Use movie ID as seed for consistent colors per movie
    // Golden angle approximation for good color distribution
    const hue = (movieId * 137.508) % 360;
    return `hsla(${hue}, 70%, 30%, 0.6)`;
  }
}
