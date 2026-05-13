package pubsub

import (
	"context"
	"fmt"
	"log/slog"

	"cloud.google.com/go/pubsub/v2"
	"cloud.google.com/go/pubsub/v2/apiv1/pubsubpb"
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
