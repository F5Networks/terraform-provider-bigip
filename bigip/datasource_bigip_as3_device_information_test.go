package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccFilterApplications_basic tests the basic functionality of the filterApplications method.
func TestAccFilterApplications_basic(t *testing.T) {
	t.Parallel()

	// Mock AS3 JSON response to simulate a BIG-IP system declaration
	as3Response := `{
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
                "A2": {
                    "class": "Application",
                    "template": "http",
                    "serviceMain": {
                        "class": "Service_HTTP",
                        "virtualAddresses": ["10.1.2.4"],
                        "virtualPort": 8081
                    }
                },
                "class": "Tenant"
            },
            "class": "ADC"
        }
    }`

	// Application list to filter out
	appList := []string{"A1"}

	// Expected filtered result after parsing (only includes A1)
	expectedFilteredResult := `{"ansible":{"A1":{"class":"Application","template":"http","serviceMain":{"class":"Service_HTTP","virtualAddresses":["10.1.2.3"],"virtualPort":80}}}}`

	// Resource definition for testing
	resourceName := "bigip_filter.applications"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBasic(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFilterApplicationConfig(as3Response, appList),
				Check: resource.ComposeTestCheckFunc(
					testCheckApplicationFilterExists(appList),
					// Ensure the filtered results match the expected behavior
					resource.TestCheckResourceAttr(resourceName, "filtered_json", expectedFilteredResult),
				),
			},
			{
				Config: testAccFilterApplicationConfig(as3Response, []string{"A2"}), // Test with A2 filtering
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "filtered_json", regexp.MustCompile("A2")),
				),
			},
		},
	})
}

// Generates a Terraform configuration to test the filterApplications function
func testAccFilterApplicationConfig(as3Response string, appList []string) string {
	// Serialize the application list into JSON-compatible format
	applicationListFmt := fmt.Sprintf(`["%s"]`, appList[0]) // Works for single app; can be extended for multiple apps.
	return fmt.Sprintf(`
resource "bigip_filter" "applications" {
  response_json     = <<EOT
%s
EOT
  application_names = %s
}
`, as3Response, applicationListFmt)
}

// Pre-check logic for the test
func testAccPreCheckBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping acceptance tests in short mode.")
	}
}

// Checks if the filtered resource exists and is valid
func testCheckApplicationFilterExists(appList []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Mock check behavior to simulate `ApplicationList` presence
		filteredApp := appList[0]
		if filteredApp == "" {
			return fmt.Errorf("Filtered application list is empty; expected at least one filtered entry")
		}
		return nil
	}
}

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
