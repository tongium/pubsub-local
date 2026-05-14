# Pub/Sub local

## Description

A local Google Cloud Pub/Sub emulator setup that automatically configures topics from YAML configuration files and saves incoming messages to the file system. This project enables local development and testing of Pub/Sub-based applications without requiring a Google Cloud account.

**Key Features:**
- **Topic Configuration**: Automatically creates topics and subscriptions via YAML files.
- **Message Persistence**: Saves incoming messages to disk as JSON for local inspection.
- **HTMX Web Viewer**: Modern, minimalist viewer (Muji-inspired) for browsing captured messages.
- **Advanced Viewer Features**:
  - **Syntax Highlighting**: Beautifully formatted JSON payloads using Prism.js.
  - **Dark Mode**: System-aware theme switching.
  - **Collapsible Topics**: Organizable sidebar with persistent state.
  - **Time-Based Sorting**: Messages sorted by publication time (newest first).
  - **Local Timestamps**: Human-readable message timing in your local timezone.
  - **Keyboard Navigation**: Flip through messages quickly with Arrow keys.
- **Batch Publishing**: CLI tool to publish multiple messages from JSONL files.
- **Docker Ready**: Pre-configured for both local development (mise) and CI/CD (GitHub Actions).

## Quick Start

### Using mise (Local Development)

Run the emulator, listener, and web viewer together:

```bash
mise run up
```

This starts the Pub/Sub emulator in Docker, the Go listener, and the web viewer (at http://localhost:8682) in parallel.

### Using Docker Compose (Portable)

To run as a standalone service in your project:

```yaml
# docker-compose.yml
services:
  pubsub-local:
    image: ghcr.io/tongium/pubsub-local
    environment:
      PUBSUB_PROJECT_ID: test-project
    ports:
      - "8681:8681"
      - "8682:8682"
    volumes:
      - ./pubsub/settings.yaml:/app/settings.yaml
      - ./pubsub/messages:/app/messages
```

Run it:
```bash
docker compose up -d
```
## Docker Configuration

### Build the Image
```bash
docker build -t pubsub-local .
```

### Manual Run
```bash
docker run -d \
  -p 8681:8681 \
  -p 8682:8682 \
  -v $(pwd)/settings.yaml:/app/settings.yaml \
  -v $(pwd)/messages:/app/messages \
  pubsub-local
```

## Manual Setup (Local)

If you prefer to run them separately:

**Emulator (Docker):**

```bash
mise run emulator
```

**Listener:**

```bash
mise run listener
```

**Web Viewer:**

```bash
mise run view
```

## Publish (JSONL)

Use the `publish` task to batch publish messages:

```bash
# From a file
mise run publish -- --jsonl ./data/example.jsonl

# From stdin
cat ./data/example.jsonl | mise run publish
```