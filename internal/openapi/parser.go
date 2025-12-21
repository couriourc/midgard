package openapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/midgard/gateway/internal/database"
)

// OpenAPISpec represents the OpenAPI specification structure
type OpenAPISpec struct {
	OpenAPI string                 `json:"openapi"`
	Info    map[string]interface{} `json:"info"`
	Paths   map[string]interface{} `json:"paths"`
	Servers []map[string]string    `json:"servers"`
}

// ParseOpenAPIFromURL fetches and parses OpenAPI spec from URL
func ParseOpenAPIFromURL(url string) (*OpenAPISpec, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch OpenAPI spec: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return ParseOpenAPIFromJSON(body)
}

// ParseOpenAPIFromJSON parses OpenAPI spec from JSON bytes
func ParseOpenAPIFromJSON(data []byte) (*OpenAPISpec, error) {
	var spec OpenAPISpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI JSON: %w", err)
	}

	return &spec, nil
}

// ExtractEndpoints extracts endpoints from OpenAPI spec
func ExtractEndpoints(spec *OpenAPISpec, baseURL string) ([]database.Endpoint, error) {
	var endpoints []database.Endpoint

	// Determine base URL
	if baseURL == "" && len(spec.Servers) > 0 {
		baseURL = spec.Servers[0]["url"]
	}

	for path, pathItem := range spec.Paths {
		pathMap, ok := pathItem.(map[string]interface{})
		if !ok {
			continue
		}

		// Common HTTP methods
		methods := []string{"get", "post", "put", "delete", "patch", "head", "options"}

		for _, method := range methods {
			if operation, exists := pathMap[method]; exists {
				opMap, ok := operation.(map[string]interface{})
				if !ok {
					continue
				}

				endpoint := database.Endpoint{
					Path:   path,
					Method: strings.ToUpper(method),
				}

				// Extract summary
				if summary, ok := opMap["summary"].(string); ok {
					endpoint.Summary = summary
				}

				// Extract description
				if description, ok := opMap["description"].(string); ok {
					endpoint.Description = description
				}

				endpoints = append(endpoints, endpoint)
			}
		}
	}

	return endpoints, nil
}

// ImportOpenAPI imports OpenAPI spec and returns endpoints
func ImportOpenAPI(openAPIURL string, openAPIJSON []byte, baseURL string) ([]database.Endpoint, error) {
	var spec *OpenAPISpec
	var err error

	if openAPIURL != "" {
		spec, err = ParseOpenAPIFromURL(openAPIURL)
	} else if len(openAPIJSON) > 0 {
		spec, err = ParseOpenAPIFromJSON(openAPIJSON)
	} else {
		return nil, fmt.Errorf("either openAPIURL or openAPIJSON must be provided")
	}

	if err != nil {
		return nil, err
	}

	return ExtractEndpoints(spec, baseURL)
}

