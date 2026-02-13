# Justfile for trykkeri-api (run `just --list` to see all recipes)

# Run full stack (app + Grafana/Loki/Promtail) with watch (rebuild on change)
default:
    docker compose --profile observability watch

# Run the application locally
run:
    go run ./cmd/server

# Run tests
test:
    go test ./...

# Build Docker image
docker-build:
    docker build -t trykkeri-api:latest .

# Run Docker container (maps 8080:8080)
docker-run:
    docker run -p 8080:8080 trykkeri-api:latest

# Run Docker container with dev env (human-readable logs)
docker-run-dev:
    docker run -p 8080:8080 \
        -e JSON_LOGS=false \
        trykkeri-api:latest

# Run tests (unit tests; run locally)
docker-test:
    go test ./...

# Clean build artifacts
clean:
    go clean -cache
    rm -f trykkeri-api
