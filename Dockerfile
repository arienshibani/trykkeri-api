# Builder stage
FROM golang:1.22-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/trykkeri-api ./cmd/server

# Runtime stage
FROM debian:bookworm-slim

# Install wkhtmltopdf and fonts
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    fonts-noto \
    fonts-noto-cjk \
    fonts-liberation \
    fontconfig \
    libx11-6 \
    libxext6 \
    libxrender1 \
    libfontconfig1 \
    libjpeg62-turbo \
    libpng16-16 \
    libssl3 \
    ca-certificates \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -r appuser && \
    useradd -r -g appuser -u 1000 appuser && \
    mkdir -p /app && \
    chown -R appuser:appuser /app

# Copy binary from builder
COPY --from=builder /app/trykkeri-api /app/trykkeri-api

# Set working directory
WORKDIR /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set environment defaults
ENV PORT=8080
ENV MAX_BODY_BYTES=2000000
ENV RENDER_TIMEOUT_MS=30000
ENV WKHTMLTOPDF_PATH=wkhtmltopdf
ENV ALLOW_NET=false
ENV JSON_LOGS=false

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Run the binary
CMD ["/app/trykkeri-api"]
