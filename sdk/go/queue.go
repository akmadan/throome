package throome

import (
	"context"
	"fmt"
)

// QueueClient provides queue/message broker operations
type QueueClient struct {
	clusterClient *ClusterClient
}

// Publish publishes a message to a topic
func (q *QueueClient) Publish(ctx context.Context, topic string, message []byte) error {
	req := QueuePublishRequest{
		Topic:   topic,
		Message: message,
	}

	path := fmt.Sprintf("/api/v1/clusters/%s/queue/publish", q.clusterClient.clusterID)
	return q.clusterClient.client.request(ctx, "POST", path, req, nil)
}

// Subscribe subscribes to a topic (placeholder - requires WebSocket/long-polling implementation)
func (q *QueueClient) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	return fmt.Errorf("subscribe not yet implemented in SDK - use direct Kafka consumer")
}
