package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/akshitmadan/throome/internal/logger"
	"github.com/akshitmadan/throome/internal/utils"
	"github.com/akshitmadan/throome/pkg/cluster"
	"go.uber.org/zap"
)

var (
	Version   = "0.1.0"
	BuildTime = "unknown"

	// Global flags
	clustersDir string
	verbose     bool

	// Command-specific flags
	clusterName string
	serviceName string
	serviceType string
	host        string
	port        int
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "throome-cli",
	Short: "Throome CLI - Manage Throome clusters",
	Long:  `Throome CLI is a command-line tool for managing Throome gateway clusters.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger
		if err := logger.InitLogger(verbose); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Throome CLI v%s (built: %s)\n", Version, BuildTime)
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create-cluster",
	Short: "Create a new cluster",
	Long:  `Create a new Throome cluster with the specified configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if clusterName == "" {
			fmt.Println("Error: cluster name is required (use --name)")
			os.Exit(1)
		}

		// Initialize cluster manager
		manager := cluster.NewManager(clustersDir)

		// Generate cluster ID from name
		clusterID := utils.SanitizeClusterName(clusterName)

		// Validate cluster ID
		if err := utils.ValidateClusterID(clusterID); err != nil {
			fmt.Printf("Error: invalid cluster ID generated from name: %v\n", err)
			os.Exit(1)
		}

		// Create default config
		config := cluster.DefaultConfig(clusterID, clusterName)

		// Create cluster
		createdID, err := manager.Create(clusterName, config)
		if err != nil {
			fmt.Printf("Error creating cluster: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Cluster created successfully!\n")
		fmt.Printf("  Cluster ID: %s\n", createdID)
		fmt.Printf("  Name: %s\n", clusterName)
		fmt.Printf("  Config: %s\n", filepath.Join(clustersDir, createdID, "config.yaml"))
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("  1. Edit the config file to add services\n")
		fmt.Printf("  2. Restart the Throome gateway to load the cluster\n")
	},
}

var listClustersCmd = &cobra.Command{
	Use:   "list-clusters",
	Short: "List all clusters",
	Run: func(cmd *cobra.Command, args []string) {
		manager := cluster.NewManager(clustersDir)

		clusterIDs, err := manager.List()
		if err != nil {
			fmt.Printf("Error listing clusters: %v\n", err)
			os.Exit(1)
		}

		if len(clusterIDs) == 0 {
			fmt.Println("No clusters found.")
			return
		}

		// Load cluster details
		if err := manager.LoadAll(); err != nil {
			logger.Error("Failed to load clusters", zap.Error(err))
		}

		configs := manager.GetAllConfigs()

		// Print table
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CLUSTER ID\tNAME\tSERVICES\tCREATED")
		fmt.Fprintln(w, "----------\t----\t--------\t-------")

		for _, id := range clusterIDs {
			config := configs[id]
			if config != nil {
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\n",
					config.ClusterID,
					config.Name,
					len(config.Services),
					config.CreatedAt.Format("2006-01-02"),
				)
			}
		}

		w.Flush()
	},
}

var getClusterCmd = &cobra.Command{
	Use:   "get-cluster [cluster-id]",
	Short: "Get cluster details",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterID := args[0]

		manager := cluster.NewManager(clustersDir)

		config, err := manager.Get(clusterID)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Print cluster details
		fmt.Printf("Cluster: %s\n", config.Name)
		fmt.Printf("ID: %s\n", config.ClusterID)
		fmt.Printf("Description: %s\n", config.Description)
		fmt.Printf("Created: %s\n", config.CreatedAt.Format(time.RFC3339))
		fmt.Printf("Updated: %s\n", config.UpdatedAt.Format(time.RFC3339))
		fmt.Printf("\nServices (%d):\n", len(config.Services))

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tHOST\tPORT")
		fmt.Fprintln(w, "----\t----\t----\t----")

		for name, svc := range config.Services {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", name, svc.Type, svc.Host, svc.Port)
		}

		w.Flush()

		fmt.Printf("\nRouting Strategy: %s\n", config.Routing.Strategy)
		fmt.Printf("Health Checks: %v\n", config.Health.Enabled)
		fmt.Printf("AI Optimization: %v\n", config.AI.Enabled)
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete-cluster [cluster-id]",
	Short: "Delete a cluster",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clusterID := args[0]

		manager := cluster.NewManager(clustersDir)

		// Confirm deletion
		fmt.Printf("Are you sure you want to delete cluster '%s'? (yes/no): ", clusterID)
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "yes" {
			fmt.Println("Deletion cancelled.")
			return
		}

		if err := manager.Delete(clusterID); err != nil {
			fmt.Printf("Error deleting cluster: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Cluster '%s' deleted successfully!\n", clusterID)
	},
}

var validateConfigCmd = &cobra.Command{
	Use:   "validate-config [config-file]",
	Short: "Validate a cluster configuration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		configPath := args[0]

		// Create a temporary loader
		loader := cluster.NewLoader(filepath.Dir(configPath))

		// Load and validate
		config, err := loader.Load(filepath.Base(filepath.Dir(configPath)))
		if err != nil {
			fmt.Printf("✗ Configuration is invalid: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Configuration is valid!\n")
		fmt.Printf("  Cluster ID: %s\n", config.ClusterID)
		fmt.Printf("  Services: %d\n", len(config.Services))
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&clustersDir, "clusters-dir", "./clusters", "Path to clusters directory")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Create cluster flags
	createClusterCmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	createClusterCmd.MarkFlagRequired("name")

	// Add commands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(createClusterCmd)
	rootCmd.AddCommand(listClustersCmd)
	rootCmd.AddCommand(getClusterCmd)
	rootCmd.AddCommand(deleteClusterCmd)
	rootCmd.AddCommand(validateConfigCmd)
}
