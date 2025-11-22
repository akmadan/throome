package monitor

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// ActivityLog represents a single service interaction
type ActivityLog struct {
	ID           string            `json:"id"`
	Timestamp    time.Time         `json:"timestamp"`
	ClusterID    string            `json:"cluster_id"`
	ServiceName  string            `json:"service_name"`
	ServiceType  string            `json:"service_type"`
	Operation    string            `json:"operation"`               // GET, SET, SELECT, PUBLISH, etc.
	Command      string            `json:"command"`                 // Full command/query
	Parameters   []interface{}     `json:"parameters,omitempty"`    // Query parameters
	Duration     int64             `json:"duration"`                // Duration in milliseconds
	Status       string            `json:"status"`                  // success, error
	Response     string            `json:"response,omitempty"`      // Result summary
	Error        string            `json:"error,omitempty"`         // Error message if failed
	RowsAffected int64             `json:"rows_affected,omitempty"` // For SQL queries
	ClientInfo   map[string]string `json:"client_info,omitempty"`   // Additional context
}

// ActivityBuffer is a thread-safe circular buffer for activity logs
type ActivityBuffer struct {
	logs     []*ActivityLog
	maxSize  int
	position int
	mu       sync.RWMutex
}

// NewActivityBuffer creates a new activity buffer with specified max size
func NewActivityBuffer(maxSize int) *ActivityBuffer {
	return &ActivityBuffer{
		logs:    make([]*ActivityLog, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds a new activity log to the buffer
func (ab *ActivityBuffer) Add(log *ActivityLog) {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	// Ensure ID is set
	if log.ID == "" {
		log.ID = uuid.New().String()
	}

	// If buffer is not full yet, append
	if len(ab.logs) < ab.maxSize {
		ab.logs = append(ab.logs, log)
	} else {
		// Buffer is full, overwrite oldest entry
		ab.logs[ab.position] = log
		ab.position = (ab.position + 1) % ab.maxSize
	}
}

// GetRecent returns the most recent n activity logs
func (ab *ActivityBuffer) GetRecent(limit int) []*ActivityLog {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	if limit <= 0 || limit > len(ab.logs) {
		limit = len(ab.logs)
	}

	result := make([]*ActivityLog, 0, limit)

	// If buffer hasn't wrapped yet
	if len(ab.logs) < ab.maxSize {
		start := len(ab.logs) - limit
		if start < 0 {
			start = 0
		}
		for i := len(ab.logs) - 1; i >= start; i-- {
			result = append(result, ab.logs[i])
		}
	} else {
		// Buffer has wrapped, need to handle circular logic
		pos := ab.position - 1
		if pos < 0 {
			pos = ab.maxSize - 1
		}

		for i := 0; i < limit; i++ {
			result = append(result, ab.logs[pos])
			pos--
			if pos < 0 {
				pos = ab.maxSize - 1
			}
		}
	}

	return result
}

// GetByCluster returns recent activity logs for a specific cluster
func (ab *ActivityBuffer) GetByCluster(clusterID string, limit int) []*ActivityLog {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	const maxLimit = 1000
	if limit <= 0 {
		limit = 0
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	// Don't allocate more than available logs
	if limit > len(ab.logs) {
		limit = len(ab.logs)
	}

	result := make([]*ActivityLog, 0, limit)
	count := 0

	// Iterate from newest to oldest
	if len(ab.logs) < ab.maxSize {
		for i := len(ab.logs) - 1; i >= 0 && count < limit; i-- {
			if ab.logs[i].ClusterID == clusterID {
				result = append(result, ab.logs[i])
				count++
			}
		}
	} else {
		pos := ab.position - 1
		if pos < 0 {
			pos = ab.maxSize - 1
		}

		checked := 0
		for checked < len(ab.logs) && count < limit {
			if ab.logs[pos].ClusterID == clusterID {
				result = append(result, ab.logs[pos])
				count++
			}
			pos--
			if pos < 0 {
				pos = ab.maxSize - 1
			}
			checked++
		}
	}

	return result
}

// GetByService returns recent activity logs for a specific service
func (ab *ActivityBuffer) GetByService(clusterID, serviceName string, limit int) []*ActivityLog {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	result := make([]*ActivityLog, 0, limit)
	count := 0

	// Iterate from newest to oldest
	if len(ab.logs) < ab.maxSize {
		for i := len(ab.logs) - 1; i >= 0 && count < limit; i-- {
			if ab.logs[i].ClusterID == clusterID && ab.logs[i].ServiceName == serviceName {
				result = append(result, ab.logs[i])
				count++
			}
		}
	} else {
		pos := ab.position - 1
		if pos < 0 {
			pos = ab.maxSize - 1
		}

		checked := 0
		for checked < len(ab.logs) && count < limit {
			if ab.logs[pos].ClusterID == clusterID && ab.logs[pos].ServiceName == serviceName {
				result = append(result, ab.logs[pos])
				count++
			}
			pos--
			if pos < 0 {
				pos = ab.maxSize - 1
			}
			checked++
		}
	}

	return result
}

// Clear removes all activity logs from the buffer
func (ab *ActivityBuffer) Clear() {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	ab.logs = make([]*ActivityLog, 0, ab.maxSize)
	ab.position = 0
}

// Size returns the current number of logs in the buffer
func (ab *ActivityBuffer) Size() int {
	ab.mu.RLock()
	defer ab.mu.RUnlock()
	return len(ab.logs)
}

// ActivityFilters defines filters for querying activity logs
type ActivityFilters struct {
	ClusterID   string
	ServiceName string
	ServiceType string
	Operation   string
	Status      string // success, error
	Since       *time.Time
	Limit       int
}

// Filter applies filters to activity logs
func (ab *ActivityBuffer) Filter(filters ActivityFilters) []*ActivityLog {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	limit := filters.Limit
	if limit <= 0 {
		limit = 100
	}

	result := make([]*ActivityLog, 0, limit)
	count := 0

	// Helper function to check if log matches filters
	matches := func(log *ActivityLog) bool {
		if filters.ClusterID != "" && log.ClusterID != filters.ClusterID {
			return false
		}
		if filters.ServiceName != "" && log.ServiceName != filters.ServiceName {
			return false
		}
		if filters.ServiceType != "" && log.ServiceType != filters.ServiceType {
			return false
		}
		if filters.Operation != "" && log.Operation != filters.Operation {
			return false
		}
		if filters.Status != "" && log.Status != filters.Status {
			return false
		}
		if filters.Since != nil && log.Timestamp.Before(*filters.Since) {
			return false
		}
		return true
	}

	// Iterate from newest to oldest
	if len(ab.logs) < ab.maxSize {
		for i := len(ab.logs) - 1; i >= 0 && count < limit; i-- {
			if matches(ab.logs[i]) {
				result = append(result, ab.logs[i])
				count++
			}
		}
	} else {
		pos := ab.position - 1
		if pos < 0 {
			pos = ab.maxSize - 1
		}

		checked := 0
		for checked < len(ab.logs) && count < limit {
			if matches(ab.logs[pos]) {
				result = append(result, ab.logs[pos])
				count++
			}
			pos--
			if pos < 0 {
				pos = ab.maxSize - 1
			}
			checked++
		}
	}

	return result
}
