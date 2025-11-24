package throome

import (
	"context"
	"fmt"
)

// DBClient provides database operations
type DBClient struct {
	clusterClient *ClusterClient
}

// Execute executes a SQL statement without returning results
func (d *DBClient) Execute(ctx context.Context, query string, args ...interface{}) error {
	req := DBQueryRequest{
		Query: query,
		Args:  args,
	}

	path := fmt.Sprintf("/api/v1/clusters/%s/db/execute", d.clusterClient.clusterID)
	return d.clusterClient.client.request(ctx, "POST", path, req, nil)
}

// Query executes a SQL query and returns results
func (d *DBClient) Query(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	req := DBQueryRequest{
		Query: query,
		Args:  args,
	}

	var resp DBQueryResponse
	path := fmt.Sprintf("/api/v1/clusters/%s/db/query", d.clusterClient.clusterID)
	if err := d.clusterClient.client.request(ctx, "POST", path, req, &resp); err != nil {
		return nil, err
	}

	return resp.Rows, nil
}

// QueryRow executes a query that returns a single row
func (d *DBClient) QueryRow(ctx context.Context, query string, args ...interface{}) (map[string]interface{}, error) {
	rows, err := d.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows returned")
	}

	return rows[0], nil
}

