# Detailed Architecture

## System Components

### 1. Frontend (Angular)
- **Location**: `frontend/`
- **Purpose**: User interface for comparisons and recommendations
- **Key Features**:
  - Comparison view (side-by-side top 4s)
  - Voting mechanism
  - Recommendations display
  - User authentication UI

### 2. Backend API (Go)
- **Location**: `backend/`
- **Framework**: Gin
- **Endpoints**:
  - `GET /api/v1/comparison` - Get a comparison pair
  - `POST /api/v1/votes` - Submit a vote
  - `GET /api/v1/recommendations/:user_id` - Get user recommendations
  - `GET /api/v1/movies/:id` - Get movie details
  - `POST /api/v1/auth/register` - User registration
  - `POST /api/v1/auth/login` - User login

### 3. Database (PostgreSQL)
- **Schema**: See `backend/migrations/001_initial_schema.sql`
- **Tables**:
  - `movies` - Movie catalog
  - `top4_sets` - Scraped top 4 sets
  - `comparisons` - Active comparison pairs
  - `votes` - User votes
  - `users` - User accounts
  - `recommendations` - Generated recommendations
  - `movie_features` - ML features

### 4. ML Service (PyTorch)
- **Location**: `ml-service/`
- **Model**: Matrix factorization for collaborative filtering
- **Training**: Batch training on vote events
- **Output**: User-movie recommendation scores

### 5. Scraper (Python)
- **Location**: `scraper/`
- **Purpose**: Scrape Letterboxd profiles for top 4s
- **Output**: Populates `top4_sets` table

## Data Flow Diagrams

### Vote Flow
```
User → Angular → Go API → PostgreSQL
                ↓
            ML Service (polls DB)
```

### Comparison Generation
```
Angular → Go API → In-memory cache (check)
                    ↓ (miss)
                 PostgreSQL → Cache → Angular
```

### Recommendation Generation
```
Angular → Go API → PostgreSQL (check existing)
                    ↓ (stale/missing)
                 ML Service → PostgreSQL → Angular
```

## Deployment Architecture

### Development
- Local PostgreSQL installation
- Local services (Go API, Angular dev server, Python services)

### Production (Suggested)
- **Frontend**: CDN/Static hosting (Vercel, Netlify)
- **Backend**: Cloud functions or VPS
- **Database**: Managed PostgreSQL (RDS, Cloud SQL, or VPS)
- **ML Service**: Cloud function or VPS service

## Performance Considerations

1. **Comparison Selection**: Use ELO-based matching for balanced comparisons
2. **Caching Strategy**: Cache active comparisons in-memory with TTL
3. **Model Training**: Batch training every N votes (configurable)
4. **Database Indexing**: Indexes on foreign keys and frequently queried columns
5. **API Rate Limiting**: Per-user and per-IP rate limits

## Security

1. **Authentication**: JWT tokens with refresh mechanism
2. **Input Validation**: Validate all user inputs
3. **SQL Injection**: Use parameterized queries
4. **CORS**: Configure for frontend domain only
5. **Rate Limiting**: Prevent abuse

