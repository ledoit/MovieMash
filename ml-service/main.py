"""
PyTorch-based recommendation engine for MovieMash
Uses collaborative filtering with matrix factorization
"""

import os
import json
import logging
import time
from typing import List, Dict, Tuple
import numpy as np
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
import psycopg2
from psycopg2.extras import RealDictCursor

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class MovieRecommendationModel(nn.Module):
    """Matrix factorization model for collaborative filtering"""
    
    def __init__(self, num_users: int, num_movies: int, embedding_dim: int = 50):
        super(MovieRecommendationModel, self).__init__()
        self.user_embedding = nn.Embedding(num_users, embedding_dim)
        self.movie_embedding = nn.Embedding(num_movies, embedding_dim)
        
        # Initialize embeddings
        nn.init.normal_(self.user_embedding.weight, std=0.01)
        nn.init.normal_(self.movie_embedding.weight, std=0.01)
    
    def forward(self, user_ids: torch.Tensor, movie_ids: torch.Tensor) -> torch.Tensor:
        user_emb = self.user_embedding(user_ids)
        movie_emb = self.movie_embedding(movie_ids)
        return (user_emb * movie_emb).sum(dim=1)


class VoteDataset(Dataset):
    """Dataset for vote-based interactions"""
    
    def __init__(self, votes: List[Dict]):
        self.votes = votes
    
    def __len__(self):
        return len(self.votes)
    
    def __getitem__(self, idx):
        vote = self.votes[idx]
        return (
            torch.LongTensor([vote['user_id']]),
            torch.LongTensor([vote['movie_id']]),
            torch.FloatTensor([vote['score']])
        )


class RecommendationEngine:
    """Main recommendation engine"""
    
    def __init__(self, db_url: str):
        self.db_url = db_url
        self.model = None
        self.user_to_idx = {}
        self.movie_to_idx = {}
        self.idx_to_user = {}
        self.idx_to_movie = {}
        
    def connect_db(self):
        """Connect to PostgreSQL"""
        return psycopg2.connect(self.db_url)
    
    def load_data(self) -> Tuple[List[Dict], int, int]:
        """Load votes and build user/movie mappings"""
        conn = self.connect_db()
        cur = conn.cursor(cursor_factory=RealDictCursor)
        
        # Load votes
        cur.execute("""
            SELECT v.user_id, t.movie_ids, v.winner_set_id, v.comparison_id
            FROM votes v
            JOIN comparisons c ON v.comparison_id = c.id
            JOIN top4_sets t ON (c.set_a_id = t.id OR c.set_b_id = t.id)
            WHERE v.timestamp > NOW() - INTERVAL '30 days'
        """)
        
        votes = []
        users = set()
        movies = set()
        
        for row in cur.fetchall():
            user_id = row['user_id']
            winner_set_id = row['winner_set_id']
            movie_ids = row['movie_ids']
            
            users.add(user_id)
            
            # Create positive interactions for winner set
            for movie_id in movie_ids:
                if movie_id:
                    movies.add(movie_id)
                    votes.append({
                        'user_id': user_id,
                        'movie_id': movie_id,
                        'score': 1.0  # Positive interaction
                    })
        
        cur.close()
        conn.close()
        
        # Build mappings
        self.user_to_idx = {uid: idx for idx, uid in enumerate(sorted(users))}
        self.movie_to_idx = {mid: idx for idx, mid in enumerate(sorted(movies))}
        self.idx_to_user = {idx: uid for uid, idx in self.user_to_idx.items()}
        self.idx_to_movie = {idx: mid for mid, idx in self.movie_to_idx.items()}
        
        logger.info(f"Loaded {len(votes)} votes, {len(users)} users, {len(movies)} movies")
        
        return votes, len(users), len(movies)
    
    def train_model(self, epochs: int = 10, batch_size: int = 64, lr: float = 0.01):
        """Train the recommendation model"""
        votes, num_users, num_movies = self.load_data()
        
        if len(votes) == 0:
            logger.warning("No votes found, skipping training")
            return
        
        # Update indices in votes
        for vote in votes:
            vote['user_id'] = self.user_to_idx[vote['user_id']]
            vote['movie_id'] = self.movie_to_idx[vote['movie_id']]
        
        # Initialize model
        self.model = MovieRecommendationModel(num_users, num_movies, embedding_dim=50)
        criterion = nn.MSELoss()
        optimizer = optim.Adam(self.model.parameters(), lr=lr)
        
        # Create dataset and dataloader
        dataset = VoteDataset(votes)
        dataloader = DataLoader(dataset, batch_size=batch_size, shuffle=True)
        
        # Training loop
        self.model.train()
        for epoch in range(epochs):
            total_loss = 0
            for user_ids, movie_ids, scores in dataloader:
                optimizer.zero_grad()
                predictions = self.model(user_ids.squeeze(), movie_ids.squeeze())
                loss = criterion(predictions, scores.squeeze())
                loss.backward()
                optimizer.step()
                total_loss += loss.item()
            
            logger.info(f"Epoch {epoch+1}/{epochs}, Loss: {total_loss/len(dataloader):.4f}")
    
    def generate_recommendations(self, user_id: int, top_k: int = 20) -> List[Tuple[int, float]]:
        """Generate recommendations for a user"""
        if self.model is None:
            logger.error("Model not trained yet")
            return []
        
        if user_id not in self.user_to_idx:
            logger.warning(f"User {user_id} not in training data")
            return []
        
        self.model.eval()
        user_idx = self.user_to_idx[user_id]
        
        # Get all movies
        all_movie_indices = torch.LongTensor(list(range(len(self.movie_to_idx))))
        user_indices = torch.LongTensor([user_idx] * len(all_movie_indices))
        
        with torch.no_grad():
            scores = self.model(user_indices, all_movie_indices)
        
        # Get top K recommendations
        top_scores, top_indices = torch.topk(scores, min(top_k, len(scores)))
        
        recommendations = []
        for score, idx in zip(top_scores, top_indices):
            movie_id = self.idx_to_movie[idx.item()]
            recommendations.append((movie_id, score.item()))
        
        return recommendations
    
    def save_recommendations_to_db(self, user_id: int, recommendations: List[Tuple[int, float]]):
        """Save recommendations to PostgreSQL"""
        conn = self.connect_db()
        cur = conn.cursor()
        
        # Delete old recommendations
        cur.execute("DELETE FROM recommendations WHERE user_id = %s", (user_id,))
        
        # Insert new recommendations
        for movie_id, score in recommendations:
            cur.execute("""
                INSERT INTO recommendations (user_id, movie_id, score, updated_at)
                VALUES (%s, %s, %s, NOW())
                ON CONFLICT (user_id, movie_id) 
                DO UPDATE SET score = %s, updated_at = NOW()
            """, (user_id, movie_id, score, score))
        
        conn.commit()
        cur.close()
        conn.close()
        
        logger.info(f"Saved {len(recommendations)} recommendations for user {user_id}")


