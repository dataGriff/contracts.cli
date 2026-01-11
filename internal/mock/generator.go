package mock

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dataGriff/contracts.cli/internal/contracts"
	"gopkg.in/yaml.v3"
)

// Generator generates mock data from contracts
type Generator struct{}

// NewGenerator creates a new mock generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateMock generates mock data based on the contract type
func (g *Generator) GenerateMock(contract *contracts.Contract) (string, error) {
	switch contract.Type {
	case contracts.OpenAPI:
		return g.generateOpenAPIMock(contract)
	case contracts.ODCS:
		return g.generateODCSMock(contract)
	case contracts.AsyncAPI:
		return g.generateAsyncAPIMock(contract)
	default:
		return "", fmt.Errorf("unsupported contract type: %s", contract.Type)
	}
}

func (g *Generator) generateOpenAPIMock(contract *contracts.Contract) (string, error) {
	var spec map[string]interface{}
	if err := yaml.Unmarshal([]byte(contract.Content), &spec); err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	mock := map[string]interface{}{
		"type":        "OpenAPI Mock",
		"description": "Mock data generated from OpenAPI specification",
		"endpoints":   []map[string]interface{}{},
	}

	// Extract paths
	if paths, ok := spec["paths"].(map[string]interface{}); ok {
		for path, pathItem := range paths {
			if pathObj, ok := pathItem.(map[string]interface{}); ok {
				for method := range pathObj {
					if method != "parameters" {
						endpoint := map[string]interface{}{
							"path":   path,
							"method": strings.ToUpper(method),
							"mockResponse": map[string]interface{}{
								"status": 200,
								"body":   "Mock response for " + strings.ToUpper(method) + " " + path,
							},
						}
						mock["endpoints"] = append(mock["endpoints"].([]map[string]interface{}), endpoint)
					}
				}
			}
		}
	}

	result, err := json.MarshalIndent(mock, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (g *Generator) generateODCSMock(contract *contracts.Contract) (string, error) {
	var spec map[string]interface{}
	if err := yaml.Unmarshal([]byte(contract.Content), &spec); err != nil {
		return "", fmt.Errorf("failed to parse ODCS spec: %w", err)
	}

	mock := map[string]interface{}{
		"type":        "ODCS Mock",
		"description": "Mock data generated from Open Data Contract Specification",
		"records":     []map[string]interface{}{},
	}

	// Extract models
	if models, ok := spec["models"].(map[string]interface{}); ok {
		for modelName, modelDef := range models {
			if modelObj, ok := modelDef.(map[string]interface{}); ok {
				record := map[string]interface{}{
					"model": modelName,
					"data":  map[string]interface{}{},
				}

				if fields, ok := modelObj["fields"].(map[string]interface{}); ok {
					for fieldName, fieldDef := range fields {
						if fieldObj, ok := fieldDef.(map[string]interface{}); ok {
							fieldType := "string"
							if ft, ok := fieldObj["type"].(string); ok {
								fieldType = ft
							}
							record["data"].(map[string]interface{})[fieldName] = fmt.Sprintf("mock_%s_%s", fieldType, fieldName)
						}
					}
				}

				mock["records"] = append(mock["records"].([]map[string]interface{}), record)
			}
		}
	}

	result, err := json.MarshalIndent(mock, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (g *Generator) generateAsyncAPIMock(contract *contracts.Contract) (string, error) {
	var spec map[string]interface{}
	if err := yaml.Unmarshal([]byte(contract.Content), &spec); err != nil {
		return "", fmt.Errorf("failed to parse AsyncAPI spec: %w", err)
	}

	mock := map[string]interface{}{
		"type":        "AsyncAPI Mock",
		"description": "Mock events generated from AsyncAPI specification",
		"events":      []map[string]interface{}{},
	}

	// Extract channels
	if channels, ok := spec["channels"].(map[string]interface{}); ok {
		for channelName, channelDef := range channels {
			if channelObj, ok := channelDef.(map[string]interface{}); ok {
				event := map[string]interface{}{
					"channel": channelName,
				}

				if desc, ok := channelObj["description"].(string); ok {
					event["description"] = desc
				}

				// Try to extract message schema
				if publish, ok := channelObj["publish"].(map[string]interface{}); ok {
					event["operation"] = "publish"
					if message, ok := publish["message"].(map[string]interface{}); ok {
						event["message"] = message
					}
				}

				if subscribe, ok := channelObj["subscribe"].(map[string]interface{}); ok {
					event["operation"] = "subscribe"
					if message, ok := subscribe["message"].(map[string]interface{}); ok {
						event["message"] = message
					}
				}

				event["mockPayload"] = map[string]interface{}{
					"timestamp": time.Now().Format(time.RFC3339),
					"data":      "Mock event data for " + channelName,
				}

				mock["events"] = append(mock["events"].([]map[string]interface{}), event)
			}
		}
	}

	result, err := json.MarshalIndent(mock, "", "  ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}
