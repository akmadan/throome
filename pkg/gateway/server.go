package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/akshitmadan/throome/internal/config"
	"github.com/akshitmadan/throome/internal/logger"
	"go.uber.org/zap"
)

// Server represents the HTTP server for the gateway
type Server struct {
	config  *config.AppConfig
	gateway *Gateway
	router  *mux.Router
	server  *http.Server
}

// NewServer creates a new HTTP server
func NewServer(cfg *config.AppConfig, gateway *Gateway) *Server {
	s := &Server{
		config:  cfg,
		gateway: gateway,
		router:  mux.NewRouter(),
	}

	s.setupRoutes()

	return s
}

// setupRoutes sets up HTTP routes
func (s *Server) setupRoutes() {
	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Cluster management
	api.HandleFunc("/clusters", s.handleListClusters).Methods("GET")
	api.HandleFunc("/clusters", s.handleCreateCluster).Methods("POST")
	api.HandleFunc("/clusters/{cluster_id}", s.handleGetCluster).Methods("GET")
	api.HandleFunc("/clusters/{cluster_id}", s.handleDeleteCluster).Methods("DELETE")

	// Health and metrics
	api.HandleFunc("/health", s.handleHealth).Methods("GET")
	api.HandleFunc("/clusters/{cluster_id}/health", s.handleClusterHealth).Methods("GET")
	api.HandleFunc("/clusters/{cluster_id}/metrics", s.handleClusterMetrics).Methods("GET")

	// Prometheus metrics endpoint
	if s.config.Monitoring.Enabled {
		s.router.Handle(s.config.Monitoring.MetricsPath, promhttp.Handler())
	}

	// Root endpoint
	s.router.HandleFunc("/", s.handleRoot).Methods("GET")

	// Middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.corsMiddleware)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	logger.Info("Starting HTTP server", zap.String("addr", addr))

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

// HTTP Handlers

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": "Throome Gateway",
		"version": "0.1.0",
		"status":  "running",
	}
	s.jsonResponse(w, http.StatusOK, response)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	}
	s.jsonResponse(w, http.StatusOK, response)
}

func (s *Server) handleListClusters(w http.ResponseWriter, r *http.Request) {
	clusters, err := s.gateway.ListClusters()
	if err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to list clusters", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"clusters": clusters,
		"count":    len(clusters),
	})
}

func (s *Server) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.errorResponse(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// TODO: Convert req.Config to cluster.Config
	// For now, return not implemented
	s.errorResponse(w, http.StatusNotImplemented, "Cluster creation via API not yet implemented", nil)
}

func (s *Server) handleGetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	config, err := s.gateway.GetClusterConfig(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, config)
}

func (s *Server) handleDeleteCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	if err := s.gateway.DeleteCluster(r.Context(), clusterID); err != nil {
		s.errorResponse(w, http.StatusInternalServerError, "Failed to delete cluster", err)
		return
	}

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Cluster deleted successfully",
	})
}

func (s *Server) handleClusterHealth(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	router, err := s.gateway.GetRouter(clusterID)
	if err != nil {
		s.errorResponse(w, http.StatusNotFound, "Cluster not found", err)
		return
	}

	healthStatuses := router.HealthCheckAll(r.Context())

	s.jsonResponse(w, http.StatusOK, map[string]interface{}{
		"cluster_id": clusterID,
		"services":   healthStatuses,
	})
}

func (s *Server) handleClusterMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clusterID := vars["cluster_id"]

	metrics := s.gateway.GetCollector().GetClusterMetrics(clusterID)
	if metrics == nil {
		s.errorResponse(w, http.StatusNotFound, "No metrics found for cluster", nil)
		return
	}

	s.jsonResponse(w, http.StatusOK, metrics)
}

// Middleware

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the next handler
		next.ServeHTTP(w, r)

		// Log request
		logger.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Duration("duration", time.Since(start)),
		)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper methods

func (s *Server) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) errorResponse(w http.ResponseWriter, status int, message string, err error) {
	response := map[string]interface{}{
		"error":  message,
		"status": status,
	}

	if err != nil {
		response["details"] = err.Error()
	}

	s.jsonResponse(w, status, response)
}
