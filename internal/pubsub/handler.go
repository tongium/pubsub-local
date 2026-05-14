package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/pubsub/v2"
	"github.com/tongium/pubsub-local/pkg/models"
)

func CreateHandler(topicID string) func(context.Context, *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		slog.Info("Received message", "topicID", topicID, "messageID", msg.ID)
		if err := saveMessageToFile(topicID, msg); err != nil {
			slog.Error("Error saving message", "topicID", topicID, "messageID", msg.ID, "error", err)
		}
		msg.Ack()
	}
}

func saveMessageToFile(topicID string, msg *pubsub.Message) error {
	dir := filepath.Join("messages", topicID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var data any
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		// If not JSON, save as string
		data = string(msg.Data)
	}

	record := models.MessageRecord{
		PublishTime: msg.PublishTime.Format(time.RFC3339),
		MessageID:   msg.ID,
		Attributes:  msg.Attributes,
		OrderingKey: msg.OrderingKey,
		Data:        data,
	}

	filename := filepath.Join(dir, fmt.Sprintf("%s.json", msg.ID))
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(record)
}
