# Fetch Movie Posters

This script fetches movie poster URLs from The Movie Database (TMDB) API and updates your database.

## Setup

1. **Get a free TMDB API key:**
   - Go to https://www.themoviedb.org/
   - Sign up for a free account
   - Go to Settings → API
   - Request an API key (it's free and instant)

2. **Add API key to your `.env` file:**
   ```bash
   TMDB_API_KEY=your_api_key_here
   ```

## Usage

Run the script from the backend directory:

```bash
cd backend
go run cmd/fetch-posters/main.go
```

The script will:
- Find all movies in your database without poster URLs
- Query TMDB API for each movie (by title and year)
- Update the database with poster URLs
- Show progress for each movie

## Example Output

```
Found 20 movies without poster URLs

[1/20] Fetching poster for: The Shawshank Redemption (1994)... ✓ Updated
[2/20] Fetching poster for: The Godfather (1972)... ✓ Updated
[3/20] Fetching poster for: Pulp Fiction (1994)... ✓ Updated
...

✅ Successfully updated 18/20 movies with poster URLs

Poster URLs are now available in the database!
```

## Notes

- The script includes a 250ms delay between requests to be respectful to TMDB's API
- If a movie isn't found or has no poster, it will skip it
- Poster URLs are stored as full URLs (e.g., `https://image.tmdb.org/t/p/w500/poster.jpg`)
- The script only updates movies that don't already have poster URLs