def poll_votes_and_train():
    """Poll database for new votes and trigger model updates"""
    engine = RecommendationEngine(os.getenv('DATABASE_URL'))
    last_vote_id = 0
    batch_size = 100  # Retrain after N votes
    poll_interval = 60  # Poll every 60 seconds
    
    logger.info("ML service started, polling database for votes...")
    
    while True:
        try:
            conn = engine.connect_db()
            cur = conn.cursor()
            
            # Get new votes since last check
            cur.execute("""
                SELECT COUNT(*) FROM votes 
                WHERE id > %s
            """, (last_vote_id,))
            
            new_vote_count = cur.fetchone()[0]
            
            if new_vote_count > 0:
                logger.info(f"Found {new_vote_count} new votes")
                
                # Get latest vote ID
                cur.execute("SELECT MAX(id) FROM votes")
                last_vote_id = cur.fetchone()[0] or 0
                
                # Retrain model if we have enough new votes
                if new_vote_count >= batch_size:
                    logger.info(f"Retraining model after {new_vote_count} new votes...")
                    engine.train_model(epochs=5)
                    logger.info("Model retraining complete")
            
            cur.close()
            conn.close()
            
            time.sleep(poll_interval)
            
        except Exception as e:
            logger.error(f"Error polling votes: {e}")
            time.sleep(poll_interval)


if __name__ == "__main__":
    import sys
    
    if len(sys.argv) > 1 and sys.argv[1] == "train":
        # One-time training
        engine = RecommendationEngine(os.getenv('DATABASE_URL', 
            'postgres://user:password@localhost:5432/proji?sslmode=disable'))
        engine.train_model(epochs=10)
    else:
        # Continuous polling
        poll_votes_and_train()

