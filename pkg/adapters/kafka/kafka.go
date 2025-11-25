package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/akmadan/throome/pkg/adapters"
	"github.com/akmadan/throome/pkg/cluster"
)

// KafkaAdapter implements the QueueAdapter interface for Kafka
type KafkaAdapter struct {
	*adapters.BaseAdapter
	config    *cluster.ServiceConfig
	writer    *kafka.Writer
	readers   map[string]*kafka.Reader
	handlers  map[string]adapters.MessageHandler
	stopChans map[string]chan struct{}
}

// NewKafkaAdapter creates a new Kafka adapter
func NewKafkaAdapter(config *cluster.ServiceConfig) (adapters.Adapter, error) {
	adapter := &KafkaAdapter{
		BaseAdapter: adapters.NewBaseAdapter(config),
		config:      config,
		readers:     make(map[string]*kafka.Reader),
		handlers:    make(map[string]adapters.MessageHandler),
		stopChans:   make(map[string]chan struct{}),
	}
	return adapter, nil
}

// Connect establishes a connection to Kafka
func (k *KafkaAdapter) Connect(ctx context.Context) error {
	brokers := []string{fmt.Sprintf("%s:%d", k.config.Host, k.config.Port)}

	// Create a writer for publishing messages
	k.writer = &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		MaxAttempts:  3,
	}

	// Test connection by listing topics
	conn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", k.config.Host, k.config.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	k.SetConnected(true)
	return nil
}

// Disconnect closes all Kafka connections
func (k *KafkaAdapter) Disconnect(ctx context.Context) error {
	// Stop all consumers
	for topic, stopChan := range k.stopChans {
		close(stopChan)
		delete(k.stopChans, topic)
	}

	// Close all readers
	for topic, reader := range k.readers {
		_ = reader.Close() // Ignore errors during cleanup
		delete(k.readers, topic)
	}

	// Close writer
	if k.writer != nil {
		if err := k.writer.Close(); err != nil {
			return err
		}
	}

	k.SetConnected(false)
	return nil
}

// Ping checks if the Kafka connection is alive
func (k *KafkaAdapter) Ping(ctx context.Context) error {
	start := time.Now()

	conn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", k.config.Host, k.config.Port))
	duration := time.Since(start)

	if err != nil {
		k.RecordRequest(duration, false)
		k.LogActivity("PING", "PING", duration, err, "")
		return err
	}
	defer conn.Close()

	k.RecordRequest(duration, true)
	k.LogActivity("PING", "PING", duration, nil, "PONG")
	return nil
}

// HealthCheck performs a health check
func (k *KafkaAdapter) HealthCheck(ctx context.Context) (*adapters.HealthStatus, error) {
	start := time.Now()
	err := k.Ping(ctx)
	responseTime := time.Since(start)

	status := &adapters.HealthStatus{
		Healthy:      err == nil,
		ResponseTime: responseTime,
		LastChecked:  time.Now(),
	}

	if err != nil {
		status.ErrorMessage = err.Error()
	}

	return status, nil
}

// Publish publishes a message to a topic
func (k *KafkaAdapter) Publish(ctx context.Context, topic string, message []byte) error {
	start := time.Now()

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: message,
		Time:  time.Now(),
	})

	duration := time.Since(start)
	k.RecordRequest(duration, err == nil)

	// Log activity
	command := fmt.Sprintf("PUBLISH to topic '%s' (size: %d bytes)", topic, len(message))
	response := ""
	if err == nil {
		response = fmt.Sprintf("Message published successfully to topic '%s'", topic)
	}
	k.LogActivity("PUBLISH", command, duration, err, response)

	return err
}

// PublishWithKey publishes a message with a key to a topic
func (k *KafkaAdapter) PublishWithKey(ctx context.Context, topic string, key, message []byte) error {
	start := time.Now()

	err := k.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: message,
		Time:  time.Now(),
	})

	duration := time.Since(start)
	k.RecordRequest(duration, err == nil)

	// Log activity
	command := fmt.Sprintf("PUBLISH to topic '%s' with key '%s' (size: %d bytes)", topic, string(key), len(message))
	response := ""
	if err == nil {
		response = fmt.Sprintf("Message published successfully to topic '%s' with key", topic)
	}
	k.LogActivity("PUBLISH_WITH_KEY", command, duration, err, response)

	return err
}

// Subscribe subscribes to a topic
func (k *KafkaAdapter) Subscribe(ctx context.Context, topic string, handler adapters.MessageHandler) error {
	start := time.Now()

	if _, exists := k.readers[topic]; exists {
		err := fmt.Errorf("already subscribed to topic: %s", topic)
		k.LogActivity("SUBSCRIBE", fmt.Sprintf("SUBSCRIBE to topic '%s'", topic), time.Since(start), err, "")
		return err
	}

	brokers := []string{fmt.Sprintf("%s:%d", k.config.Host, k.config.Port)}

	// Create a reader for this topic
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        "throome-gateway", // Default group ID
		MinBytes:       10e3,              // 10KB
		MaxBytes:       10e6,              // 10MB
		CommitInterval: time.Second,
	})

	k.readers[topic] = reader
	k.handlers[topic] = handler

	// Start consuming messages in a goroutine
	stopChan := make(chan struct{})
	k.stopChans[topic] = stopChan

	go k.consumeMessages(ctx, topic, reader, handler, stopChan)

	duration := time.Since(start)
	command := fmt.Sprintf("SUBSCRIBE to topic '%s' with group 'throome-gateway'", topic)
	response := fmt.Sprintf("Successfully subscribed to topic '%s'", topic)
	k.LogActivity("SUBSCRIBE", command, duration, nil, response)

	return nil
}

