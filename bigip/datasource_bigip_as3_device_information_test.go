package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
)

// filterApplications filters applications from AS3 JSON based on the application list.
func filterApplications(as3Resp string, applicationList []string) (string, error) {
	// Unmarshal the AS3 JSON response into a map for filtering
	as3Json := make(map[string]interface{})
	if err := json.Unmarshal([]byte(as3Resp), &as3Json); err != nil {
		// Log parsing error and return immediately
		log.Printf("[ERROR] Failed to parse AS3 JSON response: %v", err)
		return "", fmt.Errorf("failed to parse AS3 JSON response: %w", err)
	}

	// Ensure the response contains the "declaration" element
	declaration, ok := as3Json["declaration"].(map[string]interface{})
	if !ok {
		log.Printf("[ERROR] Missing 'declaration' in AS3 response")
		return "", fmt.Errorf("missing 'declaration' in AS3 response")
	}

	// Perform application filtering
	filteredAs3Json := make(map[string]interface{})
	for tenantName, tenant := range declaration {
		tenantMap, ok := tenant.(map[string]interface{})
		if !ok {
			log.Printf("[WARN] Tenant '%s' is not a valid object, skipping", tenantName)
			continue
		}

		// Look for applications inside the tenant
		for _, appName := range applicationList {
			if app, exists := tenantMap[appName]; exists {
				log.Printf("[INFO] Application '%s' found in tenant '%s'", appName, tenantName)
				if _, ok := filteredAs3Json[tenantName]; !ok {
					filteredAs3Json[tenantName] = make(map[string]interface{})
				}
				filteredAs3Json[tenantName].(map[string]interface{})[appName] = app
			} else {
				log.Printf("[WARN] Application '%s' not found in tenant '%s'", appName, tenantName)
			}
		}
	}

	// Marshal the filtered JSON back to a string
	filteredJsonBytes, err := json.Marshal(filteredAs3Json)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal filtered AS3 JSON: %v", err)
		return "", fmt.Errorf("failed to process AS3 configuration: %w", err)
	}

	// Log the filtered output and return
	log.Printf("[INFO] Filtered AS3 JSON: %s", string(filteredJsonBytes))
	return string(filteredJsonBytes), nil
}

// -- Unit Test for the filterApplications Function --
func TestFilterApplications(t *testing.T) {
	tests := []struct {
		name            string
		as3Response     string
		applicationList []string
		expectedResult  string
		expectError     bool
	}{
		// Valid JSON with one application
		{
			name: "Single application exists",
			as3Response: `{
                "action": "deploy",
                "class": "AS3",
                "declaration": {
                    "ansible": {
                        "A1": {
                            "class": "Application",
                            "template": "http",
                            "serviceMain": {
                                "class": "Service_HTTP",
                                "virtualAddresses": ["10.1.2.3"],
                                "virtualPort": 80
                            }
                        },
                        "class": "Tenant"
                    }
                }
            }`,
			applicationList: []string{"A1"},
			expectedResult:  `{"ansible":{"A1":{"class":"Application","template":"http","serviceMain":{"class":"Service_HTTP","virtualAddresses":["10.1.2.3"],"virtualPort":80}}}}`,
			expectError:     false,
		},
		// Valid JSON, but application does not exist
		{
			name: "Application does not exist",
			as3Response: `{
                "action": "deploy",
                "class": "AS3",
                "declaration": {
                    "ansible": {
                        "class": "Tenant"
                    }
                }
            }`,
			applicationList: []string{"A2"},
			expectedResult:  `{}`,
			expectError:     false,
		},
		// Malformed JSON response
		{
			name:            "Malformed JSON",
			as3Response:     `{"action": "deploy", "class": "AS3", "declaration": `,
			applicationList: []string{"A1"},
			expectedResult:  "",
			expectError:     true,
		},
		// Empty declaration key
		{
			name:            "Empty declaration",
			as3Response:     `{"action": "deploy", "class": "AS3", "declaration": {}}`,
			applicationList: []string{"A1"},
			expectedResult:  `{}`,
			expectError:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call `filterApplications`
			result, err := filterApplications(tc.as3Response, tc.applicationList)

			// Check if an error was expected
			if (err != nil) != tc.expectError {
				t.Fatalf("Expected error: %v, got: %v", tc.expectError, err)
			}

			// Compare the output only when no error is expected
			if !tc.expectError && !equalJSON(tc.expectedResult, result) {
				t.Errorf("Expected result:\n%s\nGot:\n%s", tc.expectedResult, result)
			}
		})
	}
}

// Helper function for order-independent JSON comparison
func equalJSON(expected, actual string) bool {
	var expectedMap map[string]interface{}
	var actualMap map[string]interface{}

	// Unmarshal both JSON strings
	if err := json.Unmarshal([]byte(expected), &expectedMap); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(actual), &actualMap); err != nil {
		return false
	}

	// Compare the resulting maps (order-independent check)
	return reflect.DeepEqual(expectedMap, actualMap)
}
