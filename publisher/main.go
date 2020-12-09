package main

import (
	"context"
	"encoding/json"
	"flag"
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
	topicID := os.Getenv("PUBSUB_TOPIC_ID")

	data := Data{EventType: "create"}
	flag.StringVar(&data.EventID, "id", "", "")
	flag.Parse()

	// Creates a client.
	// PUBSUB_EMULATOR_HOSTの環境変数が定義されていると、Client作成時にそのホストに対し繋いでくれる
	// https://github.com/googleapis/google-cloud-go/blob/1bc65d9124ba22db5bec4c71b6378c27dfc04724/pubsub/pubsub.go#L59-L65
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}

	topic := client.Topic(topicID)
	msg, err := createMessage(data)
	if err != nil {
		fmt.Printf("failed to create message: %s\n", err)
		return
	}
	result := topic.Publish(ctx, msg)
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("failed to publish: %s\n", err)
		return
	}
	fmt.Printf("message created.\n")
}

func createMessage(data Data) (*pubsub.Message, error) {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("failed to marshal.\n", data)
		return nil, err
	}

	return &pubsub.Message{Data: b}, nil
}
