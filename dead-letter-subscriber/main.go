package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

type Data struct {
	EventType string `json:"eventType"`
	EventID   string `json:"eventId"`
}

func main() {
	ctx := context.Background()

	projectID := os.Getenv("PUBSUB_PROJECT_ID")
	subscriptionID := os.Getenv("PUBSUB_DEAD_SUBSCRIPTION_ID")

	// Creates a client.
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	subscription := client.Subscription(subscriptionID)
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("got message: %s %q\n", msg.ID, string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		log.Fatalf("failed to messge: %v", err)
	}
}
