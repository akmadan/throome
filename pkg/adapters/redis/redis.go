package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/akmadan/throome/pkg/adapters"
	"github.com/akmadan/throome/pkg/cluster"
)

// RedisAdapter implements the CacheAdapter interface for Redis
type RedisAdapter struct {
	*adapters.BaseAdapter
	config cluster.ServiceConfig
	client *redis.Client
}

// NewRedisAdapter creates a new Redis adapter
func NewRedisAdapter(config cluster.ServiceConfig) (adapters.Adapter, error) {
	adapter := &RedisAdapter{
		BaseAdapter: adapters.NewBaseAdapter(config),
		config:      config,
	}
	return adapter, nil
}

// Connect establishes a connection to Redis
func (r *RedisAdapter) Connect(ctx context.Context) error {
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.config.Host, r.config.Port),
		Password: r.config.Password,
		DB:       0, // default DB
	}

	// Get DB from options if specified
	if db, ok := r.config.Options["db"].(int); ok {
		options.DB = db
	}

	// Configure pool
	if r.config.Pool.MaxConnections > 0 {
		options.PoolSize = r.config.Pool.MaxConnections
	}
	if r.config.Pool.MinConnections > 0 {
		options.MinIdleConns = r.config.Pool.MinConnections
	}
	if r.config.Pool.MaxIdleTime > 0 {
		options.IdleTimeout = time.Duration(r.config.Pool.MaxIdleTime) * time.Second
	}

	r.client = redis.NewClient(options)

	// Test connection
	if err := r.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	r.SetConnected(true)
	return nil
}

// Disconnect closes the Redis connection
func (r *RedisAdapter) Disconnect(ctx context.Context) error {
	if r.client != nil {
		err := r.client.Close()
		r.SetConnected(false)
		return err
	}
	return nil
}

// Ping checks if the Redis connection is alive
func (r *RedisAdapter) Ping(ctx context.Context) error {
	start := time.Now()
	err := r.client.Ping(ctx).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// HealthCheck performs a health check
func (r *RedisAdapter) HealthCheck(ctx context.Context) (*adapters.HealthStatus, error) {
	start := time.Now()
	err := r.Ping(ctx)
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

// Get retrieves a value from Redis
func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := r.client.Get(ctx, key).Result()
	r.RecordRequest(time.Since(start), err == nil || err == redis.Nil)

	if err == redis.Nil {
		return "", nil // Key doesn't exist
	}

	return val, err
}

// Set sets a value in Redis
func (r *RedisAdapter) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Set(ctx, key, value, expiration).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// Delete deletes a key from Redis
func (r *RedisAdapter) Delete(ctx context.Context, key string) error {
	start := time.Now()
	err := r.client.Del(ctx, key).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// Exists checks if a key exists in Redis
func (r *RedisAdapter) Exists(ctx context.Context, key string) (bool, error) {
	start := time.Now()
	count, err := r.client.Exists(ctx, key).Result()
	r.RecordRequest(time.Since(start), err == nil)
	return count > 0, err
}

// Keys returns keys matching a pattern
func (r *RedisAdapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	start := time.Now()
	keys, err := r.client.Keys(ctx, pattern).Result()
	r.RecordRequest(time.Since(start), err == nil)
	return keys, err
}

// TTL returns the time-to-live of a key
func (r *RedisAdapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	start := time.Now()
	ttl, err := r.client.TTL(ctx, key).Result()
	r.RecordRequest(time.Since(start), err == nil)
	return ttl, err
}

// Expire sets expiration on a key
func (r *RedisAdapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	start := time.Now()
	err := r.client.Expire(ctx, key, expiration).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// Additional Redis-specific operations

// HSet sets a field in a hash
func (r *RedisAdapter) HSet(ctx context.Context, key string, field string, value string) error {
	start := time.Now()
	err := r.client.HSet(ctx, key, field, value).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// HGet gets a field from a hash
func (r *RedisAdapter) HGet(ctx context.Context, key string, field string) (string, error) {
	start := time.Now()
	val, err := r.client.HGet(ctx, key, field).Result()
	r.RecordRequest(time.Since(start), err == nil || err == redis.Nil)

	if err == redis.Nil {
		return "", nil
	}

	return val, err
}

// LPush pushes values to the head of a list
func (r *RedisAdapter) LPush(ctx context.Context, key string, values ...string) error {
	start := time.Now()
	err := r.client.LPush(ctx, key, values).Err()
	r.RecordRequest(time.Since(start), err == nil)
	return err
}

// RPop removes and returns the last element of a list
func (r *RedisAdapter) RPop(ctx context.Context, key string) (string, error) {
	start := time.Now()
	val, err := r.client.RPop(ctx, key).Result()
	r.RecordRequest(time.Since(start), err == nil || err == redis.Nil)

	if err == redis.Nil {
		return "", nil
	}

	return val, err
}

// Incr increments a counter
func (r *RedisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	start := time.Now()
	val, err := r.client.Incr(ctx, key).Result()
	r.RecordRequest(time.Since(start), err == nil)
	return val, err
}

// Ensure RedisAdapter implements CacheAdapter
var _ adapters.CacheAdapter = (*RedisAdapter)(nil)
