package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	webprismv1 "github.com/ChangChunCheng/webprism/gen/go/v1"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Manage API specifications",
	Long:  "Upload, list, get, and delete OpenAPI specifications",
}

var specUploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload an OpenAPI specification",
	Args:  cobra.ExactArgs(1),
	RunE:  runSpecUpload,
}

var specListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all API specifications",
	RunE:  runSpecList,
}

var specGetCmd = &cobra.Command{
	Use:   "get <spec-id>",
	Short: "Get API specification details",
	Args:  cobra.ExactArgs(1),
	RunE:  runSpecGet,
}

var specDeleteCmd = &cobra.Command{
	Use:   "delete <spec-id>",
	Short: "Delete an API specification",
	Args:  cobra.ExactArgs(1),
	RunE:  runSpecDelete,
}

func init() {
	specCmd.AddCommand(specUploadCmd)
	specCmd.AddCommand(specListCmd)
	specCmd.AddCommand(specGetCmd)
	specCmd.AddCommand(specDeleteCmd)

	// Flags for upload
	specUploadCmd.Flags().StringP("name", "n", "", "API name")
	_ = specUploadCmd.MarkFlagRequired("name")

	// Flags for list
	specListCmd.Flags().IntP("limit", "l", 10, "Maximum number of specs to return")
}

func runSpecUpload(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	name, _ := cmd.Flags().GetString("name")

	// Read spec file
	specData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read spec file: %w", err)
	}

	// Determine format (assume JSON for now)
	format := webprismv1.SpecFormat_SPEC_FORMAT_JSON

	// Connect to gRPC server
	client, conn, err := connectSpecService(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Upload spec
	resp, err := client.UploadSpec(context.Background(), &webprismv1.UploadSpecRequest{
		Name:        name,
		SpecContent: string(specData),
		Format:      format,
	})
	if err != nil {
		return fmt.Errorf("failed to upload spec: %w", err)
	}

	// Print result
	result, _ := json.MarshalIndent(resp.Spec, "", "  ")
	printSuccess(fmt.Sprintf("API spec uploaded successfully:\n%s", string(result)))

	return nil
}

func runSpecList(cmd *cobra.Command, args []string) error {
	limit, _ := cmd.Flags().GetInt("limit")

	// Connect to gRPC server
	client, conn, err := connectSpecService(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// List specs
	resp, err := client.ListSpecs(context.Background(), &webprismv1.ListSpecsRequest{
		PageSize: int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to list specs: %w", err)
	}

	// Print result
	result, _ := json.MarshalIndent(resp.Specs, "", "  ")
	printJSON(string(result))

	return nil
}

func runSpecGet(cmd *cobra.Command, args []string) error {
	specID := args[0]

	// Connect to gRPC server
	client, conn, err := connectSpecService(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get spec
	resp, err := client.GetSpec(context.Background(), &webprismv1.GetSpecRequest{
		Id: specID,
	})
	if err != nil {
		return fmt.Errorf("failed to get spec: %w", err)
	}

	// Print result
	result, _ := json.MarshalIndent(resp.Spec, "", "  ")
	printJSON(string(result))

	return nil
}

func runSpecDelete(cmd *cobra.Command, args []string) error {
	specID := args[0]

	// Connect to gRPC server
	client, conn, err := connectSpecService(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Delete spec
	_, err = client.DeleteSpec(context.Background(), &webprismv1.DeleteSpecRequest{
		Id: specID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete spec: %w", err)
	}

	printSuccess(fmt.Sprintf("API spec deleted successfully: %s", specID))

	return nil
}

func connectSpecService(cmd *cobra.Command) (webprismv1.SpecServiceClient, *grpc.ClientConn, error) {
	serverAddr := getServerAddress(cmd)

	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := webprismv1.NewSpecServiceClient(conn)
	return client, conn, nil
}
