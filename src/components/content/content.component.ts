import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ComparisonComponent } from '../comparison/comparison.component';
import { RecommendationsComponent } from '../recommendations/recommendations.component';
import { LeaderboardComponent } from '../leaderboard/leaderboard.component';

@Component({
  selector: 'app-content',
  standalone: true,
  imports: [CommonModule, ComparisonComponent, RecommendationsComponent, LeaderboardComponent],
  templateUrl: './content.component.html',
  styleUrl: './content.component.scss',
})
export class ContentComponent {
  activeTab: 'leaderboard' | 'gladiators' | 'recommendations' = 'gladiators';

  setActiveTab(tab: 'leaderboard' | 'gladiators' | 'recommendations') {
    this.activeTab = tab;
  }
}

