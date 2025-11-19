package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// ClusterIDPattern is the regex pattern for valid cluster IDs
// Must start and end with alphanumeric, can contain hyphens in the middle
var ClusterIDPattern = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

// ValidateClusterID validates a cluster ID
func ValidateClusterID(id string) error {
	if id == "" {
		return fmt.Errorf("cluster ID cannot be empty")
	}

	if len(id) < 3 || len(id) > 32 {
		return fmt.Errorf("cluster ID must be between 3 and 32 characters")
	}

	if !ClusterIDPattern.MatchString(id) {
		return fmt.Errorf("cluster ID must contain only lowercase letters, numbers, and hyphens, and must start and end with a letter or number")
	}

	return nil
}

// ValidateClusterName validates a cluster name
func ValidateClusterName(name string) error {
	if name == "" {
		return fmt.Errorf("cluster name cannot be empty")
	}

	if len(name) > 64 {
		return fmt.Errorf("cluster name cannot exceed 64 characters")
	}

	return nil
}

// SanitizeClusterName converts a cluster name to a valid cluster ID
func SanitizeClusterName(name string) string {
	// Convert to lowercase
	id := strings.ToLower(name)

	// Replace spaces and underscores with hyphens
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "_", "-")

	// Remove invalid characters
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	id = result.String()

	// Remove leading/trailing hyphens
	id = strings.Trim(id, "-")

	// Collapse multiple hyphens
	for strings.Contains(id, "--") {
		id = strings.ReplaceAll(id, "--", "-")
	}

	// Truncate if too long
	if len(id) > 32 {
		id = id[:32]
		// Remove trailing hyphen if any
		id = strings.TrimRight(id, "-")
	}

	// Ensure minimum length
	if len(id) < 3 {
		id = id + "-01"
	}

	return id
}

// ValidatePort validates a port number
func ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

// ValidateHost validates a host address
func ValidateHost(host string) error {
	if host == "" {
		return fmt.Errorf("host cannot be empty")
	}
	return nil
}
