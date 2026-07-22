# Pub/Sub Local

## Instant Local Pub/Sub Development

Develop and test your Google Cloud Pub/Sub applications without a cloud account or internet connection. This tool provides a fully functional local environment that mimics production behavior while giving you complete visibility into your message flows.

### Why use this?
- **Zero Configuration**: No service accounts or GCP projects to manage.
- **Visual Inspection**: See exactly what's happening with a clean, real-time web interface.
- **Reliable Testing**: Capture every message to disk so you can verify payloads and ordering.
- **Fast Feedback**: Instantly clear messages and restart your flow in seconds.
- **Isolated Environment**: Keep your local development entirely decoupled from production or staging.

## Quick Start

Add this to your project's `docker-compose.yml` to start the entire stack (Emulator, Listener, and Web Viewer):

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

Access the **Web Viewer** at:
👉 **http://localhost:8682**

To connect your application to the local Pub/Sub emulator, set:

```bash
export PUBSUB_EMULATOR_HOST=localhost:8681
```

---

### Features at a glance:
- **Automatic Setup**: Your topics and subscriptions are ready the moment you start.
- **Instant Search**: Browse messages with localized timestamps and keyboard navigation.
- **Dark Mode Support**: Easy on the eyes for long development sessions.
- **Data Persistence**: Messages are saved as plain JSON files for easy debugging.

---

## Examples

You can find example configuration and data files in the [docs](./docs) directory:
- [settings.example.yaml](./docs/settings.example.yaml): Template for configuring topics and subscriptions.
- [example.jsonl](./docs/example.jsonl): Sample data for batch publishing.

### Publish messages via API (curl)

The web service exposes:

```text
POST http://localhost:8682/api/publish
```

Send JSONL (recommended) using the sample file:

```bash
curl -X POST "http://localhost:8682/api/publish" \
  -H "Content-Type: application/jsonl" \
  --data-binary @docs/example.jsonl
```

Send JSON array payload (`[]PublishEntry`) with `application/json`:

```bash
curl -X POST "http://localhost:8682/api/publish" \
  -H "Content-Type: application/json" \
  --data-raw '[{"topic":"topic-1","payload":{"id":10,"message":"Hello API"},"attributes":{"key":"value"}}]'
```

Notes:
- For `Content-Type: application/jsonl`, each line must be one complete JSON record.
- For `Content-Type: application/json`, the body must be a JSON array of publish entries.
- To send a single entry with `application/json`, wrap it in an array with one item.

---

## Local Development

If you have cloned this repository and want to run it locally for development:

### Using mise
```bash
mise run up
```

### Using local Docker Compose
```bash
docker compose up -d
```

---

## License

This project is licensed under the [MIT License](LICENSE).
