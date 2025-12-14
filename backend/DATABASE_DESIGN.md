# MovieMash Database Design

## Overview

The database uses a **normalized relational design** that efficiently stores movies and top 4 sets without data duplication.

## Schema Structure

### 1. `movies` Table (Master Movie Data)
- **Purpose**: Single source of truth for all movie information
- **Fields**: `id`, `title`, `year`, `director`, `poster_url`, `letterboxd_id`, `created_at`
- **Benefits**: 
  - No duplication - each movie stored once
  - Easy to update (e.g., fix poster URL in one place)
  - Consistent data across all top 4 sets

### 2. `top4_sets` Table (User Selections)
- **Purpose**: Stores which 4 movies a user selected
- **Fields**: `id`, `user_id`, `movie_ids` (integer array), `created_at`
- **Key Design**: Uses PostgreSQL array `INTEGER[]` to store 4 movie IDs
- **Benefits**:
  - References movies by ID (foreign key relationship)
  - No duplication of movie data
  - Efficient storage (just 4 integers per set)

### 3. Relationship Pattern

```
movies (1) ←─── (many) top4_sets.movie_ids[]
```

- Each `top4_sets.movie_ids` array contains 4 references to `movies.id`
- When you query a top 4 set, you JOIN to get full movie details
- Multiple top 4 sets can reference the same movie (economies of scale!)

## Example: How It Works

### Scenario: 3 users all have "The Godfather" in their top 4

**Without normalization (bad):**
```
top4_set_1: {movie: "The Godfather", year: 1972, director: "Coppola", poster: "..."}
top4_set_2: {movie: "The Godfather", year: 1972, director: "Coppola", poster: "..."}
top4_set_3: {movie: "The Godfather", year: 1972, director: "Coppola", poster: "..."}
```
**Storage**: 3x movie data = wasteful!

**With normalization (current design):**
```
movies table:
  id=1, title="The Godfather", year=1972, director="Coppola", poster="..."

top4_sets table:
  id=1, movie_ids=[1, 5, 12, 23]  ← references movie id=1
  id=2, movie_ids=[1, 8, 15, 31]  ← references movie id=1
  id=3, movie_ids=[1, 3, 7, 19]   ← references movie id=1
```
**Storage**: 1x movie data + 3x array references = efficient!

## Query Pattern

When fetching a comparison, the API:
1. Gets `top4_sets` with `movie_ids` arrays
2. JOINs to `movies` table to get full details:
   ```sql
   SELECT m.* 
   FROM movies m 
   WHERE m.id = ANY(top4_sets.movie_ids)
   ```
3. Returns complete movie objects with posters, directors, etc.

## Benefits of This Design

✅ **No Data Duplication**: Each movie stored once
✅ **Easy Updates**: Fix poster URL once, all top 4 sets see the update
✅ **Scalable**: As you scrape more top 4s, popular movies are reused
✅ **Data Integrity**: Foreign key constraints ensure valid references
✅ **Efficient Queries**: PostgreSQL arrays + JOINs are fast

## Future Considerations

When scraping Letterboxd top 4s:
1. **Check if movie exists** before inserting (by `letterboxd_id` or `title+year`)
2. **Insert movie** if new (with poster from TMDB)
3. **Create top4_set** referencing existing movie IDs
4. **Result**: Popular movies (appearing in many top 4s) are stored once, referenced many times

This is exactly how relational databases are designed to work! 🎯

