# Project Objective: Pub/Sub Local
A standalone local development environment for Google Cloud Pub/Sub: emulator, message-capturing listener (with topic auto-discovery), and a Go + HTMX web viewer. See `AGENTS.md` for the full guidelines.

# Modern Go Guidelines
This project uses Go 1.27. Stick to modern Go best practices and freely use language features up to and including this version.

- **Skill Reference**: Apply guidelines from the `@use-modern-go` skill.
- **Workflow**: Always run `go fix ./...` to ensure code remains compliant with modern standards.
- **Logging**: Exclusively use **`log/slog`** for all structured logging.
- **Routing**: Use enhanced `http.ServeMux` patterns (e.g., `"GET /api/{id}"`).
- **Concurrency**: Use `sync.WaitGroup.Go` instead of manual `Add`/`Done`.
