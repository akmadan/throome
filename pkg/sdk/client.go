package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is the Throome SDK client
type Client struct {
	gatewayURL string
	clusterID  string
	httpClient *http.Client
}

// NewClient creates a new Throome SDK client
func NewClient(gatewayURL, clusterID string) *Client {
	return &Client{
		gatewayURL: gatewayURL,
		clusterID:  clusterID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// DB returns a database client
func (c *Client) DB() *DBClient {
	return &DBClient{client: c}
}

// Cache returns a cache client
func (c *Client) Cache() *CacheClient {
	return &CacheClient{client: c}
}

// Queue returns a queue client
func (c *Client) Queue() *QueueClient {
	return &QueueClient{client: c}
}

// Health checks the health of the cluster
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/clusters/%s/health", c.gatewayURL, c.clusterID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check failed: %s", resp.Status)
	}

	var healthResp HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		return nil, err
	}

	return &healthResp, nil
}

// request makes an HTTP request to the gateway
func (c *Client) request(ctx context.Context, method, endpoint string, body, result interface{}) error {
	url := fmt.Sprintf("%s/api/v1/clusters/%s/%s", c.gatewayURL, c.clusterID, endpoint)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp) //nolint:errcheck // Error response decode is best-effort
		return fmt.Errorf("request failed: %s - %v", resp.Status, errResp)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}
	}

	return nil
}

// Response types

// HealthResponse represents a health check response
type HealthResponse struct {
	ClusterID string                   `json:"cluster_id"`
	Services  map[string]ServiceHealth `json:"services"`
}

// ServiceHealth represents the health of a service
type ServiceHealth struct {
	Healthy      bool   `json:"healthy"`
	ResponseTime int64  `json:"response_time"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// DBClient provides database operations
type DBClient struct {
	client *Client
}

// Execute executes a query
func (d *DBClient) Execute(ctx context.Context, query string, args ...interface{}) error {
	req := map[string]interface{}{
		"query": query,
		"args":  args,
	}

	return d.client.request(ctx, "POST", "db/execute", req, nil)
}

// Query executes a query and returns results
func (d *DBClient) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	req := map[string]interface{}{
		"query": query,
		"args":  args,
	}

	var result struct {
		Rows []map[string]interface{} `json:"rows"`
	}

	if err := d.client.request(ctx, "POST", "db/query", req, &result); err != nil {
		return nil, err
	}

	return result.Rows, nil
}

// CacheClient provides cache operations
type CacheClient struct {
	client *Client
}

// Get retrieves a value from cache
func (c *CacheClient) Get(ctx context.Context, key string) (string, error) {
	req := map[string]interface{}{
		"key": key,
	}

	var result struct {
		Value string `json:"value"`
	}

	if err := c.client.request(ctx, "POST", "cache/get", req, &result); err != nil {
		return "", err
	}

	return result.Value, nil
}

// Set sets a value in cache
func (c *CacheClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	req := map[string]interface{}{
		"key":        key,
		"value":      value,
		"expiration": expiration.Seconds(),
	}

	return c.client.request(ctx, "POST", "cache/set", req, nil)
}

// Delete deletes a key from cache
func (c *CacheClient) Delete(ctx context.Context, key string) error {
	req := map[string]interface{}{
		"key": key,
	}

	return c.client.request(ctx, "POST", "cache/delete", req, nil)
}

// QueueClient provides queue operations
type QueueClient struct {
	client *Client
}

// Publish publishes a message to a topic
func (q *QueueClient) Publish(ctx context.Context, topic string, message []byte) error {
	req := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}

	return q.client.request(ctx, "POST", "queue/publish", req, nil)
}

// Subscribe subscribes to a topic
func (q *QueueClient) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	// Note: This would typically use WebSocket or long-polling
	// For now, this is a placeholder
	return fmt.Errorf("subscribe not yet implemented in SDK")
}
