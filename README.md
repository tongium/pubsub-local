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

### 1. Add the service to your `docker-compose.yml`

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

### 2. Create `pubsub/settings.yaml`

```yaml
auto_discover: true   # capture topics your app creates at runtime (default)
topics:
  - id: topic-1
    subscribers:
      - id: topic-1-sub
  - id: topic-2
```

See [settings.example.yaml](./docs/settings.example.yaml) for all options.

### 3. Point your app at the emulator

```bash
export PUBSUB_EMULATOR_HOST=localhost:8681
export PUBSUB_PROJECT_ID=test-project
```

### 4. Open the viewer

👉 **http://localhost:8682**

---

## How messages are captured

For every topic, Pub/Sub Local creates a subscription named `<topic-id>-echo` and
saves each message it receives as a JSON file under `/app/messages`. The web viewer reads those files.

- **Configured topics** (in `settings.yaml`) are created at startup.
- **Auto-discovery**: topics created later by your own application are detected
  (every `discover_interval`, default 2s) and get an `-echo` subscription automatically,
  so you can view their messages without editing any config. Messages published
  *before* the subscription exists are not captured, so create your topic before publishing
  (the usual case when your app calls `CreateTopic` on startup).
- Set `auto_discover: false` to capture only the topics listed in `settings.yaml`.

## Settings reference

| Key | Default | Description |
| --- | --- | --- |
| `auto_discover` | `true` | Auto-create an `-echo` capture subscription for topics created at runtime. |
| `discover_interval` | `2s` | How often to scan for new topics (Go duration). |
| `topics[].id` | – | Topic to create at startup. |
| `topics[].subscribers[].id` | – | Extra subscription to create on that topic. |

Environment variables: `PUBSUB_PROJECT_ID` (default `test-project`), `PUBSUB_EMULATOR_HOST` (for clients), `SETTINGS_PATH` (default `settings.yaml`).

---

### Features at a glance:
- **Automatic Setup**: Your topics and subscriptions are ready the moment you start.
- **Auto-Discovery**: Topics your app creates on the fly get a capture subscription automatically.
- **Instant Search**: Browse messages with localized timestamps and keyboard navigation.
- **Dark Mode Support**: Easy on the eyes for long development sessions.
- **Data Persistence**: Messages are saved as plain JSON files for easy debugging.

---

## Publishing test messages

Example files live in [docs](./docs):
- [settings.example.yaml](./docs/settings.example.yaml): annotated settings template.
- [example.jsonl](./docs/example.jsonl): sample data for batch publishing.

### Via the HTTP API

```text
POST http://localhost:8682/api/publish
```

JSONL (one entry per line, `application/jsonl`):

```bash
curl -X POST "http://localhost:8682/api/publish" \
  -H "Content-Type: application/jsonl" \
  --data-binary @docs/example.jsonl
```

JSON array (`application/json`):

```bash
curl -X POST "http://localhost:8682/api/publish" \
  -H "Content-Type: application/json" \
  --data-raw '[{"topic":"topic-1","payload":{"id":10,"message":"Hello API"},"attributes":{"key":"value"}}]'
```

Each entry has `topic`, `payload`, and optional `attributes` and `ordering_key`.
A single entry sent as `application/json` must still be wrapped in an array.

---

## Local development

Clone the repo and run the stack from source (requires Go 1.27+):

```bash
mise run dev          # emulator (podman), listener and web viewer
docker compose up -d  # or: build and run the all-in-one image
```

Individual tasks: `mise run emulator | listener | view | publish`.
Copy `docs/settings.example.yaml` to `settings.yaml` first.

Releases: pushing a `v*.*.*` tag publishes the image to `ghcr.io/tongium/pubsub-local`.

---

## License

This project is licensed under the [MIT License](LICENSE).
