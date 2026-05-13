package main

import (
	"cmp"
	"context"
	"flag"
	"log/slog"
	"os"

	"cloud.google.com/go/pubsub/v2"
	"github.com/joho/godotenv"
	internalpubsub "github.com/tongium/pubsub-local/internal/pubsub"
)

func main() {
	_ = godotenv.Load()

	projectID := cmp.Or(os.Getenv("PUBSUB_PROJECT_ID"), "test-project")

	jsonlPath := flag.String("jsonl", "", "Path to JSONL file (defaults to stdin)")
	flag.Parse()

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	var input *os.File
	if *jsonlPath != "" {
		f, err := os.Open(*jsonlPath)
		if err != nil {
			slog.Error("Failed to open file", "path", *jsonlPath, "error", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	} else {
		input = os.Stdin
	}

	failures, err := internalpubsub.PublishFromReader(ctx, client, projectID, input)
	if err != nil {
		slog.Error("Error during publishing", "error", err)
		os.Exit(1)
	}

	if failures > 0 {
		slog.Warn("Publishing completed with failures", "count", failures)
		os.Exit(1)
	}

	slog.Info("Publish complete")
}
