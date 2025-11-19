package utils

import (
	"testing"
)

func TestValidateClusterID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"valid lowercase", "test-cluster-01", false},
		{"valid with numbers", "app123", false},
		{"valid single word", "production", false},
		{"empty string", "", true},
		{"too short", "ab", true},
		{"too long", "this-is-a-very-long-cluster-id-that-exceeds-32-characters", true},
		{"uppercase", "TestCluster", true},
		{"starts with hyphen", "-test", true},
		{"ends with hyphen", "test-", true},
		{"special chars", "test@cluster", true},
		{"spaces", "test cluster", true},
		{"underscores", "test_cluster", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClusterID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClusterID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}
		})
	}
}

func TestValidateClusterName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid name", "My Application", false},
		{"valid with special chars", "App-2023 (Production)", false},
		{"empty", "", true},
		{"too long", string(make([]byte, 65)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClusterName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClusterName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeClusterName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple lowercase",
			input: "myapp",
			want:  "myapp",
		},
		{
			name:  "uppercase to lowercase",
			input: "MyApp",
			want:  "myapp",
		},
		{
			name:  "spaces to hyphens",
			input: "My Application",
			want:  "my-application",
		},
		{
			name:  "underscores to hyphens",
			input: "my_app_name",
			want:  "my-app-name",
		},
		{
			name:  "remove special characters",
			input: "my@app#2023!",
			want:  "myapp2023",
		},
		{
			name:  "multiple hyphens collapsed",
			input: "my---app",
			want:  "my-app",
		},
		{
			name:  "trim hyphens",
			input: "-myapp-",
			want:  "myapp",
		},
		{
			name:  "too long truncate",
			input: "this-is-a-very-long-application-name-that-needs-truncation",
			want:  "this-is-a-very-long-application",
		},
		{
			name:  "too short add suffix",
			input: "ab",
			want:  "ab-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeClusterName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeClusterName(%q) = %q, want %q", tt.input, got, tt.want)
			}
			
			// Verify result is a valid cluster ID
			if err := ValidateClusterID(got); err != nil {
				t.Errorf("SanitizeClusterName(%q) produced invalid ID %q: %v", tt.input, got, err)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid port 80", 80, false},
		{"valid port 8080", 8080, false},
		{"valid port 65535", 65535, false},
		{"invalid port 0", 0, true},
		{"invalid port negative", -1, true},
		{"invalid port too high", 65536, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestValidateHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{"valid localhost", "localhost", false},
		{"valid IP", "127.0.0.1", false},
		{"valid domain", "example.com", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHost(tt.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateHost(%q) error = %v, wantErr %v", tt.host, err, tt.wantErr)
			}
		})
	}
}

