#!/bin/bash

# Set default project ID if not provided
PROJECT_ID=${PUBSUB_PROJECT_ID:-test-project}

gcloud beta emulators pubsub start --project="$PROJECT_ID" --host-port=0.0.0.0:8681 &

# Wait for emulator to start (using bash builtin /dev/tcp to avoid external dependencies)
echo "Waiting for Pub/Sub emulator to start on localhost:8681..."
MAX_WAIT=30
WAIT_COUNT=0
while ! (echo > /dev/tcp/localhost/8681) >/dev/null 2>&1; do
    sleep 0.1
    WAIT_COUNT=$((WAIT_COUNT + 1))
    if [ $WAIT_COUNT -ge $((MAX_WAIT * 10)) ]; then
        echo "Error: Pub/Sub emulator failed to start after ${MAX_WAIT}s"
        exit 1
    fi
done
echo "Pub/Sub emulator is ready!"

./pubsub-listener -settings "${SETTINGS_PATH:-settings.yaml}" &

./pubsub-web &

# Wait for any process to exit
wait -n

# Exit with the status of the first process to exit
exit $?
