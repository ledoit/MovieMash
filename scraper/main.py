"""
Letterboxd Top 4 Scraper
Scrapes user profiles to extract their top 4 movies
"""

import os
import time
import logging
import psycopg2
from psycopg2.extras import RealDictCursor
import requests
from bs4 import BeautifulSoup
from dotenv import load_dotenv

load_dotenv()

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class LetterboxdScraper:
    """Scraper for Letterboxd top 4 movies"""
    
    def __init__(self, db_url: str):
        self.db_url = db_url
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36'
        })
    
    def connect_db(self):
        """Connect to PostgreSQL"""
        return psycopg2.connect(self.db_url)
    
    def scrape_user_top4(self, username: str) -> list:
        """Scrape top 4 movies from a Letterboxd user profile"""
        url = f"https://letterboxd.com/{username}/"
        
        try:
            response = self.session.get(url, timeout=10)
            response.raise_for_status()
            
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Find top 4 section (Letterboxd displays top 4 on profile)
            # This selector may need adjustment based on actual HTML structure
            top4_elements = soup.select('.poster-list li')
            
            movie_ids = []
            for element in top4_elements[:4]:  # Take first 4
                # Extract movie ID from data attributes or links
                # This is a placeholder - actual implementation depends on Letterboxd HTML
                movie_link = element.find('a')
                if movie_link:
                    href = movie_link.get('href', '')
                    # Extract movie identifier from href
                    # Example: /film/movie-name/
                    movie_ids.append(href)
            
            return movie_ids[:4]
        
        except Exception as e:
            logger.error(f"Error scraping {username}: {e}")
            return []
    
    def get_or_create_movie(self, letterboxd_id: str, title: str = None, year: int = None) -> int:
        """Get or create a movie in the database"""
        conn = self.connect_db()
        cur = conn.cursor()
        
        # Try to find existing movie
        cur.execute("SELECT id FROM movies WHERE letterboxd_id = %s", (letterboxd_id,))
        result = cur.fetchone()
        
        if result:
            cur.close()
            conn.close()
            return result[0]
        
        # Create new movie
        cur.execute("""
            INSERT INTO movies (letterboxd_id, title, year, created_at)
            VALUES (%s, %s, %s, NOW())
            RETURNING id
        """, (letterboxd_id, title, year))
        
        movie_id = cur.fetchone()[0]
        conn.commit()
        cur.close()
        conn.close()
        
        return movie_id
    
    def save_top4_set(self, username: str, movie_ids: list):
        """Save a top 4 set to the database"""
        if len(movie_ids) != 4:
            logger.warning(f"Invalid top 4 for {username}: {movie_ids}")
            return
        
        conn = self.connect_db()
        cur = conn.cursor()
        
        # Check if set already exists
        cur.execute("""
            SELECT id FROM top4_sets 
            WHERE user_letterboxd_id = %s 
            ORDER BY scraped_at DESC 
            LIMIT 1
        """, (username,))
        
        existing = cur.fetchone()
        
        if existing:
            logger.info(f"Top 4 already exists for {username}")
            cur.close()
            conn.close()
            return
        
        # Insert new top 4 set
        cur.execute("""
            INSERT INTO top4_sets (user_letterboxd_id, movie_ids, scraped_at)
            VALUES (%s, %s, NOW())
        """, (username, movie_ids))
        
        conn.commit()
        cur.close()
        conn.close()
        
        logger.info(f"Saved top 4 for {username}: {movie_ids}")
    
    def scrape_batch(self, usernames: list, delay: float = 2.0):
        """Scrape multiple users with rate limiting"""
        for username in usernames:
            logger.info(f"Scraping {username}...")
            movie_ids = self.scrape_user_top4(username)
            
            if movie_ids:
                # Convert letterboxd IDs to database movie IDs
                db_movie_ids = []
                for lb_id in movie_ids:
                    movie_id = self.get_or_create_movie(lb_id)
                    db_movie_ids.append(movie_id)
                
                self.save_top4_set(username, db_movie_ids)
            
            time.sleep(delay)  # Rate limiting


if __name__ == "__main__":
    db_url = os.getenv('DATABASE_URL', 
        'postgres://user:password@localhost:5432/proji?sslmode=disable')
    
    scraper = LetterboxdScraper(db_url)
    
    # Example: Scrape some popular Letterboxd users
    # In production, this would be a scheduled job or queue-based
    usernames = [
        # Add Letterboxd usernames here
    ]
    
    if usernames:
        scraper.scrape_batch(usernames, delay=3.0)
    else:
        logger.info("No usernames provided. Add usernames to scrape.")

