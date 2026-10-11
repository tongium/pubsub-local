package main

import (
	"cmp"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		tracked = make(map[string]bool)
	)

	// capture creates the "<topic>-echo" subscription and starts persisting its messages.
	capture := func(topicID string) {
		mu.Lock()
		defer mu.Unlock()
		if tracked[topicID] {
			return
		}

		echoSubID := topicID + "-echo"
		sub, err := internalpubsub.EnsureSubscription(ctx, client, projectID, echoSubID, topicID)
		if err != nil {
			slog.Error("Failed to ensure subscription", "subscriptionID", echoSubID, "error", err)
			return
		}
		tracked[topicID] = true

		wg.Go(func() {
			handler := internalpubsub.CreateHandler(topicID)
			if err := sub.Receive(ctx, handler); err != nil {
				slog.Error("Subscriber receiver error", "subscriberID", sub.ID(), "error", err)
			}
		})
	}

	for _, t := range settings.Topics {
		if _, err := internalpubsub.EnsureTopic(ctx, client, projectID, t.ID); err != nil {
			slog.Error("Failed to ensure topic", "topicID", t.ID, "error", err)
			continue
		}

		capture(t.ID)

		for _, s := range t.Subscribers {
			if _, err := internalpubsub.EnsureSubscription(ctx, client, projectID, s.ID, t.ID); err != nil {
				slog.Error("Failed to ensure subscription", "subscriptionID", s.ID, "error", err)
			}
		}
	}

	if settings.AutoDiscoverEnabled() {
		interval := 2 * time.Second
		if settings.DiscoverInterval != "" {
			d, err := time.ParseDuration(settings.DiscoverInterval)
			if err != nil || d <= 0 {
				slog.Warn("Invalid discover_interval, using default", "value", settings.DiscoverInterval)
			} else {
				interval = d
			}
		}

		wg.Go(func() {
			slog.Info("Auto-discovering topics", "interval", interval)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				ids, err := internalpubsub.ListTopicIDs(ctx, client, projectID)
				if err != nil && ctx.Err() == nil {
					slog.Error("Failed to list topics", "error", err)
				}
				for _, id := range ids {
					capture(id)
				}
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		})
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
