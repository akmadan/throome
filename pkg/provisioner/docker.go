package provisioner

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/akmadan/throome/internal/logger"
	"github.com/akmadan/throome/pkg/cluster"
	"go.uber.org/zap"
)

// DockerProvisioner handles Docker container lifecycle
type DockerProvisioner struct {
	client *client.Client
}

// ServiceContainer represents a provisioned container
type ServiceContainer struct {
	ContainerID string
	Name        string
	Type        string
	Port        int
	Status      string
}

// NewDockerProvisioner creates a new Docker provisioner
func NewDockerProvisioner() (*DockerProvisioner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerProvisioner{
		client: cli,
	}, nil
}

// ProvisionService provisions a new service container
func (p *DockerProvisioner) ProvisionService(ctx context.Context, serviceName string, config *cluster.ServiceConfig) (*ServiceContainer, error) {
	logger.Info("Provisioning service",
		zap.String("name", serviceName),
		zap.String("type", config.Type),
		zap.Int("port", config.Port),
	)

	// Determine image and environment based on service type
	var imageName string
	var env []string
	var healthCheck *container.HealthConfig

	switch config.Type {
	case "postgres":
		imageName = "postgres:17-alpine"
		env = []string{
			fmt.Sprintf("POSTGRES_USER=%s", getOrDefault(config.Username, "postgres")),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", getOrDefault(config.Password, "password")),
			fmt.Sprintf("POSTGRES_DB=%s", getOrDefault(config.Database, "postgres")),
		}
		healthCheck = &container.HealthConfig{
			Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
			Interval: 5 * time.Second,
			Timeout:  3 * time.Second,
			Retries:  3,
		}

	case "redis":
		imageName = "redis:7-alpine"
		env = []string{}
		if config.Password != "" {
			env = append(env, fmt.Sprintf("REDIS_PASSWORD=%s", config.Password))
		}
		healthCheck = &container.HealthConfig{
			Test:     []string{"CMD", "redis-cli", "ping"},
			Interval: 5 * time.Second,
			Timeout:  3 * time.Second,
			Retries:  3,
		}

	case "kafka":
		imageName = "apache/kafka:latest"
		env = []string{
			"KAFKA_BROKER_ID=1",
			fmt.Sprintf("KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:%d", config.Port),
			"KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT",
			"KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT",
			"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
		}
		healthCheck = &container.HealthConfig{
			Test:     []string{"CMD-SHELL", "kafka-broker-api-versions.sh --bootstrap-server localhost:9092 || exit 1"},
			Interval: 10 * time.Second,
			Timeout:  5 * time.Second,
			Retries:  5,
		}

	default:
		return nil, fmt.Errorf("unsupported service type: %s", config.Type)
	}

	// Pull image if not present
	logger.Info("Pulling Docker image", zap.String("image", imageName))
	reader, err := p.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer reader.Close()
	// Discard pull output
	_, _ = io.Copy(io.Discard, reader)

	// Create container configuration
	containerName := fmt.Sprintf("throome-%s", serviceName)
	
	// Port binding
	exposedPorts := nat.PortSet{
		nat.Port(fmt.Sprintf("%d/tcp", config.Port)): struct{}{},
	}
	portBindings := nat.PortMap{
		nat.Port(fmt.Sprintf("%d/tcp", getInternalPort(config.Type))): []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: fmt.Sprintf("%d", config.Port),
			},
		},
	}

	// Create container
	resp, err := p.client.ContainerCreate(ctx,
		&container.Config{
			Image:        imageName,
			Env:          env,
			ExposedPorts: exposedPorts,
			Healthcheck:  healthCheck,
			Labels: map[string]string{
				"throome.managed": "true",
				"throome.service": serviceName,
				"throome.type":    config.Type,
			},
		},
		&container.HostConfig{
			PortBindings: portBindings,
			RestartPolicy: container.RestartPolicy{
				Name: container.RestartPolicyUnlessStopped,
			},
		},
		nil,
		nil,
		containerName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Start container
	if err := p.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	logger.Info("Container started successfully",
		zap.String("container_id", resp.ID[:12]),
		zap.String("name", containerName),
		zap.Int("port", config.Port),
	)

	return &ServiceContainer{
		ContainerID: resp.ID,
		Name:        serviceName,
		Type:        config.Type,
		Port:        config.Port,
		Status:      "running",
	}, nil
}

// StopService stops a running container
func (p *DockerProvisioner) StopService(ctx context.Context, containerID string) error {
	timeout := 10
	return p.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RestartService restarts a container
func (p *DockerProvisioner) RestartService(ctx context.Context, containerID string) error {
	timeout := 10
	return p.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RemoveService stops and removes a container
func (p *DockerProvisioner) RemoveService(ctx context.Context, containerID string) error {
	// Stop first
	_ = p.StopService(ctx, containerID)
	
	// Remove
	return p.client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force: true,
	})
}

// GetContainerStatus gets the status of a container
func (p *DockerProvisioner) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	inspect, err := p.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return inspect.State.Status, nil
}

// Close closes the Docker client
func (p *DockerProvisioner) Close() error {
	return p.client.Close()
}

// Helper functions

func getOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getInternalPort(serviceType string) int {
	switch serviceType {
	case "postgres":
		return 5432
	case "redis":
		return 6379
	case "kafka":
		return 9092
	default:
		return 8080
	}
}