// consumeMessages consumes messages from a topic
func (k *KafkaAdapter) consumeMessages(ctx context.Context, topic string, reader *kafka.Reader, handler adapters.MessageHandler, stopChan chan struct{}) {
	for {
		select {
		case <-stopChan:
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				// Handle error (log it)
				continue
			}

			// Convert to our Message type
			message := &adapters.Message{
				Topic:     msg.Topic,
				Key:       msg.Key,
				Value:     msg.Value,
				Timestamp: msg.Time,
				Offset:    msg.Offset,
				Headers:   make(map[string]string),
			}

			// Copy headers
			for _, header := range msg.Headers {
				message.Headers[header.Key] = string(header.Value)
			}

			// Call handler, ignore errors to continue processing
			_ = handler(ctx, message)
		}
	}
}

// Unsubscribe unsubscribes from a topic
func (k *KafkaAdapter) Unsubscribe(ctx context.Context, topic string) error {
	start := time.Now()

	// Stop the consumer
	if stopChan, exists := k.stopChans[topic]; exists {
		close(stopChan)
		delete(k.stopChans, topic)
	}

	// Close the reader
	if reader, exists := k.readers[topic]; exists {
		if err := reader.Close(); err != nil {
			k.LogActivity("UNSUBSCRIBE", fmt.Sprintf("UNSUBSCRIBE from topic '%s'", topic), time.Since(start), err, "")
			return err
		}
		delete(k.readers, topic)
	}

	// Remove handler
	delete(k.handlers, topic)

	duration := time.Since(start)
	command := fmt.Sprintf("UNSUBSCRIBE from topic '%s'", topic)
	response := fmt.Sprintf("Successfully unsubscribed from topic '%s'", topic)
	k.LogActivity("UNSUBSCRIBE", command, duration, nil, response)

	return nil
}

// CreateTopic creates a new topic
func (k *KafkaAdapter) CreateTopic(ctx context.Context, topic string, config map[string]interface{}) error {
	start := time.Now()

	conn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", k.config.Host, k.config.Port))
	if err != nil {
		k.LogActivity("CREATE_TOPIC", fmt.Sprintf("CREATE TOPIC '%s'", topic), time.Since(start), err, "")
		return err
	}
	defer conn.Close()

	numPartitions := 1
	replicationFactor := 1

	if np, ok := config["num_partitions"].(int); ok {
		numPartitions = np
	}
	if rf, ok := config["replication_factor"].(int); ok {
		replicationFactor = rf
	}

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = conn.CreateTopics(topicConfig)
	duration := time.Since(start)

	// Log activity
	command := fmt.Sprintf("CREATE TOPIC '%s' (partitions: %d, replication: %d)", topic, numPartitions, replicationFactor)
	response := ""
	if err == nil {
		response = fmt.Sprintf("Topic '%s' created successfully", topic)
	}
	k.LogActivity("CREATE_TOPIC", command, duration, err, response)

	return err
}

// DeleteTopic deletes a topic
func (k *KafkaAdapter) DeleteTopic(ctx context.Context, topic string) error {
	start := time.Now()

	conn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", k.config.Host, k.config.Port))
	if err != nil {
		k.LogActivity("DELETE_TOPIC", fmt.Sprintf("DELETE TOPIC '%s'", topic), time.Since(start), err, "")
		return err
	}
	defer conn.Close()

	err = conn.DeleteTopics(topic)
	duration := time.Since(start)

	// Log activity
	command := fmt.Sprintf("DELETE TOPIC '%s'", topic)
	response := ""
	if err == nil {
		response = fmt.Sprintf("Topic '%s' deleted successfully", topic)
	}
	k.LogActivity("DELETE_TOPIC", command, duration, err, response)

	return err
}

// ListTopics lists all topics
func (k *KafkaAdapter) ListTopics(ctx context.Context) ([]string, error) {
	start := time.Now()

	conn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", k.config.Host, k.config.Port))
	if err != nil {
		k.LogActivity("LIST_TOPICS", "LIST TOPICS", time.Since(start), err, "")
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		k.LogActivity("LIST_TOPICS", "LIST TOPICS", time.Since(start), err, "")
		return nil, err
	}

	// Extract unique topic names
	topicMap := make(map[string]bool)
	for i := range partitions {
		partition := partitions[i]
		topicMap[partition.Topic] = true
	}

	topics := make([]string, 0, len(topicMap))
	for topic := range topicMap {
		topics = append(topics, topic)
	}

	duration := time.Since(start)
	command := "LIST TOPICS"
	response := fmt.Sprintf("Found %d topics", len(topics))
	k.LogActivity("LIST_TOPICS", command, duration, nil, response)

	return topics, nil
}

// Ensure KafkaAdapter implements QueueAdapter
var _ adapters.QueueAdapter = (*KafkaAdapter)(nil)
