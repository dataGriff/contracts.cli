package mock

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dataGriff/contracts.cli/internal/contracts"
	"github.com/go-gota/gota/dataframe"
	"gopkg.in/yaml.v3"
)

// Generator generates mock data from contracts
type Generator struct{}

// NewGenerator creates a new mock generator
func NewGenerator() *Generator {
	rand.Seed(time.Now().UnixNano())
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

// GenerateMockDataFrame generates mock data as a dataframe for ODCS contracts
func (g *Generator) GenerateMockDataFrame(contract *contracts.Contract, rows int) (dataframe.DataFrame, error) {
	if contract.Type != contracts.ODCS {
		return dataframe.DataFrame{}, fmt.Errorf("dataframe generation only supported for ODCS contracts")
	}
	return g.generateODCSDataFrame(contract, rows)
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

// generateODCSDataFrame generates a dataframe with mock data from an ODCS contract
func (g *Generator) generateODCSDataFrame(contract *contracts.Contract, rows int) (dataframe.DataFrame, error) {
	var spec map[string]interface{}
	if err := yaml.Unmarshal([]byte(contract.Content), &spec); err != nil {
		return dataframe.DataFrame{}, fmt.Errorf("failed to parse ODCS spec: %w", err)
	}

	// Extract models
	if models, ok := spec["models"].(map[string]interface{}); ok {
		for _, modelDef := range models {
			if modelObj, ok := modelDef.(map[string]interface{}); ok {
				// Create CSV data in memory
				var buf bytes.Buffer
				writer := csv.NewWriter(&buf)

				// Get fields to determine columns
				var headers []string
				var fieldTypes []string
				if fields, ok := modelObj["fields"].(map[string]interface{}); ok {
					for fieldName, fieldDef := range fields {
						headers = append(headers, fieldName)
						if fieldObj, ok := fieldDef.(map[string]interface{}); ok {
							fieldType := "string"
							if ft, ok := fieldObj["type"].(string); ok {
								fieldType = ft
							}
							fieldTypes = append(fieldTypes, fieldType)
						} else {
							fieldTypes = append(fieldTypes, "string")
						}
					}
				}

				// Write headers
				if err := writer.Write(headers); err != nil {
					return dataframe.DataFrame{}, err
				}

				// Generate mock rows
				for i := 0; i < rows; i++ {
					row := make([]string, len(headers))
					for j, header := range headers {
						row[j] = g.generateMockValue(header, fieldTypes[j], i)
					}
					if err := writer.Write(row); err != nil {
						return dataframe.DataFrame{}, err
					}
				}

				writer.Flush()

				// Create dataframe from CSV
				df := dataframe.ReadCSV(strings.NewReader(buf.String()))
				return df, nil
			}
		}
	}

	return dataframe.DataFrame{}, fmt.Errorf("no models found in ODCS spec")
}

// generateMockValue generates a mock value based on the field type
func (g *Generator) generateMockValue(fieldName, fieldType string, index int) string {
	switch strings.ToLower(fieldType) {
	case "uuid":
		return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
			rand.Uint32(), rand.Uint32()&0xffff, rand.Uint32()&0xffff,
			rand.Uint32()&0xffff, rand.Uint64()&0xffffffffffff)
	case "int", "integer":
		return fmt.Sprintf("%d", rand.Intn(1000)+1)
	case "decimal", "float", "double":
		return fmt.Sprintf("%.2f", rand.Float64()*1000)
	case "boolean", "bool":
		return fmt.Sprintf("%t", rand.Intn(2) == 1)
	case "timestamp", "datetime":
		randomTime := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)
		return randomTime.Format(time.RFC3339)
	case "date":
		randomTime := time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour)
		return randomTime.Format("2006-01-02")
	case "string":
		if strings.Contains(strings.ToLower(fieldName), "status") {
			statuses := []string{"pending", "confirmed", "shipped", "delivered", "cancelled"}
			return statuses[rand.Intn(len(statuses))]
		}
		if strings.Contains(strings.ToLower(fieldName), "email") {
			return fmt.Sprintf("user%d@example.com", index+1)
		}
		if strings.Contains(strings.ToLower(fieldName), "name") {
			return fmt.Sprintf("Name_%d", index+1)
		}
		return fmt.Sprintf("mock_%s_%d", fieldName, index+1)
	default:
		return fmt.Sprintf("mock_%s_%d", fieldName, index+1)
	}
}
