package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kubernetes/mcp-server/config"
	"github.com/kubernetes/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Watchcorev1namespacedsecretlistHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		queryParams := make([]string, 0)
		// Handle multiple authentication parameters
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("continue=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/api/v1/watch/namespaces/%s/secrets%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to create request", err), nil
		}
		// Set authentication based on auth type
		// Handle multiple authentication parameters
		// API key already added to query string
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Request failed", err), nil
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to read response body", err), nil
		}

		if resp.StatusCode >= 400 {
			return mcp.NewToolResultError(fmt.Sprintf("API error: %s", body)), nil
		}
		// Use properly typed response
		var result models.Iok8sapimachinerypkgapismetav1WatchEvent
		if err := json.Unmarshal(body, &result); err != nil {
			// Fallback to raw text if unmarshaling fails
			return mcp.NewToolResultText(string(body)), nil
		}

		prettyJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to format JSON", err), nil
		}

		return mcp.NewToolResultText(string(prettyJSON)), nil
	}
}

func CreateWatchcorev1namespacedsecretlistTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_api_v1_watch_namespaces_namespace_secrets",
		mcp.WithDescription("watch individual changes to a list of Secret. deprecated: use the 'watch' parameter with a list operation instead."),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Watchcorev1namespacedsecretlistHandler(cfg),
	}
}
