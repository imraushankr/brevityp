# ====================== BUILDER STAGE ======================
FROM golang:1.24-alpine AS builder

# Install build tools and Task
RUN apk add --no-cache git bash make \
    && sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

WORKDIR /app
COPY . .

# Install dependencies
RUN go mod download

# Build the application
RUN task build

# ====================== DEVELOPMENT STAGE ======================
FROM builder AS dev

# Install live-reload tools
RUN go install github.com/cosmtrek/air@latest \
    && go install github.com/cespare/reflex@latest

# Create required directories
RUN mkdir -p /app/data /app/logs /app/uploads

EXPOSE 8080
CMD ["task", "dev"]

# ====================== PRODUCTION STAGE ======================
FROM alpine:latest AS prod

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache bash sqlite

# Copy built binary and assets
COPY --from=builder /app/brevity .
COPY --from=builder /app/src/migrations ./migrations
COPY --from=builder /app/.env.example .env

# Create data directories
RUN mkdir -p /app/data /app/logs /app/uploads

EXPOSE 8080
CMD ["/app/brevity"]