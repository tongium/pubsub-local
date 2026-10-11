package pubsub

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub/v2"
)

type PublishEntry struct {
	Topic       string            `json:"topic"`
	TopicID     string            `json:"topic_id"`
	TargetTopic string            `json:"target_topic"`
	Payload     any               `json:"payload"`
	Data        any               `json:"data"`
	Attributes  map[string]string `json:"attributes"`
	Headers     map[string]string `json:"headers"`
	OrderingKey string            `json:"ordering_key"`
}

type publishInput struct {
	entry      PublishEntry
	entryIndex int
}

func PublishFromReader(ctx context.Context, client *pubsub.Client, projectID string, r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	var failures int

	var inputs []publishInput
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var entry PublishEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			slog.Error("Failed to decode publish entry", "index", lineNum, "error", err)
			failures++
			continue
		}
		inputs = append(inputs, publishInput{entry: entry, entryIndex: lineNum})
	}

	if err := scanner.Err(); err != nil {
		return failures, err
	}

	publishFailures, err := publishInputs(ctx, client, projectID, inputs)
	if err != nil {
		return failures + publishFailures, err
	}

	return failures + publishFailures, nil
}

func PublishEntries(ctx context.Context, client *pubsub.Client, projectID string, entries []PublishEntry) (int, error) {
	inputs := make([]publishInput, 0, len(entries))
	for idx, entry := range entries {
		inputs = append(inputs, publishInput{entry: entry, entryIndex: idx + 1})
	}
	return publishInputs(ctx, client, projectID, inputs)
}

func publishInputs(ctx context.Context, client *pubsub.Client, projectID string, inputs []publishInput) (int, error) {
	var wg sync.WaitGroup
	var failures int
	var mu sync.Mutex
	addFailure := func() {
		mu.Lock()
		failures++
		mu.Unlock()
	}

	for _, input := range inputs {
		entry := input.entry
		entryIndex := input.entryIndex

		topicID := entry.Topic
		if topicID == "" {
			topicID = entry.TopicID
		}
		if topicID == "" {
			topicID = entry.TargetTopic
		}
		if topicID == "" {
			slog.Error("Missing topic", "index", entryIndex)
			addFailure()
			continue
		}

		payload := entry.Payload
		if payload == nil {
			payload = entry.Data
		}
		if payload == nil {
			slog.Error("Missing payload/data", "index", entryIndex)
			addFailure()
			continue
		}

		data, err := json.Marshal(payload)
		if err != nil {
			slog.Error("Failed to marshal payload", "index", entryIndex, "error", err)
			addFailure()
			continue
		}

		attrs := entry.Attributes
		if attrs == nil {
			attrs = entry.Headers
		}
		if attrs == nil {
			attrs = make(map[string]string)
		}

		publisher, err := EnsureTopic(ctx, client, projectID, topicID)
		if err != nil {
			slog.Error("Failed to ensure topic", "index", entryIndex, "topicID", topicID, "error", err)
			addFailure()
			continue
		}

		wg.Go(func() {
			res := publisher.Publish(ctx, &pubsub.Message{
				Data:        data,
				Attributes:  attrs,
				OrderingKey: entry.OrderingKey,
			})
			if _, err := res.Get(ctx); err != nil {
				slog.Error("Publish failed", "index", entryIndex, "topicID", topicID, "error", err)
				addFailure()
			}
		})
	}

	wg.Wait()
	return failures, nil
}
