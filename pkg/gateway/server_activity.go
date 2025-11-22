package gateway

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akmadan/throome/pkg/monitor"
	"github.com/gorilla/mux"
)

// handleGetActivity returns global activity logs
func (s *Server) handleGetActivity(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()

	filters := monitor.ActivityFilters{
		ClusterID:   query.Get("cluster_id"),
		ServiceType: query.Get("service_type"),
		Operation:   query.Get("operation"),
		Status:      query.Get("status"),
		Limit:       100, // default
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			if limit > 1000 {
				limit = 1000 // max limit
			}
			filters.Limit = limit
		}
	}

	// Parse since timestamp
	if sinceStr := query.Get("since"); sinceStr != "" {
		if since, err := time.Parse(time.RFC3339, sinceStr); err == nil {
			filters.Since = &since
		}
	}

	// Get activity buffer
	buffer := s.gateway.GetActivityBuffer()

	// Apply filters
	activities := buffer.Filter(filters)

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"activities": activities,
		"count":      len(activities),
		"filters":    filters,
	})
}

// handleGetClusterActivity returns activity logs for a specific cluster
func (s *Server) handleGetClusterActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	// Parse query parameters
	query := r.URL.Query()
	limit := 100

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 1000 {
				l = 1000
			}
			limit = l
		}
	}

	// Get activity buffer
	buffer := s.gateway.GetActivityBuffer()

	// Get activities for this cluster
	activities := buffer.GetByCluster(clusterID, limit)

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"activities": activities,
		"count":      len(activities),
		"cluster_id": clusterID,
	})
}

// handleGetServiceActivity returns activity logs for a specific service
func (s *Server) handleGetServiceActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]
	serviceName := vars["service_name"]

	// Parse query parameters
	query := r.URL.Query()
	limit := 100

	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l > 1000 {
				l = 1000
			}
			limit = l
		}
	}

	// Get activity buffer
	buffer := s.gateway.GetActivityBuffer()

	// Get activities for this service
	activities := buffer.GetByService(clusterID, serviceName, limit)

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"activities":   activities,
		"count":        len(activities),
		"cluster_id":   clusterID,
		"service_name": serviceName,
	})
}
