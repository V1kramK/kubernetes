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

func Replaceappsv1namespacedcontrollerrevisionHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		var requestBody models.Iok8sapiappsv1ControllerRevision
		
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
		url := fmt.Sprintf("%s/apis/apps/v1/namespaces/%s/controllerrevisions/%s%s", cfg.BaseURL, queryString)
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
		var result models.Iok8sapiappsv1ControllerRevision
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

func CreateReplaceappsv1namespacedcontrollerrevisionTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("put_apis_apps_v1_namespaces_namespace_controllerrevisions_name",
		mcp.WithDescription("replace the specified ControllerRevision"),
		mcp.WithString("dryRun", mcp.Description("When present, indicates that modifications should not be persisted. An invalid or unrecognized dryRun directive will result in an error response and no further processing of the request. Valid values are: - All: all dry run stages will be processed")),
		mcp.WithString("fieldManager", mcp.Description("fieldManager is a name associated with the actor or entity that is making these changes. The value must be less than or 128 characters long, and only contain printable characters, as defined by https://golang.org/pkg/unicode/#IsPrint.")),
		mcp.WithString("fieldValidation", mcp.Description("fieldValidation instructs the server on how to handle objects in the request (POST/PUT/PATCH) containing unknown or duplicate fields. Valid values are: - Ignore: This will ignore any unknown fields that are silently dropped from the object, and will ignore all but the last duplicate field that the decoder encounters. This is the default behavior prior to v1.23. - Warn: This will send a warning via the standard warning response header for each unknown field that is dropped from the object, and for each duplicate field that is encountered. The request will still succeed if there are no other errors, and will only persist the last of any duplicate fields. This is the default in v1.23+ - Strict: This will fail the request with a BadRequest error if any unknown fields would be dropped from the object, or if any duplicate fields are present. The error returned from the server will contain all unknown and duplicate fields encountered.")),
		mcp.WithNumber("revision", mcp.Required(), mcp.Description("Input parameter: Revision indicates the revision of the state represented by Data.")),
		mcp.WithString("apiVersion", mcp.Description("Input parameter: APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources")),
		mcp.WithObject("data", mcp.Description("Input parameter: RawExtension is used to hold extensions in external versions.\n\nTo use this, make a field which has RawExtension as its type in your external, versioned struct, and Object in your internal struct. You also need to register your various plugin types.\n\n// Internal package:\n\n	type MyAPIObject struct {\n		runtime.TypeMeta `json:\",inline\"`\n		MyPlugin runtime.Object `json:\"myPlugin\"`\n	}\n\n	type PluginA struct {\n		AOption string `json:\"aOption\"`\n	}\n\n// External package:\n\n	type MyAPIObject struct {\n		runtime.TypeMeta `json:\",inline\"`\n		MyPlugin runtime.RawExtension `json:\"myPlugin\"`\n	}\n\n	type PluginA struct {\n		AOption string `json:\"aOption\"`\n	}\n\n// On the wire, the JSON will look something like this:\n\n	{\n		\"kind\":\"MyAPIObject\",\n		\"apiVersion\":\"v1\",\n		\"myPlugin\": {\n			\"kind\":\"PluginA\",\n			\"aOption\":\"foo\",\n		},\n	}\n\nSo what happens? Decode first uses json or yaml to unmarshal the serialized data into your external MyAPIObject. That causes the raw JSON to be stored, but not unpacked. The next step is to copy (using pkg/conversion) into the internal struct. The runtime package's DefaultScheme has conversion functions installed which will unpack the JSON stored in RawExtension, turning it into the correct object type, and storing it in the Object. (TODO: In the case where the object is of an unknown type, a runtime.Unknown object will be created and stored.)")),
		mcp.WithString("kind", mcp.Description("Input parameter: Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds")),
		mcp.WithObject("metadata", mcp.Description("Input parameter: ObjectMeta is metadata that all persisted resources must have, which includes all objects users must create.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Replaceappsv1namespacedcontrollerrevisionHandler(cfg),
	}
}
