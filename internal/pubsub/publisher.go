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

func PublishFromReader(ctx context.Context, client *pubsub.Client, projectID string, r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	var wg sync.WaitGroup
	var failures int
	var mu sync.Mutex

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var entry PublishEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			slog.Error("Failed to unmarshal entry", "line", lineNum, "error", err)
			failures++
			continue
		}

		topicID := entry.Topic
		if topicID == "" {
			topicID = entry.TopicID
		}
		if topicID == "" {
			topicID = entry.TargetTopic
		}

		if topicID == "" {
			slog.Error("Missing topic", "line", lineNum)
			failures++
			continue
		}

		payload := entry.Payload
		if payload == nil {
			payload = entry.Data
		}
		if payload == nil {
			slog.Error("Missing payload/data", "line", lineNum)
			failures++
			continue
		}

		data, err := json.Marshal(payload)
		if err != nil {
			slog.Error("Failed to marshal payload", "line", lineNum, "error", err)
			failures++
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
			slog.Error("Failed to ensure topic", "line", lineNum, "topicID", topicID, "error", err)
			failures++
			continue
		}

		wg.Add(1)
		go func(p *pubsub.Publisher, d []byte, oKey string, a map[string]string, lNum int, tID string) {
			defer wg.Done()
			res := p.Publish(ctx, &pubsub.Message{
				Data:        d,
				Attributes:  a,
				OrderingKey: oKey,
			})
			_, err := res.Get(ctx)
			if err != nil {
				slog.Error("Publish failed", "line", lNum, "topicID", tID, "error", err)
				mu.Lock()
				failures++
				mu.Unlock()
			}
		}(publisher, data, entry.OrderingKey, attrs, lineNum, topicID)
	}

	if err := scanner.Err(); err != nil {
		return failures, err
	}

	wg.Wait()
	return failures, nil
}
