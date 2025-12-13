import { Routes } from '@angular/router';
import { ComparisonComponent } from './comparison/comparison.component';
import { RecommendationsComponent } from './recommendations/recommendations.component';
import { LeaderboardComponent } from './leaderboard/leaderboard.component';

export const routes: Routes = [
  { path: '', component: ComparisonComponent },
  { path: 'recommendations', component: RecommendationsComponent },
  { path: 'leaderboard', component: LeaderboardComponent },
];

