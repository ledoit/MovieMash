import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface Movie {
  id: number;
  title: string;
  year: number;
  director: string;
  poster: string;
  win_rate?: number;
  appearances?: number;
  wins?: number;
}

interface Top4Set {
  id: number;
  movies: Movie[];
  win_rate?: number;
  appearances?: number;
  wins?: number;
}

@Component({
  selector: 'app-leaderboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './leaderboard.component.html',
  styleUrl: './leaderboard.component.scss',
})
export class LeaderboardComponent implements OnInit {
  viewMode: 'top4' | 'movies' = 'top4';
  top4Sets: Top4Set[] = [];
  movies: Movie[] = [];
  isLoading = false;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadTop4Sets();
    this.loadMovies();
  }

  loadTop4Sets() {
    this.isLoading = true;
    this.http.get<Top4Set[]>(`${environment.apiUrl}/v1/leaderboard/top4`).subscribe({
      next: (data) => {
        this.top4Sets = data;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading top4 leaderboard:', error);
        this.isLoading = false;
      },
    });
  }

  loadMovies() {
    this.http.get<Movie[]>(`${environment.apiUrl}/v1/leaderboard/movies`).subscribe({
      next: (data) => {
        this.movies = data;
      },
      error: (error) => {
        console.error('Error loading movies leaderboard:', error);
      },
    });
  }

  switchView(mode: 'top4' | 'movies') {
    this.viewMode = mode;
  }

  openMoviePage(movie: Movie) {
    // Construct IMDB search URL
    const searchQuery = encodeURIComponent(`${movie.title} ${movie.year}`);
    const imdbUrl = `https://www.imdb.com/find?q=${searchQuery}`;
    window.open(imdbUrl, '_blank');
  }
}
