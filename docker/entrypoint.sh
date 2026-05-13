#!/bin/bash

# Set default project ID if not provided
PROJECT_ID=${PUBSUB_PROJECT_ID:-test-project}

gcloud beta emulators pubsub start --project="$PROJECT_ID" --host-port=0.0.0.0:8681 &

# Wait for emulator to start
sleep 1

./pubsub-listener -settings "${SETTINGS_PATH:-settings.yaml}" &

./pubsub-web &

# Wait for any process to exit
wait -n

# Exit with the status of the first process to exit
exit $?
