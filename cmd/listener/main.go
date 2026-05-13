package main

import (
	"cmp"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cloud.google.com/go/pubsub/v2"
	"github.com/joho/godotenv"
	"github.com/tongium/pubsub-local/internal/config"
	internalpubsub "github.com/tongium/pubsub-local/internal/pubsub"
)

func main() {
	_ = godotenv.Load()

	projectID := cmp.Or(os.Getenv("PUBSUB_PROJECT_ID"), "test-project")

	settingsPath := flag.String("settings", "settings.yaml", "Path to configuration file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		slog.Error("Failed to create client", "error", err)
		os.Exit(1)
	}
	defer client.Close()

	settings, err := config.LoadSettings(*settingsPath)
	if err != nil {
		slog.Error("Failed to load settings", "path", *settingsPath, "error", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, t := range settings.Topics {
		_, err := internalpubsub.EnsureTopic(ctx, client, projectID, t.ID)
		if err != nil {
			slog.Error("Failed to ensure topic", "topicID", t.ID, "error", err)
			continue
		}

		// Create echo subscription
		echoSubID := fmt.Sprintf("%s-echo", t.ID)
		sub, err := internalpubsub.EnsureSubscription(ctx, client, projectID, echoSubID, t.ID)
		if err != nil {
			slog.Error("Failed to ensure subscription", "subscriptionID", echoSubID, "error", err)
			continue
		}

		wg.Add(1)
		go func(topicID string, s *pubsub.Subscriber) {
			defer wg.Done()
			handler := internalpubsub.CreateHandler(topicID)
			err := s.Receive(ctx, handler)
			if err != nil {
				slog.Error("Subscriber receiver error", "subscriberID", s.ID(), "error", err)
			}
		}(t.ID, sub)

		// Create other subscriptions
		for _, s := range t.Subscribers {
			_, err := internalpubsub.EnsureSubscription(ctx, client, projectID, s.ID, t.ID)
			if err != nil {
				slog.Error("Failed to ensure subscription", "subscriptionID", s.ID, "error", err)
			}
		}
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutting down...")
		cancel()
	}()

	slog.Info("Listening for messages...")
	wg.Wait()
}
