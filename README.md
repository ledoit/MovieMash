# MovieMash

A FaceMash-style application for comparing Letterboxd top 4 movie sets.

## Tech Stack

- **Frontend**: Angular 18 (standalone components)
- **Backend**: Go (Gin)
- **Database**: PostgreSQL

## Development

### Frontend

```bash
# Install dependencies
npm install

# Start development server
npm start
# or
npm run dev

# Build for production
npm run build
```

The frontend runs on `http://localhost:4200` and proxies API requests to `http://localhost:8080`.

### Database and ingest

```bash
docker compose up -d
cd backend
go run ./cmd/migrate
go run ./cmd/ingest -json ../ingest/letterboxd.json
go run ./cmd/api
```

`ingest/letterboxd.json` is a public-profile export (top 4s and a page of logged films). `go run ./cmd/ingest -fetch someuser` only prints favorites when Letterboxd returns the page. A Cloudflare challenge stops that command.

### Backend

See backend documentation for Go server setup.

## Project Structure

```
proji/
├── src/              # Angular source files
│   ├── app.component.ts
│   ├── main.ts
│   ├── index.html
│   └── styles.css
├── angular.json      # Angular CLI configuration
├── tsconfig.json     # TypeScript configuration
└── proxy.conf.json   # API proxy configuration
```

