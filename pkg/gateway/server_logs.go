package gateway

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/mux"
)

// handleGetServiceLogs returns Docker container logs for a service
func (s *Server) handleGetServiceLogs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]
	serviceName := vars["service_name"]

	// Get cluster config to find container ID
	cfg, err := s.gateway.GetClusterManager().Get(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	// Find service in cluster
	serviceConfig, exists := cfg.Services[serviceName]
	if !exists {
		s.errorResponse(w, http.StatusNotFound, "Service not found in cluster", nil)
		return
	}

	// Check if service has a container ID
	if serviceConfig.ContainerID == "" {
		s.errorResponse(w, http.StatusBadRequest, "Service is not provisioned by Throome", nil)
		return
	}

	// Parse query parameters
	tailLines := 100 // default
	if tailStr := r.URL.Query().Get("tail"); tailStr != "" {
		if tail, err := strconv.Atoi(tailStr); err == nil && tail > 0 {
			tailLines = tail
		}
	}

	timestamps := r.URL.Query().Get("timestamps") == "true"

	// Create Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to create Docker client", err)
		return
	}
	defer dockerClient.Close()

	// Get container logs
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: timestamps,
		Tail:       strconv.Itoa(tailLines),
	}

	logs, err := dockerClient.ContainerLogs(ctx, serviceConfig.ContainerID, options)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to get container logs", err)
		return
	}
	defer logs.Close()

	// Read logs
	logBytes, err := io.ReadAll(logs)
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to read logs", err)
		return
	}

	// Return logs as plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(logBytes)
}

// handleGetServiceInfo returns service information including container status
func (s *Server) handleGetServiceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]
	serviceName := vars["service_name"]

	// Get cluster config
	cfg, err := s.gateway.GetClusterManager().Get(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	// Find service in cluster
	serviceConfig, exists := cfg.Services[serviceName]
	if !exists {
		s.errorResponse(w, http.StatusNotFound, "Service not found in cluster", nil)
		return
	}

	// Build response
	response := map[string]interface{}{
		"cluster_id":   clusterID,
		"cluster_name": cfg.Name,
		"service_name": serviceName,
		"type":         serviceConfig.Type,
		"host":         serviceConfig.Host,
		"port":         serviceConfig.Port,
		"container_id": serviceConfig.ContainerID,
	}

	// If service has a container, get its status
	if serviceConfig.ContainerID != "" {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err == nil {
			defer dockerClient.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			inspect, err := dockerClient.ContainerInspect(ctx, serviceConfig.ContainerID)
			if err == nil {
				response["container_status"] = inspect.State.Status
				response["container_running"] = inspect.State.Running
				response["container_started_at"] = inspect.State.StartedAt
				response["container_image"] = inspect.Config.Image
			}
		}
	}

	// Add database-specific fields for PostgreSQL
	if serviceConfig.Type == "postgres" {
		response["database"] = serviceConfig.Database
		response["username"] = serviceConfig.Username
	}

	s.jsonResponse(w, http.StatusOK, response)
}
