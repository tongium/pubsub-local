# Pub/Sub local

## Description

A local Google Cloud Pub/Sub emulator setup that automatically configures topics from YAML configuration files and saves incoming messages to the file system. This project enables local development and testing of Pub/Sub-based applications without requiring a Google Cloud account.

**Key Features:**
- Topic configuration via YAML files
- Automatic message persistence to disk
- Batch message publishing from JSONL format

## Quick Start

### Using mise (Local Development)

Run the emulator, listener, and web viewer together:

```bash
mise run up
```

This starts the Pub/Sub emulator in Docker, the Go listener, and the web viewer (at http://localhost:8080) in parallel.

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
      - "8080:8080"
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
  -p 8080:8080 \
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