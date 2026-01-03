package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/ChangChunCheng/webprism/internal/infrastructure/config"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "webprism",
	Short: "WEBPRISM - Web Proxy Request Integration Service Manager",
	Long: `WEBPRISM is an intelligent API integration middleware that accepts OpenAPI/Swagger
specifications and creates proxy interfaces automatically.

Features:
  - Upload and manage OpenAPI 2.0 and 3.0 specifications
  - Configure authentication (API Key, Bearer Token, Basic, OAuth2)
  - Execute proxied API requests
  - Health check monitoring
  - Multiple interfaces: HTTP REST, gRPC, CLI, MCP Server`,
	PersistentPreRunE: initializeConfig,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./deployments/config.yaml)")
	rootCmd.PersistentFlags().StringP("server", "s", "localhost:9090", "WEBPRISM gRPC server address (host:port)")

	// Add subcommands
	rootCmd.AddCommand(specCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(proxyCmd)
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(versionCmd)
}

// initializeConfig initializes the configuration.
func initializeConfig(cmd *cobra.Command, args []string) error {
	// Load config (CLI can work without config file)
	_, _ = config.Load(cfgFile)
	return nil
}

// getServerAddress returns the server address from flags.
func getServerAddress(cmd *cobra.Command) string {
	server, _ := cmd.Flags().GetString("server")
	// Remove http:// or https:// prefix for gRPC connection
	server = strings.TrimPrefix(server, "http://")
	server = strings.TrimPrefix(server, "https://")
	return server
}

// printJSON prints a JSON string to stdout.
func printJSON(data string) {
	fmt.Println(data)
}

// printSuccess prints a success message to stdout.
func printSuccess(msg string) {
	fmt.Println(msg)
}
