#!/bin/zsh

# Create a new tmux session
SESSION="pubsub"

# Kill existing session if it exists
tmux kill-session -t $SESSION 2>/dev/null

# Create new session
tmux new-session -d -s $SESSION -c "/Users/tongium/Workspaces/pubsub-local"

# Split window vertically (left and right panes)
tmux split-window -h -t $SESSION

# Run main.py in the left pane
tmux send-keys -t $SESSION:0.0 "gcloud beta emulators pubsub start --project=test-project --host-port=0.0.0.0:8681" Enter

# Run publish.py in the right pane
tmux send-keys -t $SESSION:0.1 "uv run main.py" Enter

# Attach to the tmux session
tmux attach-session -t $SESSION
