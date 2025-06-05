package bigip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBigipAs3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBigipAs3Read,
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The specific AS3 tenant to retrieve configuration for.",
			},
			"applications": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of applications to retrieve from the tenant. Leave empty to fetch all applications for the tenant.",
			},
			"as3_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "JSON string representation of the retrieved AS3 declaration. This is useful for debugging and downstream processing.",
			},
		},
	}
}

func dataSourceBigipAs3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	// Resetting ID to ensure proper handling during the read operation
	d.SetId("")

	// Retrieve tenant name and applications list from the schema
	tenant := d.Get("tenant").(string)
	log.Printf("[INFO] Reading AS3 configuration for tenant: %s", tenant)

	// Extract and build the application list
	applicationList := extractApplications(d)
	log.Printf("[INFO] Application list: %v", applicationList)

	// Validate and normalize the tenant name
	normalizedTenant, err := validateAndNormalizeTenant(tenant)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Normalized tenant name: %s", normalizedTenant)

	// API call to fetch AS3 configuration for the tenant
	as3Resp, err := fetchAs3Configuration(client, normalizedTenant, applicationList)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Successfully fetched AS3 configuration for tenant: %s", normalizedTenant)

	// Filter the JSON response if application filtering is needed
	filteredJSON, err := filterAs3JSON(as3Resp, applicationList)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("as3_json", filteredJSON)

	// Set the tenant as the resource ID
	d.SetId(normalizedTenant)
	log.Printf("[INFO] Successfully processed AS3 configuration for tenant: %s", normalizedTenant)

	return nil
}

func extractApplications(d *schema.ResourceData) []string {
	applicationsRaw := d.Get("applications").([]interface{})
	var applicationList []string
	if len(applicationsRaw) > 0 {
		for _, app := range applicationsRaw {
			applicationList = append(applicationList, app.(string))
		}
	} else {
		log.Printf("[INFO] No specific applications provided, retrieving all applications")
	}
	return applicationList
}

func validateAndNormalizeTenant(tenant string) (string, error) {
	if strings.TrimSpace(tenant) == "" {
		return "", fmt.Errorf("tenant name cannot be empty")
	}

	// Validate tenant name format
	re := regexp.MustCompile(`^(/([a-zA-Z0-9? ,_-]+)/)?([a-zA-Z0-9? ,._-]+)$`)
	if !re.MatchString(tenant) {
		log.Printf("[ERROR] Tenant name '%s' is invalid", tenant)
		return "", fmt.Errorf("tenant name '%s' is invalid. Expected '/partition/tenant_name' or 'tenant_name'.", tenant)
	}

	// Normalize tenant name for API
	if strings.HasPrefix(tenant, "/") {
		tenant = strings.TrimPrefix(tenant, "/")
		tenant = strings.ReplaceAll(tenant, "/", "~") // Convert /partition/tenant to partition~tenant
		log.Printf("[DEBUG] Tenant name normalized to: %s", tenant)
	}
	return tenant, nil
}

func fetchAs3Configuration(client *bigip.BigIP, tenant string, applicationList []string) (string, error) {
	applications := strings.Join(applicationList, ",") // Join application list into CSV for filtering
	log.Printf("[INFO] Fetching AS3 configuration for tenant: %s with applications: %s", tenant, applications)

	as3Resp, err := client.GetAs3(tenant, applications, false)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			log.Printf("[WARN] Tenant '%s' not found on BIG-IP system", tenant)
			return "", fmt.Errorf("tenant '%s' not found on BIG-IP system", tenant)
		} else if strings.Contains(err.Error(), "401") {
			log.Printf("[ERROR] Unauthorized access for tenant '%s'. Check credentials.", tenant)
			return "", fmt.Errorf("unauthorized access when fetching AS3 configuration for tenant '%s'", tenant)
		} else {
			log.Printf("[ERROR] Failed to fetch AS3 configuration for tenant '%s': %v", tenant, err)
			return "", fmt.Errorf("failed to fetch AS3 configuration for tenant '%s': %v", tenant, err)
		}
	}

	if strings.TrimSpace(as3Resp) == "" {
		log.Printf("[WARN] No AS3 configuration found for tenant '%s'", tenant)
		return "", fmt.Errorf("no AS3 configuration found for tenant '%s'", tenant)
	}

	return as3Resp, nil
}

func filterAs3JSON(as3Resp string, applicationList []string) (string, error) {
	if len(applicationList) == 0 {
		log.Printf("[INFO] Returning full AS3 JSON as no application filtering is specified")
		return as3Resp, nil
	}

	// Print applicationList and AS3 response for debugging purposes
	log.Printf("[INFO] Application list: %v", applicationList)
	log.Printf("[INFO] AS3 JSON response: %s", as3Resp)

	// Unmarshal the AS3 JSON response into a map for filtering
	as3Json := make(map[string]interface{})
	if err := json.Unmarshal([]byte(as3Resp), &as3Json); err != nil {
		log.Printf("[ERROR] Failed to parse AS3 JSON response: %v", err)
		return "", fmt.Errorf("failed to parse AS3 JSON response: %w", err)
	}

	// Ensure the response contains the expected "declaration" and tenant-level details
	declaration, ok := as3Json["declaration"].(map[string]interface{})
	if !ok {
		log.Printf("[ERROR] Missing 'declaration' in AS3 response")
		return "", fmt.Errorf("missing 'declaration' in AS3 response")
	}

	// Filter based on the application list
	filteredAs3Json := make(map[string]interface{})

	// Traverse each tenant (like "ansible") in the declaration
	for tenantName, tenant := range declaration {
		tenantMap, ok := tenant.(map[string]interface{})
		if !ok {
			log.Printf("[WARN] Tenant '%s' is not a valid object, skipping", tenantName)
			continue
		}

		// Look for applications inside the tenant and filter them
		for _, appName := range applicationList {
			if app, exists := tenantMap[appName]; exists {
				// Assign filtered application to the filtered JSON
				if _, ok := filteredAs3Json[tenantName]; !ok {
					filteredAs3Json[tenantName] = make(map[string]interface{})
				}
				filteredAs3Json[tenantName].(map[string]interface{})[appName] = app
			} else {
				log.Printf("[WARN] Application '%s' not found in tenant '%s'", appName, tenantName)
			}
		}
	}

	// Marshal filtered JSON back to a string
	filteredJsonBytes, err := json.Marshal(filteredAs3Json)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal filtered AS3 JSON: %v", err)
		return "", fmt.Errorf("failed to process AS3 configuration: %w", err)
	}

	// Log and return the filtered JSON string
	log.Printf("[INFO] Filtered AS3 JSON: %s", string(filteredJsonBytes))
	return string(filteredJsonBytes), nil
}
