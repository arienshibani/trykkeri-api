# Trykkeri

REST API for converting HTML to PDF.

## Quick Start

Clone the repository and build the Docker image.

```bash
docker build -t trykkeri-api:latest .
```

Then run the built image:

```bash
docker run -p 8080:8080 trykkeri-api:latest
```

Or use the [just](https://github.com/casey/just) runner: `just docker-build` then `just docker-run`.

**Local development (Go):**

```bash
go run ./cmd/server
```

Or build and run: `just run` (see below).

## API

- `POST /print` - Convert HTML to PDF (accepts `text/html` body)
  - Simply send a POST request with HTML content in the body to receive a PDF in response.
  - Query parameters can be used to customize the output PDF (see below).

## Example

```bash
curl -X POST "http://localhost:8080/print?filename=test.pdf&page_size=A4" \
  -H 'Content-Type: text/html' \
  -o out.pdf \
  -d '<html><head><meta charset="UTF-8"></head><body><h1>Test PDF</h1><p>This is a test document.</p></body></html>'
```

## Query Parameters

All optional:

- `filename` - Output filename (default: document.pdf)
- `base_url` - Base URL for resolving relative assets
- `page_size` - A4, Letter, etc. (default: A4)
- `margin_*_mm` - Margins in mm (default: 10)
- `dpi` - Render DPI (default: 300)
- `print_background` - Include backgrounds (default: true)
- `grayscale` - Grayscale mode (default: false)
- `portrait` - true = portrait, false = landscape (default: true)

### Environment variables:

- `PORT` (default: 8080)
- `MAX_BODY_BYTES` (default: 2000000)
- `RENDER_TIMEOUT_MS` (default: 30000)
- `WKHTMLTOPDF_PATH` (default: wkhtmltopdf)
- `ALLOW_NET` (default: false)
- `CORS_ORIGINS` (comma-separated, unset = permissive)
- `JSON_LOGS` (default: false) – structured JSON logging for production/Loki

## Monitoring and logging

The repo includes an optional observability stack (Grafana, Loki, Promtail) behind the `observability` profile.

- **App only:** `docker compose up --build` (or `docker compose watch` for rebuild-on-change).
- **App + Grafana/Loki/Promtail:** `docker compose --profile observability up --build`.

- **API:** <http://localhost:8080>  
- **Grafana:** <http://localhost:3000> (default login: `admin` / `admin`). To override, set `GRAFANA_ADMIN_USER` and `GRAFANA_ADMIN_PASSWORD` in the environment or in a `.env` file (see `.env.example`). Grafana opens on the **trykkeri-api logs** dashboard by default.

The app container is named `trykkeri-api`; logs appear in Loki under that name. The app runs with `JSON_LOGS=true` so log fields are structured for filtering.

**Note:** Promtail reads container logs from the Docker daemon. On Linux it works out of the box. On Docker Desktop (Mac/Windows), if logs do not appear, ensure the Docker socket and container log directory are available to the Promtail container (they are mounted in the default `docker-compose.yml`; some Docker Desktop setups may need the stack to run in a Linux context).