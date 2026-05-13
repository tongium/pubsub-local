# Stage 1: Build Go binaries
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binaries
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o pubsub-listener cmd/listener/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o pubsub-publish cmd/publish/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o pubsub-web cmd/web/main.go

# Stage 2: Final image with Pub/Sub emulator
FROM gcr.io/google.com/cloudsdktool/google-cloud-cli:emulators

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/pubsub-listener .
COPY --from=builder /app/pubsub-publish .
COPY --from=builder /app/pubsub-web .

# Copy templates and default configuration
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/settings.exmaple.yaml ./settings.yaml

# Copy entrypoint script
COPY docker/entrypoint.sh .
RUN chmod +x entrypoint.sh

# Environment variables for internal communication
ENV PUBSUB_EMULATOR_HOST=localhost:8681
ENV PUBSUB_PROJECT_ID=test-project

# Expose Emulator (8681) and Web Viewer (8080)
EXPOSE 8681 8080

ENTRYPOINT ["./entrypoint.sh"]
