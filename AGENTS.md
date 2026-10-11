# Project Objective: Pub/Sub Local

The objective of this project is to provide a comprehensive, standalone local development environment for Google Cloud Pub/Sub.

## Core Capabilities
- **Local Emulation**: Runs the Google Cloud Pub/Sub emulator in Docker.
- **Automated Setup**: Automatically configures topics and subscriptions based on `settings.yaml`.
- **Auto-Discovery**: Topics created at runtime get a `<topic>-echo` capture subscription automatically (`auto_discover`, `discover_interval` in `settings.yaml`).
- **Message Persistence**: A Go-based listener service captures all messages published to local topics and saves them as JSON files.
- **Web Visualization**: A lightweight Go + HTMX web viewer for browsing captured messages.
- **Task Orchestration**: Uses `mise` to manage the stack.

---

# Modern Go Development Guidelines

This project strictly adheres to modern Go best practices for **Go 1.27**.

### 1. Skill Reference
Refer to the `@use-modern-go` skill for specific syntax and pattern guidelines.

### 2. Standard Workflow
- **Code Maintenance**: Consistently use `go fix ./...` to update code to modern patterns.
- **Logging**: Use **`log/slog`** for all structured logging.
- **Web**: Utilize enhanced `http.ServeMux` patterns and `r.PathValue` where applicable.
- **Concurrency**: Use `sync.WaitGroup.Go`.

### 3. Layout
- `cmd/listener`, `cmd/web`, `cmd/publish`: binaries.
- `internal/pubsub`, `internal/config`, `pkg/models`: shared code.
- `docs/`: example settings and data; keep `README.md` in sync when behavior or settings change.
