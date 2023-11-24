package content

import (
	"context"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	pubsubClient *pubsub.Client
	topic        *pubsub.Topic
	onceClient   sync.Once
)

func initializePubsub() error {
	var err error
	onceClient.Do(func() {
		ctx := context.Background()
		pubsubClient, err = pubsub.NewClient(ctx, "takeoff-task-3")
		if err != nil {
			// Handle error creating Pub/Sub client.
			log.Fatalf("Error creating Pub/Sub client: %v", err)
			return
		}
		// Replace "your-topic-name" with the actual name of your Pub/Sub topic.
		topic = pubsubClient.Topic("create-employee-topic")
	})
	return err
}
