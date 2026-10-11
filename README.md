# Pub/Sub Local

Local Google Cloud Pub/Sub for development: emulator, message capture, and a web viewer. No GCP account needed.

## Quick Start

`docker-compose.yml`:

```yaml
services:
  pubsub-local:
    image: ghcr.io/tongium/pubsub-local
    environment:
      PUBSUB_PROJECT_ID: test-project
    ports:
      - "8681:8681" # Emulator
      - "8682:8682" # Web Viewer
    volumes:
      - ./pubsub/settings.yaml:/app/settings.yaml
      - ./pubsub/messages:/app/messages
```

`pubsub/settings.yaml`:

```yaml
topics:
  - id: topic-1
    subscribers:
      - id: topic-1-sub
```

Point your app at the emulator:

```bash
export PUBSUB_EMULATOR_HOST=localhost:8681
export PUBSUB_PROJECT_ID=test-project
```

Open the viewer: **http://localhost:8682**

## How it works

Every topic gets a `<topic>-echo` subscription; each message is saved as JSON in `/app/messages` and shown in the viewer.

Topics your app creates later are picked up automatically (no config needed). Create the topic before publishing, or earlier messages are missed.

## Settings

| Key | Default | Description |
| --- | --- | --- |
| `auto_discover` | `true` | Capture topics created at runtime |
| `discover_interval` | `2s` | How often to look for new topics |
| `topics[].id` | | Topic created at startup |
| `topics[].subscribers[].id` | | Extra subscription on that topic |

See [settings.example.yaml](./docs/settings.example.yaml).

## Publish test messages

```bash
# JSONL: one entry per line
curl -X POST localhost:8682/api/publish \
  -H "Content-Type: application/jsonl" \
  --data-binary @docs/example.jsonl

# JSON array
curl -X POST localhost:8682/api/publish \
  -H "Content-Type: application/json" \
  -d '[{"topic":"topic-1","payload":{"id":1},"attributes":{"key":"value"}}]'
```

Entry fields: `topic`, `payload`, optional `attributes`, `ordering_key`.

## Development

Requires Go 1.27+.

```bash
cp docs/settings.example.yaml settings.yaml
mise run dev
```

Pushing a `v*.*.*` tag publishes the image to `ghcr.io/tongium/pubsub-local`.

## License

[MIT](LICENSE)
