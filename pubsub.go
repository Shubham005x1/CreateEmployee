package content

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	pubsubClient *pubsub.Client
	oncePubsub   sync.Once
	topic        *pubsub.Topic
)

// initializePubsub initializes the Pub/Sub client and topic.
func initializePubsub() {
	oncePubsub.Do(func() {
		ctx := context.Background()
		// Initialize Pub/Sub client.
		pubsubClient, _ = pubsub.NewClient(ctx, "takeoff-task-3")
		// Replace "your-topic-name" with the actual name you want for your Pub/Sub topic.
		topic = pubsubClient.Topic("create-employee-topic")

		// Ensure the pubsubClient is closed after the function completes.
		// defer pubsubClient.Close()
	})
}
