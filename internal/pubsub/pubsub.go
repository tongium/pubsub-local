package pubsub

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func EnsureTopic(ctx context.Context, client *pubsub.Client, projectID, topicID string) (*pubsub.Publisher, error) {
	topicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)

	_, err := client.TopicAdminClient.GetTopic(ctx, &pubsubpb.GetTopicRequest{
		Topic: topicName,
	})
	if err == nil {
		return client.Publisher(topicID), nil
	}

	if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
		return nil, err
	}

	_, err = client.TopicAdminClient.CreateTopic(ctx, &pubsubpb.Topic{
		Name: topicName,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return client.Publisher(topicID), nil
		}
		return nil, err
	}

	slog.Info("Topic created", "topicID", topicID)
	return client.Publisher(topicID), nil
}

func EnsureSubscription(ctx context.Context, client *pubsub.Client, projectID, subID, topicID string) (*pubsub.Subscriber, error) {
	subName := fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
	topicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)

	_, err := client.SubscriptionAdminClient.GetSubscription(ctx, &pubsubpb.GetSubscriptionRequest{
		Subscription: subName,
	})
	if err == nil {
		return client.Subscriber(subID), nil
	}

	if s, ok := status.FromError(err); !ok || s.Code() != codes.NotFound {
		return nil, err
	}

	_, err = client.SubscriptionAdminClient.CreateSubscription(ctx, &pubsubpb.Subscription{
		Name:  subName,
		Topic: topicName,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.AlreadyExists {
			return client.Subscriber(subID), nil
		}
		return nil, err
	}

	slog.Info("Subscription created", "subscriptionID", subID)
	return client.Subscriber(subID), nil
}

// ListTopicIDs returns the IDs of all topics currently in the project.
func ListTopicIDs(ctx context.Context, client *pubsub.Client, projectID string) ([]string, error) {
	it := client.TopicAdminClient.ListTopics(ctx, &pubsubpb.ListTopicsRequest{
		Project: "projects/" + projectID,
	})

	prefix := fmt.Sprintf("projects/%s/topics/", projectID)
	var ids []string
	for {
		t, err := it.Next()
		if errors.Is(err, iterator.Done) {
			return ids, nil
		}
		if err != nil {
			return nil, err
		}
		ids = append(ids, strings.TrimPrefix(t.GetName(), prefix))
	}
}
