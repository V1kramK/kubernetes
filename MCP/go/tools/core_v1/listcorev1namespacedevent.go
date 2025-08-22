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

func Listcorev1namespacedeventHandler(cfg *config.APIConfig) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, ok := request.Params.Arguments.(map[string]any)
		if !ok {
			return mcp.NewToolResultError("Invalid arguments object"), nil
		}
		queryParams := make([]string, 0)
		if val, ok := args["allowWatchBookmarks"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("allowWatchBookmarks=%v", val))
		}
		if val, ok := args["continue"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("continue=%v", val))
		}
		if val, ok := args["fieldSelector"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("fieldSelector=%v", val))
		}
		if val, ok := args["labelSelector"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("labelSelector=%v", val))
		}
		if val, ok := args["limit"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("limit=%v", val))
		}
		if val, ok := args["resourceVersion"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("resourceVersion=%v", val))
		}
		if val, ok := args["resourceVersionMatch"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("resourceVersionMatch=%v", val))
		}
		if val, ok := args["sendInitialEvents"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("sendInitialEvents=%v", val))
		}
		if val, ok := args["timeoutSeconds"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("timeoutSeconds=%v", val))
		}
		if val, ok := args["watch"]; ok {
			queryParams = append(queryParams, fmt.Sprintf("watch=%v", val))
		}
		// Handle multiple authentication parameters
		if cfg.APIKey != "" {
			queryParams = append(queryParams, fmt.Sprintf("continue=%s", cfg.APIKey))
		}
		queryString := ""
		if len(queryParams) > 0 {
			queryString = "?" + strings.Join(queryParams, "&")
		}
		url := fmt.Sprintf("%s/api/v1/namespaces/%s/events%s", cfg.BaseURL, queryString)
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
		var result models.Iok8sapicorev1EventList
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

func CreateListcorev1namespacedeventTool(cfg *config.APIConfig) models.Tool {
	tool := mcp.NewTool("get_api_v1_namespaces_namespace_events",
		mcp.WithDescription("list or watch objects of kind Event"),
		mcp.WithBoolean("allowWatchBookmarks", mcp.Description("allowWatchBookmarks requests watch events with type \"BOOKMARK\". Servers that do not implement bookmarks may ignore this flag and bookmarks are sent at the server's discretion. Clients should not assume bookmarks are returned at any specific interval, nor may they assume the server will send any BOOKMARK event during a session. If this is not a watch, this field is ignored.")),
		mcp.WithString("continue", mcp.Description("The continue option should be set when retrieving more results from the server. Since this value is server defined, clients may only use the continue value from a previous query result with identical query parameters (except for the value of continue) and the server may reject a continue value it does not recognize. If the specified continue value is no longer valid whether due to expiration (generally five to fifteen minutes) or a configuration change on the server, the server will respond with a 410 ResourceExpired error together with a continue token. If the client needs a consistent list, it must restart their list without the continue field. Otherwise, the client may send another list request with the token received with the 410 error, the server will respond with a list starting from the next key, but from the latest snapshot, which is inconsistent from the previous list results - objects that are created, modified, or deleted after the first list request will be included in the response, as long as their keys are after the \"next key\".\n\nThis field is not supported when watch is true. Clients may start a watch from the last resourceVersion value returned by the server and not miss any modifications.")),
		mcp.WithString("fieldSelector", mcp.Description("A selector to restrict the list of returned objects by their fields. Defaults to everything.")),
		mcp.WithString("labelSelector", mcp.Description("A selector to restrict the list of returned objects by their labels. Defaults to everything.")),
		mcp.WithNumber("limit", mcp.Description("limit is a maximum number of responses to return for a list call. If more items exist, the server will set the `continue` field on the list metadata to a value that can be used with the same initial query to retrieve the next set of results. Setting a limit may return fewer than the requested amount of items (up to zero items) in the event all requested objects are filtered out and clients should only use the presence of the continue field to determine whether more results are available. Servers may choose not to support the limit argument and will return all of the available results. If limit is specified and the continue field is empty, clients may assume that no more results are available. This field is not supported if watch is true.\n\nThe server guarantees that the objects returned when using continue will be identical to issuing a single list call without a limit - that is, no objects created, modified, or deleted after the first request is issued will be included in any subsequent continued requests. This is sometimes referred to as a consistent snapshot, and ensures that a client that is using limit to receive smaller chunks of a very large result can ensure they see all possible objects. If objects are updated during a chunked list the version of the object that was present at the time the first list result was calculated is returned.")),
		mcp.WithString("resourceVersion", mcp.Description("resourceVersion sets a constraint on what resource versions a request may be served from. See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.\n\nDefaults to unset")),
		mcp.WithString("resourceVersionMatch", mcp.Description("resourceVersionMatch determines how resourceVersion is applied to list calls. It is highly recommended that resourceVersionMatch be set for list calls where resourceVersion is set See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for details.\n\nDefaults to unset")),
		mcp.WithBoolean("sendInitialEvents", mcp.Description("`sendInitialEvents=true` may be set together with `watch=true`. In that case, the watch stream will begin with synthetic events to produce the current state of objects in the collection. Once all such events have been sent, a synthetic \"Bookmark\" event  will be sent. The bookmark will report the ResourceVersion (RV) corresponding to the set of objects, and be marked with `\"k8s.io/initial-events-end\": \"true\"` annotation. Afterwards, the watch stream will proceed as usual, sending watch events corresponding to changes (subsequent to the RV) to objects watched.\n\nWhen `sendInitialEvents` option is set, we require `resourceVersionMatch` option to also be set. The semantic of the watch request is as following: - `resourceVersionMatch` = NotOlderThan\n  is interpreted as \"data at least as new as the provided `resourceVersion`\"\n  and the bookmark event is send when the state is synced\n  to a `resourceVersion` at least as fresh as the one provided by the ListOptions.\n  If `resourceVersion` is unset, this is interpreted as \"consistent read\" and the\n  bookmark event is send when the state is synced at least to the moment\n  when request started being processed.\n- `resourceVersionMatch` set to any other value or unset\n  Invalid error is returned.\n\nDefaults to true if `resourceVersion=\"\"` or `resourceVersion=\"0\"` (for backward compatibility reasons) and to false otherwise.")),
		mcp.WithNumber("timeoutSeconds", mcp.Description("Timeout for the list/watch call. This limits the duration of the call, regardless of any activity or inactivity.")),
		mcp.WithBoolean("watch", mcp.Description("Watch for changes to the described resources and return them as a stream of add, update, and remove notifications. Specify resourceVersion.")),
	)

	return models.Tool{
		Definition: tool,
		Handler:    Listcorev1namespacedeventHandler(cfg),
	}
}
