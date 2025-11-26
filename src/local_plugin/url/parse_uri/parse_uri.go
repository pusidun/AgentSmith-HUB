package parse_uri

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// Eval parses URI string and returns JSON object with path and query parameters.
// Args: uri string.
// Returns: map[string]interface{} containing parsed URI components.
func Eval(args ...interface{}) (interface{}, bool, error) {
	if len(args) != 1 {
		return nil, false, fmt.Errorf("parseURI requires exactly one argument")
	}

	uriStr, ok := args[0].(string)
	if !ok {
		return nil, false, fmt.Errorf("argument must be a string")
	}

	// Parse the URI
	parsedURI, err := url.Parse(uriStr)
	if err != nil {
		return nil, false, fmt.Errorf("failed to parse URI: %w", err)
	}

	// Build result map
	result := make(map[string]interface{})

	// Add path components
	if parsedURI.Path != "" {
		// Split path into segments
		pathSegments := strings.Split(strings.Trim(parsedURI.Path, "/"), "/")
		// Filter out empty segments
		filteredSegments := make([]string, 0, len(pathSegments))
		for _, seg := range pathSegments {
			if seg != "" {
				filteredSegments = append(filteredSegments, seg)
			}
		}
		result["path"] = parsedURI.Path
		result["pathSegments"] = filteredSegments
	}

	// Add query parameters
	if parsedURI.RawQuery != "" {
		queryValues := parsedURI.Query()
		queryMap := make(map[string]interface{})
		
		// Convert url.Values to map, handling multiple values for same key
		for key, values := range queryValues {
			if len(values) == 1 {
				// Single value - store as string
				queryMap[key] = values[0]
			} else {
				// Multiple values - store as array
				queryMap[key] = values
			}
		}
		result["query"] = queryMap
		result["rawQuery"] = parsedURI.RawQuery
	}

	// Add other URI components if present
	if parsedURI.Scheme != "" {
		result["scheme"] = parsedURI.Scheme
	}
	if parsedURI.Host != "" {
		result["host"] = parsedURI.Host
		result["hostname"] = parsedURI.Hostname()
		if parsedURI.Port() != "" {
			result["port"] = parsedURI.Port()
		}
	}
	if parsedURI.Fragment != "" {
		result["fragment"] = parsedURI.Fragment
	}
	if parsedURI.User != nil {
		userInfo := make(map[string]interface{})
		userInfo["username"] = parsedURI.User.Username()
		if password, ok := parsedURI.User.Password(); ok {
			userInfo["password"] = password
		}
		result["user"] = userInfo
	}

	// Convert to JSON string for consistency, then parse back to interface{}
	// This ensures proper JSON structure
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, false, fmt.Errorf("failed to marshal result to JSON: %w", err)
	}

	var jsonResult interface{}
	err = json.Unmarshal(jsonBytes, &jsonResult)
	if err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal JSON result: %w", err)
	}

	return jsonResult, true, nil
}

