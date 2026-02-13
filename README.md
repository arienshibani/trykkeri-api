# Trykkeri API 🖨️

[![Go 1.22](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)](https://go.dev/)
[![OpenAPI 3.0](https://img.shields.io/badge/OpenAPI-3.0-6BA539?logo=openapi-initiative)](https://swagger.io/specification/)
[![Docker](https://img.shields.io/badge/Docker-ready-2496ED?logo=docker)](https://www.docker.com/)
[![Grafana](https://img.shields.io/badge/Grafana-dashboard-F46800?logo=grafana)](https://grafana.com/)

REST API for turning raw HTML into PDF files. 

## Features ✨

- **Grafana dashboard** - Preconfigured with a custom dashboard for monitoring usage and errors (when run with the observability stack).
- **Scalar UI** - Interactive API docs for trying different HTML and query parameters.
- **Tunable output** - Margins, page size, filename, DPI, orientation, background printing, grayscale etc. All using query parameters.

## Usage 🚀

Send a `POST` request to `/print` with your HTML as the request body.

```bash
curl http://localhost:8080/print \
  --request POST \
  --header 'Content-Type: text/html' \
  --header 'Accept: application/pdf' \
  --data '<h1 style="color: red; text-align: center">Hello world!</h1>'
```

The response is a PDF. To save it to a file, add `-o output.pdf` or use the `filename` query parameter to control the suggested download name.

### Optional query parameters 🔧

| Parameter | Type | Description |
| ----------- | ------ | ------------- |
| `filename` | string | Suggested filename in `Content-Disposition` (default: `document.pdf`) |
| `base_url` | string | Base URL for resolving relative links and assets in the HTML |
| `page_size` | string | e.g. `A4`, `Letter` |
| `portrait` | boolean | `true` = portrait, `false` = landscape |
| `margin_top_mm` | integer | Top margin in mm |
| `margin_right_mm` | integer | Right margin in mm |
| `margin_bottom_mm` | integer | Bottom margin in mm |
| `margin_left_mm` | integer | Left margin in mm |
| `dpi` | integer | Output DPI (e.g. `300`) |
| `print_background` | boolean | Include CSS background graphics |
| `grayscale` | boolean | Render in grayscale |

Example with options:

```bash
curl 'http://localhost:8080/print?filename=report.pdf&page_size=A4&margin_top_mm=20&dpi=150' \
  --request POST \
  --header 'Content-Type: text/html' \
  --data '<html><body><h1>Report</h1></body></html>'
```

## Quickstart 🏁

Ensure you have the following installed

- [Docker](https://www.docker.com/)
- [Go](https://go.dev/)
- [Justfile](https://github.com/casey/just) (optional, but recommended)

### 1. Clone the repository

```bash
git clone https://github.com/trykkeri/trykkeri-api.git
cd trykkeri-api
```

### 2. Spin up the services

**API only** (Go, no Docker):

```bash
just run
# or: go run ./cmd/server
```

**Full stack** (API + Grafana, Loki, Promtail in Docker, with live reload):

```bash
just watch
# or: docker compose --profile observability watch
```

**Pre-built image** (from [GitHub Container Registry](https://github.com/trykkeri/trykkeri-api/pkgs/container/trykkeri-api)):

```bash
docker run -p 8080:8080 ghcr.io/trykkeri/trykkeri-api:latest
```

Images are published on every push to `main` and when you create a release. Use `:latest`, `:sha-<commit>`, or a version tag (e.g. `:v1.0.0`) after a release.

### 3. Access the services

- 🔍 **API / Scalar UI** — <http://localhost:8080> (available with either `just run` or `just watch`).
- 🪵 **Grafana dashboard** — <http://localhost:3000> (only when you use `just watch`).

## Configuration 🔧

The service can be configured using environment variables. When you run the stack with Docker Compose, set these in a `.env` file in the project root—the same file used for Grafana (e.g. `GRAFANA_ADMIN_USER`). See [.env.example](.env.example) for a template.

| Variable | Description | Default |
| ---------- | ------------- | ------- |
| `PORT` | The port the service listens on | `8080` |
| `JSON_LOGS` | Whether to log in JSON format | `false` |
| `MAX_BODY_BYTES` | The maximum body size in bytes | `2000000` |
| `RENDER_TIMEOUT_MS` | The timeout in milliseconds for rendering a PDF | `30000` |
| `WKHTMLTOPDF_PATH` | The path to the wkhtmltopdf binary | `wkhtmltopdf` |
| `ALLOW_NET` | Whether to allow network access | `false` |

## Screenshots 📸

### API
<img width="1452" height="1279" alt="image" src="https://github.com/user-attachments/assets/a6affe80-b3e3-45f6-be4e-8b7f7d6a96df" />


### Grafana
<img width="1127" height="1142" alt="image" src="https://github.com/user-attachments/assets/51758579-b8f2-4c86-b05b-dc1f552a9540" />
