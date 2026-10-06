# Profiles

Profiles are Letterboxd identities stored in this app’s PostgreSQL. They are not a separate service, frontend, or auth vendor.

## Stack

| Concern | Choice |
|---|---|
| Browser | Angular 18. The Profiles leaderboard is a view on the existing Leaderboard tab. |
| API | Go with Gin, same process as votes, ingest, and model scoring. |
| Store | PostgreSQL. One `users` row per Letterboxd username. Login is a separate `accounts` row. |
| Public facts | Python scraper of public Letterboxd HTML. Letterboxd does not offer an official API that can list a profile’s top 4, film log, and bio for this app. |
| Posters | TMDB, through the existing Go fetcher. |
| Ratings | Elo on the profile, on its top 4, and on each movie, written in the vote transaction. |
| Recommendations | PyTorch, trained on vote outcomes, top-4 membership, and watch overlap. Go asks that model for scores. |
| Later simulators | Python workers that call the same vote route. They do not get their own stack. |

`users` stays the profile table because `top4_sets.user_id` already points there. Rows with `control = 'system'` (the anonymous voter, local seed owners) are not Letterboxd profiles and do not appear on the profiles leaderboard.

## Rows

`users` (profile):

- `letterboxd_username` unique, nullable for system rows
- `control`: `simulator`, `human`, or `system`
- `elo` default 1500, `match_count`, `wins`
- `claimed_at` set when a human takes the row

`accounts` (our login):

- email, password hash
- `user_id` null until that person claims a profile

`watches`:

- `(user_id, movie_id, watched_at, letterboxd_rating)`
- scraped film log, not MovieMash vote history

`votes`:

- `actor`: `guest`, `account`, or `simulator`
- `voter_key`: guest cookie, account id, or the profile id a simulator is acting as
- `user_id` nullable; set when the voter is a profile (simulator or claimed human)

`movies` and `top4_sets` gain the same Elo columns.

## What a simulator is

A simulator is this app casting votes as a real Letterboxd username, using that person’s public top 4 and film log. It never posts to Letterboxd. Those votes are tagged `actor = simulator`.

## Claim

A person creates an `accounts` row and signs in. To take a profile they put a one-time code in their public Letterboxd bio. The scraper reads that public page. On a match:

1. `accounts.user_id` attaches to that profile.
2. `control` becomes `human` and `claimed_at` is set.
3. Every `simulator` vote cast as that profile is deleted.
4. Elo, match counts, wins, and leaderboard position for profiles, top 4s, and movies are rebuilt from the votes that remain. That profile returns to the starting rating.
5. The scraped top 4 and film log stay. Those are the person’s public Letterboxd record, not the simulation.

Signing in without a matching bio code does not attach a profile and does not wipe anyone.
