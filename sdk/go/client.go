package throome

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
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Throome SDK client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// WithTimeout sets a custom timeout for the HTTP client
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.httpClient.Timeout = timeout
	return c
}

// Cluster returns a cluster client for the specified cluster ID
func (c *Client) Cluster(clusterID string) *ClusterClient {
	return &ClusterClient{
		client:    c,
		clusterID: clusterID,
	}
}

// request makes an HTTP request to the gateway
func (c *Client) request(ctx context.Context, method, path string, body, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Health checks the health of the gateway
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var health HealthResponse
	if err := c.request(ctx, "GET", "/api/v1/health", nil, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

// ListClusters lists all clusters
func (c *Client) ListClusters(ctx context.Context) ([]Cluster, error) {
	var clusters []Cluster
	if err := c.request(ctx, "GET", "/api/v1/clusters", nil, &clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

// GetCluster gets a specific cluster
func (c *Client) GetCluster(ctx context.Context, clusterID string) (*Cluster, error) {
	var cluster Cluster
	path := fmt.Sprintf("/api/v1/clusters/%s", clusterID)
	if err := c.request(ctx, "GET", path, nil, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}

// CreateCluster creates a new cluster
func (c *Client) CreateCluster(ctx context.Context, req CreateClusterRequest) (*CreateClusterResponse, error) {
	var resp CreateClusterResponse
	if err := c.request(ctx, "POST", "/api/v1/clusters", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DeleteCluster deletes a cluster
func (c *Client) DeleteCluster(ctx context.Context, clusterID string) error {
	path := fmt.Sprintf("/api/v1/clusters/%s", clusterID)
	return c.request(ctx, "DELETE", path, nil, nil)
}

// GetActivity gets global activity logs
func (c *Client) GetActivity(ctx context.Context, filters ActivityFilters) ([]ActivityLog, error) {
	var logs []ActivityLog
	path := "/api/v1/activity"
	if filters.Limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, filters.Limit)
	}
	if err := c.request(ctx, "GET", path, nil, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// ClusterClient provides cluster-specific operations
type ClusterClient struct {
	client    *Client
	clusterID string
}

// Health checks the health of the cluster
func (cc *ClusterClient) Health(ctx context.Context) (*ClusterHealthResponse, error) {
	var health ClusterHealthResponse
	path := fmt.Sprintf("/api/v1/clusters/%s/health", cc.clusterID)
	if err := cc.client.request(ctx, "GET", path, nil, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

// Metrics gets cluster metrics
func (cc *ClusterClient) Metrics(ctx context.Context) (*MetricsResponse, error) {
	var metrics MetricsResponse
	path := fmt.Sprintf("/api/v1/clusters/%s/metrics", cc.clusterID)
	if err := cc.client.request(ctx, "GET", path, nil, &metrics); err != nil {
		return nil, err
	}
	return &metrics, nil
}

// GetActivity gets cluster-specific activity logs
func (cc *ClusterClient) GetActivity(ctx context.Context, filters ActivityFilters) ([]ActivityLog, error) {
	var logs []ActivityLog
	path := fmt.Sprintf("/api/v1/clusters/%s/activity", cc.clusterID)
	if filters.Limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, filters.Limit)
	}
	if err := cc.client.request(ctx, "GET", path, nil, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}

// Service returns a service client
func (cc *ClusterClient) Service(serviceName string) *ServiceClient {
	return &ServiceClient{
		client:      cc.client,
		clusterID:   cc.clusterID,
		serviceName: serviceName,
	}
}

// DB returns a database client
func (cc *ClusterClient) DB() *DBClient {
	return &DBClient{clusterClient: cc}
}

// Cache returns a cache client
func (cc *ClusterClient) Cache() *CacheClient {
	return &CacheClient{clusterClient: cc}
}

// Queue returns a queue client
func (cc *ClusterClient) Queue() *QueueClient {
	return &QueueClient{clusterClient: cc}
}

// ServiceClient provides service-specific operations
type ServiceClient struct {
	client      *Client
	clusterID   string
	serviceName string
}

// GetInfo gets service information
func (sc *ServiceClient) GetInfo(ctx context.Context) (*ServiceInfo, error) {
	var info ServiceInfo
	path := fmt.Sprintf("/api/v1/clusters/%s/services/%s", sc.clusterID, sc.serviceName)
	if err := sc.client.request(ctx, "GET", path, nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// GetLogs gets service Docker container logs
func (sc *ServiceClient) GetLogs(ctx context.Context, options LogOptions) (string, error) {
	path := fmt.Sprintf("/api/v1/clusters/%s/services/%s/logs", sc.clusterID, sc.serviceName)
	if options.Tail > 0 {
		path = fmt.Sprintf("%s?tail=%d", path, options.Tail)
	}
	if options.Timestamps {
		if options.Tail > 0 {
			path = fmt.Sprintf("%s&timestamps=true", path)
		} else {
			path = fmt.Sprintf("%s?timestamps=true", path)
		}
	}

	url := fmt.Sprintf("%s%s", sc.client.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := sc.client.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	logs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// GetActivity gets service-specific activity logs
func (sc *ServiceClient) GetActivity(ctx context.Context, filters ActivityFilters) ([]ActivityLog, error) {
	var logs []ActivityLog
	path := fmt.Sprintf("/api/v1/clusters/%s/services/%s/activity", sc.clusterID, sc.serviceName)
	if filters.Limit > 0 {
		path = fmt.Sprintf("%s?limit=%d", path, filters.Limit)
	}
	if err := sc.client.request(ctx, "GET", path, nil, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
