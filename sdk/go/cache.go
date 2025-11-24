package throome

import (
	"context"
	"fmt"
	"time"
)

// CacheClient provides cache operations
type CacheClient struct {
	clusterClient *ClusterClient
}

// Get retrieves a value from cache
func (c *CacheClient) Get(ctx context.Context, key string) (string, error) {
	req := CacheGetRequest{
		Key: key,
	}

	var resp CacheGetResponse
	path := fmt.Sprintf("/api/v1/clusters/%s/cache/get", c.clusterClient.clusterID)
	if err := c.clusterClient.client.request(ctx, "POST", path, req, &resp); err != nil {
		return "", err
	}

	return resp.Value, nil
}

// Set sets a value in cache
func (c *CacheClient) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	req := CacheSetRequest{
		Key:        key,
		Value:      value,
		Expiration: expiration.Seconds(),
	}

	path := fmt.Sprintf("/api/v1/clusters/%s/cache/set", c.clusterClient.clusterID)
	return c.clusterClient.client.request(ctx, "POST", path, req, nil)
}

// Delete deletes a key from cache
func (c *CacheClient) Delete(ctx context.Context, key string) error {
	req := CacheDeleteRequest{
		Key: key,
	}

	path := fmt.Sprintf("/api/v1/clusters/%s/cache/delete", c.clusterClient.clusterID)
	return c.clusterClient.client.request(ctx, "POST", path, req, nil)
}
