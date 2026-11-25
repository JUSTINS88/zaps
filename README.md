# zaps

A minimal URL shortener built with Go, PostgreSQL, and `sqlc` — usable from both the terminal (CLI) and a lightweight embedded web UI.

Live demo: [zaps.my.id](https://zaps.my.id)

## Features
- Shorten URLs with custom or auto-generated 6-character slugs
- Destination preview page before leaving `zaps.my.id` so visitors can inspect the target URL
- Public link stats page (`/stats`) to view total clicks and link metadata
- Terminal CLI (`zaps`) with zero-login local machine key (`~/.zaps.json`), terminal QR code (`-q`), and PNG export (`-S`)
- URL (`http`/`https` only) and slug validation to prevent invalid redirects and reserved-path conflicts

## Tech Stack
- **Backend:** Go (`net/http`, `html/template`, `go:embed` single-binary web UI)
- **Database:** PostgreSQL, `sqlc`, `goose` migrations, bounded connection pooling
- **Performance:** Thread-safe in-memory lookup cache (`sync.RWMutex`) & async click tracking
- **CLI:** `pflag`, `skip2/go-qrcode`

## Install CLI
```bash
go install github.com/JUSTINS88/zaps/cmd/zaps@latest
```

## Project Structure
- `cmd/zapsd` — HTTP server (`POST /api/links`, `GET /api/links/{slug}/stats`, `GET /{slug}`, `GET /stats`, `GET /`)
- `cmd/zaps` — CLI client
- `internal/web` — Embedded frontend (`index.html`, `stats.html`, `preview.html`, `zaps-mark.svg`)
- `internal/store` — `sqlc` generated database layer
- `db/` — SQL migrations (`goose`) and queries

## Running Locally
1. Start PostgreSQL and run migrations in `db/migrations/`.
2. Copy `.env.example` to `.env` and set `DB_URL`.
3. Run the server:
   ```bash
   go run ./cmd/zapsd
   ```
4. Use the CLI:
   ```bash
   zaps https://example.com -s my-slug
   zaps https://example.com -q        # print QR code in terminal
   zaps https://example.com -q -S     # print QR code and save <slug>.png
   ```

## API
- `POST /api/links` — Create a short link (`400` on invalid URL/slug, `409` if slug is taken or reserved)
  ```json
  {
    "original_link": "https://example.com",
    "slug": "optional-slug",
    "machine_key": "your-key"
  }
  ```
- `GET /api/links/{slug}/stats` — Get click count and link metadata (`404` if not found)
- `GET /{slug}` — Render destination preview page and track clicks

## Benchmark (Local)

Tested locally against PostgreSQL on the direct `302 Found` redirect handler using [`hey`](https://github.com/rakyll/hey):

```bash
hey -n 1000 -c 50 -disable-redirects http://localhost:8080/<slug>
```

| Metric | Before Cache (DB) | After Cache (`sync.RWMutex`) |
| --- | --- | --- |
| Throughput | ~2,210 req/s | **~9,110 req/s** |
| Median Latency (p50) | 2.0 ms | **1.3 ms** |
| p95 Latency | 78.0 ms | **22.1 ms** |
| Status | 1000/1000 (`302 Found`, 0 errors) | 1000/1000 (`302 Found`, 0 errors) |
