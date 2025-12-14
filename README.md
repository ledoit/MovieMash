# MovieMash

A FaceMash-style application for comparing Letterboxd top 4 movie sets.

## Tech Stack

- **Frontend**: Angular 18 (standalone components)
- **Backend**: Go (Gin framework)
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

