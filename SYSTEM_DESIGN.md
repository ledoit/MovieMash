# MovieMash - System Design

## Overview
A MovieMash-style application where users compare two Letterboxd top 4 movie sets and vote for their preference. Votes are aggregated to build a collaborative filtering recommendation engine.

## Architecture

### Components

1. **Frontend (Angular)**
   - Comparison UI (side-by-side top 4s)
   - Voting interface
   - User authentication
   - Recommendation display

2. **API Gateway (Go)**
   - REST API endpoints
   - Request routing
   - Authentication middleware
   - Rate limiting

3. **Data Pipeline**
   - Letterboxd scraper (Python)
   - Data normalization
   - Movie metadata enrichment

4. **Database (PostgreSQL)**
   - Movies catalog
   - Top 4 sets
   - Votes history
   - User profiles
   - Recommendation scores

5. **ML Service (PyTorch)**
   - Collaborative filtering model
   - Training pipeline
   - Recommendation generation
   - Model serving

## Data Flow

### Vote Flow
1. User votes on comparison → Angular frontend
2. Frontend → Go API (POST /api/votes)
3. Go API → PostgreSQL (persist vote directly)
4. ML Service polls PostgreSQL or triggered via HTTP
5. PyTorch service → PostgreSQL (update recommendations)

### Comparison Generation
1. Go API receives request for comparison
2. Check in-memory cache for active comparisons
3. If miss, query PostgreSQL for two random top 4 sets
4. Cache in memory with TTL
5. Return to frontend

### Recommendation Flow
1. User requests recommendations
2. Go API queries PostgreSQL for user's recommendation scores
3. If stale/missing, trigger PyTorch service
4. PyTorch generates recommendations
5. Store in PostgreSQL
6. Return to frontend

## Database Schema

### Tables
- `movies`: Movie metadata (id, title, year, director, genres, letterboxd_id)
- `top4_sets`: Scraped top 4 sets (id, user_letterboxd_id, movies[], scraped_at)
- `votes`: User votes (id, user_id, comparison_id, winner_set_id, timestamp)
- `comparisons`: Active comparisons (id, set_a_id, set_b_id, votes_a, votes_b, created_at)
- `users`: User accounts (id, username, email, created_at)
- `recommendations`: Generated recommendations (user_id, movie_id, score, updated_at)
- `movie_features`: ML features (movie_id, features_json)

## Technology Stack

- **Frontend**: Angular 18+, TypeScript, RxJS
- **Backend API**: Go 1.21+, Gin/Echo framework
- **Database**: PostgreSQL 15+ (local installation)
- **ML**: PyTorch 2.0+, Python 3.11+
- **Scraping**: Python, BeautifulSoup/Scrapy

## Key Algorithms

### Comparison Selection
- ELO-based matching (similar popularity)
- Avoid recent comparisons for same user
- Ensure diversity in movie genres

### Recommendation Model
- Matrix factorization (collaborative filtering)
- User-item interaction matrix
- Implicit feedback from votes
- Regular updates via database polling

## Scalability Considerations

- Horizontal scaling: Stateless Go services
- In-memory cache can be replaced with Redis if needed
- PostgreSQL read replicas
- CDN for Angular static assets
- Load balancer for API gateway

## Security

- JWT authentication
- Rate limiting per user/IP
- Input validation and sanitization
- SQL injection prevention
- CORS configuration
- HTTPS/TLS

