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

The fastest way to get everything running (Emulator, Listener, and Viewer) is using **mise**:

```bash
mise run up
```

This starts the entire stack. You can then access the **Web Viewer** at:
👉 **http://localhost:8682**

---

### Features at a glance:
- **Automatic Setup**: Your topics and subscriptions are ready the moment you start.
- **Instant Search**: Browse messages with localized timestamps and keyboard navigation.
- **Dark Mode Support**: Easy on the eyes for long development sessions.
- **Data Persistence**: Messages are saved as plain JSON files for easy debugging.
