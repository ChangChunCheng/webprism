package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"

	webprismv1 "github.com/ChangChunCheng/webprism/gen/go/v1"
)

// ==================== AUTH COMMAND ====================

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication configurations",
}

var authSetCmd = &cobra.Command{
	Use:   "set <spec-id>",
	Short: "Set authentication configuration for an API",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuthSet,
}

var authGetCmd = &cobra.Command{
	Use:   "get <spec-id>",
	Short: "Get authentication configuration",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuthGet,
}

func init() {
	authCmd.AddCommand(authSetCmd)
	authCmd.AddCommand(authGetCmd)

	authSetCmd.Flags().StringP("type", "t", "bearer", "Auth type (api_key, bearer, basic, oauth2)")
	authSetCmd.Flags().StringToStringP("credentials", "c", nil, "Credentials as key=value pairs")
	_ = authSetCmd.MarkFlagRequired("credentials")
}

func runAuthSet(cmd *cobra.Command, args []string) error {
	specID := args[0]
	authType, _ := cmd.Flags().GetString("type")
	credentials, _ := cmd.Flags().GetStringToString("credentials")

	conn, err := grpc.NewClient(getServerAddress(cmd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := webprismv1.NewAuthServiceClient(conn)

	// Convert auth type string to enum
	var authTypeEnum webprismv1.AuthType
	switch authType {
	case "api_key":
		authTypeEnum = webprismv1.AuthType_AUTH_TYPE_API_KEY
	case "bearer":
		authTypeEnum = webprismv1.AuthType_AUTH_TYPE_BEARER
	case "basic":
		authTypeEnum = webprismv1.AuthType_AUTH_TYPE_BASIC
	case "oauth2":
		authTypeEnum = webprismv1.AuthType_AUTH_TYPE_OAUTH2
	default:
		return fmt.Errorf("invalid auth type: %s", authType)
	}

	_, err = client.SetAuthConfig(context.Background(), &webprismv1.SetAuthConfigRequest{
		SpecId:      specID,
		AuthType:    authTypeEnum,
		Credentials: credentials,
	})
	if err != nil {
		return fmt.Errorf("failed to set auth config: %w", err)
	}

	printSuccess(fmt.Sprintf("Authentication configured for spec: %s", specID))
	return nil
}

func runAuthGet(cmd *cobra.Command, args []string) error {
	specID := args[0]

	conn, err := grpc.NewClient(getServerAddress(cmd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := webprismv1.NewAuthServiceClient(conn)
	resp, err := client.GetAuthConfig(context.Background(), &webprismv1.GetAuthConfigRequest{
		SpecId: specID,
	})
	if err != nil {
		return fmt.Errorf("failed to get auth config: %w", err)
	}

	result, _ := json.MarshalIndent(resp, "", "  ")
	printJSON(string(result))
	return nil
}

// ==================== PROXY COMMAND ====================

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Execute proxied API requests",
}

var proxyCallCmd = &cobra.Command{
	Use:   "call <spec-id> <operation-id>",
	Short: "Execute a proxied API call",
	Args:  cobra.ExactArgs(2),
	RunE:  runProxyCall,
}

func init() {
	proxyCmd.AddCommand(proxyCallCmd)

	proxyCallCmd.Flags().StringToStringP("param", "p", nil, "Request parameters as key=value pairs")
	proxyCallCmd.Flags().StringP("body", "b", "", "Request body as JSON string")
}

func runProxyCall(cmd *cobra.Command, args []string) error {
	specID := args[0]
	operationID := args[1]
	params, _ := cmd.Flags().GetStringToString("param")
	bodyStr, _ := cmd.Flags().GetString("body")

	conn, err := grpc.NewClient(getServerAddress(cmd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := webprismv1.NewProxyServiceClient(conn)

	// Parse body
	var bodyMap map[string]interface{}
	if bodyStr != "" {
		if err := json.Unmarshal([]byte(bodyStr), &bodyMap); err != nil {
			return fmt.Errorf("invalid body JSON: %w", err)
		}
	}

	bodyStruct, _ := structpb.NewStruct(bodyMap)

	resp, err := client.ExecuteProxy(context.Background(), &webprismv1.ExecuteProxyRequest{
		SpecId:      specID,
		OperationId: operationID,
		Parameters:  params,
		Body:        bodyStruct,
	})
	if err != nil {
		return fmt.Errorf("failed to execute proxy: %w", err)
	}

	result, _ := json.MarshalIndent(resp, "", "  ")
	printJSON(string(result))
	return nil
}

// ==================== HEALTH COMMAND ====================

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check API health status",
}

var healthCheckCmd = &cobra.Command{
	Use:   "check <spec-id>",
	Short: "Perform a health check",
	Args:  cobra.ExactArgs(1),
	RunE:  runHealthCheck,
}

var healthGetCmd = &cobra.Command{
	Use:   "get <spec-id>",
	Short: "Get health check history",
	Args:  cobra.ExactArgs(1),
	RunE:  runHealthGet,
}

func init() {
	healthCmd.AddCommand(healthCheckCmd)
	healthCmd.AddCommand(healthGetCmd)

	healthGetCmd.Flags().IntP("limit", "l", 10, "Maximum number of health checks to return")
}

func runHealthCheck(cmd *cobra.Command, args []string) error {
	specID := args[0]

	conn, err := grpc.NewClient(getServerAddress(cmd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := webprismv1.NewHealthServiceClient(conn)
	resp, err := client.CheckHealth(context.Background(), &webprismv1.CheckHealthRequest{
		SpecId: specID,
	})
	if err != nil {
		return fmt.Errorf("failed to check health: %w", err)
	}

	result, _ := json.MarshalIndent(resp.Result, "", "  ")
	printJSON(string(result))
	return nil
}

func runHealthGet(cmd *cobra.Command, args []string) error {
	specID := args[0]
	limit, _ := cmd.Flags().GetInt("limit")

	conn, err := grpc.NewClient(getServerAddress(cmd), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := webprismv1.NewHealthServiceClient(conn)
	resp, err := client.GetHealthHistory(context.Background(), &webprismv1.GetHealthHistoryRequest{
		SpecId: specID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get health history: %w", err)
	}

	result, _ := json.MarshalIndent(resp.Results, "", "  ")
	printJSON(string(result))
	return nil
}

// ==================== VERSION COMMAND ====================

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WEBPRISM v1.0.0")
	},
}
