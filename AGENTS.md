# Pub/Sub Local

Local Google Cloud Pub/Sub environment: emulator in Docker, a Go listener that saves every message as JSON (and auto-creates `<topic>-echo` subscriptions for new topics), and a Go + HTMX web viewer. Managed with `mise`.

## Layout
- `cmd/listener`, `cmd/web`, `cmd/publish`: binaries
- `internal/`, `pkg/models`: shared code
- `docs/`: example settings and data. Update `README.md` when behavior or settings change.

## Go 1.27 guidelines
- Follow the `@use-modern-go` skill; run `go fix ./...`.
- Log with `log/slog`.
- Use enhanced `http.ServeMux` patterns and `r.PathValue`.
- Use `sync.WaitGroup.Go`.
