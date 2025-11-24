package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/akmadan/throome/internal/logger"
	"github.com/akmadan/throome/pkg/adapters/postgres"
	"github.com/akmadan/throome/pkg/adapters/redis"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// Database operation request/response types
type DBExecuteRequest struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args"`
}

type DBQueryRequest struct {
	Query string        `json:"query"`
	Args  []interface{} `json:"args"`
}

type DBQueryResponse struct {
	Rows []map[string]interface{} `json:"rows"`
}

type DBExecuteResponse struct {
	RowsAffected int64 `json:"rows_affected"`
}

// Cache operation request/response types
type CacheGetRequest struct {
	Key string `json:"key"`
}

type CacheSetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"` // TTL in seconds
}

type CacheDeleteRequest struct {
	Key string `json:"key"`
}

type CacheGetResponse struct {
	Value string `json:"value"`
}

// handleDBExecute handles database execute operations (INSERT, UPDATE, DELETE, DDL)
func (s *Server) handleDBExecute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	var req DBExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Find the PostgreSQL service in the cluster
	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	var postgresService string
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Type == "postgres" {
			postgresService = serviceName
			break
		}
	}

	if postgresService == "" {
		s.errorResponse(w, http.StatusNotFound, "No PostgreSQL service found in cluster", nil)
		return
	}

	// Get the adapter
	adapter, err := s.gateway.GetAdapter(clusterID, postgresService)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get database adapter", err)
		return
	}

	// Type assert to PostgresAdapter
	pgAdapter, ok := adapter.(*postgres.PostgresAdapter)
	if !ok {
		s.errorResponse(w, http.StatusInternalServerError, "Adapter is not a PostgresAdapter", nil)
		return
	}

	// Execute the query
	result, err := pgAdapter.Execute(r.Context(), req.Query, req.Args...)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to execute query", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, DBExecuteResponse{
		RowsAffected: result.RowsAffected(),
	})
}

// handleDBQuery handles database query operations (SELECT)
func (s *Server) handleDBQuery(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	var req DBQueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Find the PostgreSQL service in the cluster
	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	var postgresService string
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Type == "postgres" {
			postgresService = serviceName
			break
		}
	}

	if postgresService == "" {
		s.errorResponse(w, http.StatusNotFound, "No PostgreSQL service found in cluster", nil)
		return
	}

	// Get the adapter
	adapter, err := s.gateway.GetAdapter(clusterID, postgresService)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get database adapter", err)
		return
	}

	// Type assert to PostgresAdapter
	pgAdapter, ok := adapter.(*postgres.PostgresAdapter)
	if !ok {
		s.errorResponse(w, http.StatusInternalServerError, "Adapter is not a PostgresAdapter", nil)
		return
	}

	// Execute the query directly with pgx to get access to pgx.Rows
	pool := pgAdapter.GetPool()
	pgxRows, err := pool.Query(r.Context(), req.Query, req.Args...)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to execute query", err)
		return
	}
	defer pgxRows.Close()

	// Use pgx.CollectRows to convert rows to maps
	result, err := pgx.CollectRows(pgxRows, pgx.RowToMap)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to collect rows", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, DBQueryResponse{
		Rows: result,
	})
}

// handleCacheGet handles cache GET operations
func (s *Server) handleCacheGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	var req CacheGetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Find the Redis service in the cluster
	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	var redisService string
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Type == "redis" {
			redisService = serviceName
			break
		}
	}

	if redisService == "" {
		s.errorResponse(w, http.StatusNotFound, "No Redis service found in cluster", nil)
		return
	}

	// Get the adapter
	adapter, err := s.gateway.GetAdapter(clusterID, redisService)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get cache adapter", err)
		return
	}

	// Type assert to RedisAdapter
	redisAdapter, ok := adapter.(*redis.RedisAdapter)
	if !ok {
		s.errorResponse(w, http.StatusInternalServerError, "Adapter is not a RedisAdapter", nil)
		return
	}

	// Get the value
	value, err := redisAdapter.Get(r.Context(), req.Key)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get key", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, CacheGetResponse{
		Value: value,
	})
}

// handleCacheSet handles cache SET operations
func (s *Server) handleCacheSet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	var req CacheSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Find the Redis service in the cluster
	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	var redisService string
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Type == "redis" {
			redisService = serviceName
			break
		}
	}

	if redisService == "" {
		s.errorResponse(w, http.StatusNotFound, "No Redis service found in cluster", nil)
		return
	}

	// Get the adapter
	adapter, err := s.gateway.GetAdapter(clusterID, redisService)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get cache adapter", err)
		return
	}

	// Type assert to RedisAdapter
	redisAdapter, ok := adapter.(*redis.RedisAdapter)
	if !ok {
		s.errorResponse(w, http.StatusInternalServerError, "Adapter is not a RedisAdapter", nil)
		return
	}

	// Set the value
	ttl := time.Duration(req.TTL) * time.Second
	if err := redisAdapter.Set(r.Context(), req.Key, req.Value, ttl); err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to set key", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}

// handleCacheDelete handles cache DELETE operations
func (s *Server) handleCacheDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	var req CacheDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Find the Redis service in the cluster
	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	var redisService string
	for serviceName, serviceConfig := range config.Services {
		if serviceConfig.Type == "redis" {
			redisService = serviceName
			break
		}
	}

	if redisService == "" {
		s.errorResponse(w, http.StatusNotFound, "No Redis service found in cluster", nil)
		return
	}

	// Get the adapter
	adapter, err := s.gateway.GetAdapter(clusterID, redisService)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get cache adapter", err)
		return
	}

	// Type assert to RedisAdapter
	redisAdapter, ok := adapter.(*redis.RedisAdapter)
	if !ok {
		s.errorResponse(w, http.StatusInternalServerError, "Adapter is not a RedisAdapter", nil)
		return
	}

	// Delete the key
	if err := redisAdapter.Delete(r.Context(), req.Key); err != nil {
		logger.Error("Failed to delete key", zap.Error(err))
		s.errorResponse(w, http.StatusInternalServerError, "Failed to delete key", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}
