package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"bytes"

	"github.com/kubernetes/mcp-server/config"
	"github.com/kubernetes/mcp-server/models"
	"github.com/mark3labs/mcp-go/mcp"
)

func Replaceeventsv1namespacedeventHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["dryRun"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("dryRun=%v", val))
		}
		if val, ok := args["fieldManager"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("fieldManager=%v", val))
		}
		if val, ok := args["fieldValidation"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("fieldValidation=%v", val))
		}
		// Handle multiple authentication parameters
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("continue=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		// Create properly typed request body using the generated schema
		var requestBody models.Iok8sapieventsv1Event
		
		// Optimized: Single marshal/unmarshal with JSON tags handling field mapping
		if argsJSON, err := json.Marshal(args); err == nil {
			if err := json.Unmarshal(argsJSON, &requestBody); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to convert arguments to request type: %v", err)), nil
			}
		} else {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal arguments: %v", err)), nil
		}
		
		bodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("Failed to encode request body", err), nil
		}
		url := fmt.Sprintf("%s/apis/events.k8s.io/v1/namespaces/%s/events/%s%s", cfg.BaseURL, queryString)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
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
		var result models.Iok8sapieventsv1Event
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

func CreateReplaceeventsv1namespacedeventTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_apis_events_k8s_io_v1_namespaces_namespace_events_name",
		mcp.WithDescription("replace the specified Event"),
		mcp.WithString("dryRun", mcp.Description("When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed")),
		mcp.WithString("fieldManager", mcp.Description("fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.")),
		mcp.WithString("fieldValidation", mcp.Description("fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered.")),
		mcp.WithObject("metadata", mcp.Description("Input parameter: ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.")),
		mcp.WithString("eventTime", mcp.Required(), mcp.Description("Input parameter: MicroTime is version of Time with microsecond level precision.")),
		mcp.WithString("note", mcp.Description("Input parameter: note is a human-readable description of the status of this operation. Maximal length of the note is 1kB, but libraries should be prepared to handle values up to 64kB.")),
		mcp.WithString("reportingInstance", mcp.Description("Input parameter: reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`. This field cannot be empty for new Events and it can have at most 128 characters.")),
		mcp.WithNumber("deprecatedCount", mcp.Description("Input parameter: deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.")),
		mcp.WithString("kind", mcp.Description("Input parameter: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds")),
		mcp.WithObject("regarding", mcp.Description("Input parameter: ObjectReference contains enough information to let you inspect or modify the referred object.")),
		mcp.WithObject("series", mcp.Description("Input parameter: EventSeries contain information on series of events, i.e. thing that was/is happening continuously for some time. How often to update the EventSeries is up to the event reporters. The default event reporter in \"k8s.io/client-go/tools/events/event_broadcaster.go\" shows how this struct is updated on heartbeats and can guide customized reporter implementations.")),
		mcp.WithObject("deprecatedSource", mcp.Description("Input parameter: EventSource contains information for an event.")),
		mcp.WithString("deprecatedFirstTimestamp", mcp.Description("Input parameter: Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.")),
		mcp.WithString("apiVersion", mcp.Description("Input parameter: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources")),
		mcp.WithString("reportingController", mcp.Description("Input parameter: reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`. This field cannot be empty for new Events.")),
		mcp.WithString("deprecatedLastTimestamp", mcp.Description("Input parameter: Time is a wrapper around time.Time which supports correct marshaling to YAML and JSON.  Wrappers are provided for many of the factory methods that the time package offers.")),
		mcp.WithString("action", mcp.Description("Input parameter: action is what action was taken/failed regarding to the regarding object. It is machine-readable. This field cannot be empty for new Events and it can have at most 128 characters.")),
		mcp.WithObject("related", mcp.Description("Input parameter: ObjectReference contains enough information to let you inspect or modify the referred object.")),
		mcp.WithString("type", mcp.Description("Input parameter: type is the type of this event (Normal, Warning), new types could be added in the future. It is machine-readable. This field cannot be empty for new Events.")),
		mcp.WithString("reason", mcp.Description("Input parameter: reason is why the action was taken. It is human-readable. This field cannot be empty for new Events and it can have at most 128 characters.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Replaceeventsv1namespacedeventHandler(cfg),
	}
}
