package monitor

import (
	"time"
)

// ActivityLogger provides methods for logging service interactions
type ActivityLogger interface {
	Log(activity *ActivityLog)
	LogOperation(clusterID, serviceName, serviceType, operation, command string, duration time.Duration, err error, response string)
}

// DefaultActivityLogger implements ActivityLogger using an ActivityBuffer
type DefaultActivityLogger struct {
	buffer *ActivityBuffer
}

// NewActivityLogger creates a new activity logger with the given buffer
func NewActivityLogger(buffer *ActivityBuffer) ActivityLogger {
	return &DefaultActivityLogger{
		buffer: buffer,
	}
}

// Log adds an activity log to the buffer
func (l *DefaultActivityLogger) Log(activity *ActivityLog) {
	if l.buffer != nil {
		l.buffer.Add(activity)
	}
}

// LogOperation is a convenience method for logging an operation
func (l *DefaultActivityLogger) LogOperation(
	clusterID, serviceName, serviceType, operation, command string,
	duration time.Duration,
	err error,
	response string,
) {
	activity := &ActivityLog{
		Timestamp:   time.Now(),
		ClusterID:   clusterID,
		ServiceName: serviceName,
		ServiceType: serviceType,
		Operation:   operation,
		Command:     command,
		Duration:    duration.Milliseconds(),
		Response:    response,
	}

	if err != nil {
		activity.Status = "error"
		activity.Error = err.Error()
	} else {
		activity.Status = "success"
	}

	l.Log(activity)
}

// NoOpActivityLogger is a logger that does nothing (for testing or when logging is disabled)
type NoOpActivityLogger struct{}

// NewNoOpActivityLogger creates a no-op activity logger
func NewNoOpActivityLogger() ActivityLogger {
	return &NoOpActivityLogger{}
}

// Log does nothing
func (l *NoOpActivityLogger) Log(activity *ActivityLog) {}

// LogOperation does nothing
func (l *NoOpActivityLogger) LogOperation(
	clusterID, serviceName, serviceType, operation, command string,
	duration time.Duration,
	err error,
	response string,
) {
}
