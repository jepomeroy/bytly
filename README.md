# Bytly

A simple [Bitly](https://bitly.com/)-style URL shortening microservice, written in Go.

## Requirements

- Go 1.25+
- PostgreSQL

## Configuration

Bytly is configured via environment variables, loaded from a `.env` file in the project root (see `.gitignore` — this file is not committed).

| Variable      | Description                          |
|---------------|---------------------------------------|
| `DB_HOST`     | PostgreSQL host                       |
| `DB_PORT`     | PostgreSQL port                       |
| `DB_USER`     | PostgreSQL user                       |
| `DB_PASSWORD` | PostgreSQL password                   |
| `DB_NAME`     | PostgreSQL database name              |
| `API_PORT`    | Port the API server listens on        |

Example `.env`:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=bytly
API_PORT=5000
```

The database schema is auto-migrated on startup.

## Running

```sh
go run .
```

Or build and run the binary:

```sh
go build -o bytly .
./bytly
```

## API

All request/response bodies are JSON.

| Method | Path          | Description                                            |
|--------|---------------|----------------------------------------------------------|
| POST   | `/bytly`      | Create a shortened URL                                   |
| GET    | `/bytly`      | List all shortened URLs                                  |
| GET    | `/bytly/:id`  | Get a shortened URL by ID                                 |
| PATCH  | `/bytly`      | Update a shortened URL                                   |
| DELETE | `/bytly/:id`  | Delete a shortened URL by ID                              |
| GET    | `/r/:bytly`   | Redirect to the target URL and increment its click count |

### Bytly object

| Field      | Type   | Description                                                        |
|------------|--------|----------------------------------------------------------------------|
| `id`       | uint64 | Auto-generated identifier                                           |
| `redirect` | string | The destination URL                                                 |
| `bytly`    | string | The short code                                                      |
| `clicked`  | uint64 | Number of times the short URL has been visited                     |
| `random`   | bool   | If `true` on create, an 8-character random short code is generated |

### Create a shortened URL

```sh
curl -X POST http://localhost:5000/bytly \
  -H 'Content-Type: application/json' \
  -d '{"redirect": "https://example.com", "bytly": "example", "random": false}'
```

To have a short code generated automatically, set `"random": true` (the supplied `bytly` value is ignored in that case).
