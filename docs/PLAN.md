# MovieMash plan

Agent tracker. Update the status line when a chunk lands. Commits stay small. Phil sees subinitiatives, not every commit.

## Status

Current: 5 — `infra/` describes a Neon free-plan database. It is not applied: no `NEON_API_KEY` in the environment. Railway stays until that database answers. UI waits on Forge.
Next: apply Neon when a key exists, then a $0 place for the API process. PyTorch after the API can reach the database.

## Goal

https://movie-mash-six.vercel.app serves real Letterboxd top 4s. A vote updates Elo for the two sets, their movies, and the two owning profiles, then loads a pair that voter has not seen. Leaderboards order profiles, top 4s, and movies by Elo. Recommendations come from a PyTorch model trained on votes, top-4 membership, and watches. A fresh clone can run the frontend, API, database, ingest, and training from the README.

## Stack

- Browser: Angular 18, TypeScript, SCSS, already on Vercel project `movie-mash`.
- API: Go, Gin, `net/http` server replaced by Gin.
- Database: PostgreSQL on Neon’s free plan, described in `infra/`. Railway config leaves once that database is applied and the API is pointed at it.
- Infra: Terraform. Frontend stays on Vercel. API compute is still open: Fly and Cloud Run both want an account, and neither is applied.
- Model: PyTorch trains offline and writes a score file. The API process loads it. Training is a command in the README.
- Ingest: Go command reading a JSON export or a saved public HTML page. A direct fetch that receives a Cloudflare challenge stops. No challenge bypass.
- Profiles: see `docs/profiles.md`.

## Hosting constraint

Spend stays at $0. Terraform codifies the free tier; it does not replace the need for a vendor account. Do not attach a card or open a paid plan.

Target when credentials exist:

- Postgres: Neon free tier (no AWS RDS).
- API container: a scale-to-zero free compute that can run the Go binary and the small PyTorch score process. Fly.io or Cloud Run, whichever applies without a paid plan.
- Until then, local Postgres via Docker Compose and `go run`.

Vercel keeps the Angular build. `environment.prod.ts` changes only after the new API answers `GET /api/v1/comparison`.

## Subinitiatives

1. **Router.** Gin serves the current four routes with the same JSON. README names Gin because the module requires it. Cache and `backend/bin` untracked.
2. **Schema.** Elo columns, `watches`, `accounts`, vote `actor` / `voter_key`, profile `control`. Seed still boots.
3. **Elo.** A vote updates sets, movies (each movie in the winner against each movie in the loser), and the two profiles. Pair pick prefers unseen `voter_key`s. `GET /api/v1/leaderboard/profiles` orders by Elo. Guest cookie is the voter key until accounts exist.
4. **Ingest.** Public Letterboxd usernames → `users`, movies deduped on `letterboxd_id` else title+year, `top4_sets`, `watches`. TMDB posters. Production rankings from those profiles. Seed remains for an empty local database.
5. **Terraform.** Replace Railway. Apply only against a free tier. Point prod API URL after a live comparison response.
6. **Forge.** Scout pack in `docs/forge/scout.md` before any product chrome, layout, or empty-state change. Unsigned pack goes to Atelier and waits for a signed cut. Current templates stay until then.
7. **UI.** Implement the signed pack: Gladiators, three leaderboards, recommendations with posters, titles, and a reason, empty and error states, bundled poster fallback.
8. **PyTorch.** Train on votes, top-4 membership, and watch overlap. `GET /api/v1/recommendations?profile=` returns ranked movies and a short reason. Re-run command documented.
9. **Claim.** Email and password session. Bio-code claim. Simulator votes for that profile deleted and Elo rebuilt. See `docs/profiles.md`.
10. **Simulators.** After real profiles and Elo exist. Workers cast `actor=simulator` votes from a profile’s top 4 and film log, through the public vote route.

## Done when

- Live Gladiators shows two real Letterboxd top 4s, a click stores a vote, Elo moves, the next pair loads.
- Leaderboard has profiles, top 4s, and movies in Elo order.
- Recommendations for a real profile come from the trained model.
- README runs frontend, API, database, ingest, and training on a fresh clone.
- Claim wipes that profile’s simulator history and rating.
- Railway is not on the boot path.
