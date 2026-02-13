# Justfile for trykkeri-api (run `just --list` to see all recipes)

# Run full stack (app + Grafana/Loki/Promtail) with watch (rebuild on change)
default:
    just --list    

# Run the whole stack in watch mode. Rebuild on code changes.
watch:
    docker compose --profile observability watch

# Run the application locally
run:
    go run ./cmd/server

# Run tests
test:
    go test ./...