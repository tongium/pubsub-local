# Pub/Sub local

## Description

A local Google Cloud Pub/Sub emulator setup that automatically configures topics from YAML configuration files and saves incoming messages to the file system. This project enables local development and testing of Pub/Sub-based applications without requiring a Google Cloud account.

**Key Features:**
- Topic configuration via YAML files
- Automatic message persistence to disk
- Local emulator setup with tmux multi-pane execution
- Batch message publishing from JSONL format

## Quick Start

Run the emulator and subscriber together:

```bash
chmod +x run.sh
./run.sh
```

This script starts the Pub/Sub emulator and runs the subscriber (`main.py`) with `uv`.

## Manual Setup

If you prefer to run the emulator and subscriber separately:

**Terminal 1 - Start emulator:**

```bash
gcloud beta emulators pubsub start --project=test-project --host-port=0.0.0.0:8681
```

**Terminal 2 - Set environment and run subscriber:**

```bash
export PUBSUB_EMULATOR_HOST=localhost:8681
export PUBSUB_PROJECT_ID=test-project
uv run main.py
```

## Subscribe

The subscriber creates topics from [settings.yaml](settings.yaml) and echoes incoming messages
to stdout and the [messages](messages) folder as JSON files.

## Publish (JSONL)

Use [publish.py](publish.py) to batch publish messages from JSONL (one JSON object
per line). Each line must include a target topic and payload.

Supported fields per line:
- `topic`, `topic_id`, or `target_topic` (string)
- `payload` or `data` (any JSON value)
- `headers` or `attributes` (object of key/value pairs, optional)
- `ordering_key` (string, optional)

Example JSONL from [data/example.jsonl](data/example.jsonl):

```json
{"topic":"test1","payload":{"id":1,"event":"test1_event"},"headers":{"source":"example","timestamp":"2026-02-14T00:00:00Z"}}
{"topic":"test2","payload":{"id":2,"event":"test2_event"},"headers":{"source":"example","timestamp":"2026-02-14T00:00:00Z"}}
{"topic":"test3","payload":{"id":3,"event":"test3_event"},"headers":{"source":"example","timestamp":"2026-02-14T00:00:00Z"}}
```

Publish from a file (with emulator running):

```bash
export PUBSUB_EMULATOR_HOST=localhost:8681
export PUBSUB_PROJECT_ID=test-project
uv run publish.py --jsonl ./data/example.jsonl
```

Or from stdin:

```bash
cat ./data/example.jsonl | uv run publish.py
```