# Project Objective: Pub/Sub Local
The objective of this project is to provide a comprehensive, standalone local development environment for Google Cloud Pub/Sub.

# Modern Go Guidelines
This project is using Go 1.22.8. Stick to modern Go best practices and freely use language features up to and including this version.

- **Skill Reference**: Apply guidelines from the `@use-modern-go` skill.
- **Workflow**: Always run `go fix ./...` to ensure code remains compliant with modern standards.
- **Logging**: Exclusively use **`log/slog`** for all structured logging.
- **Routing**: Use enhanced `http.ServeMux` patterns (e.g., `"GET /api/{id}"`).
