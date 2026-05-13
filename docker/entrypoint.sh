#!/bin/bash

# Set default project ID if not provided
PROJECT_ID=${PUBSUB_PROJECT_ID:-test-project}

slog_info() {
    echo "{\"time\":\"$(date -u +%Y-%m-%dT%H:%M:%SZ)\",\"level\":\"INFO\",\"msg\":\"$1\"}"
}

slog_info "Starting Pub/Sub Emulator on port 8681..."
gcloud beta emulators pubsub start --project="$PROJECT_ID" --host-port=0.0.0.0:8681 &

# Wait for emulator to start (simple sleep)
sleep 2

slog_info "Starting Pub/Sub Listener..."
./pubsub-listener -settings "${SETTINGS_PATH:-settings.yaml}" &

slog_info "Starting Web Viewer on port 8080..."
./pubsub-web &

# Wait for any process to exit
wait -n

# Exit with the status of the first process to exit
exit $?
