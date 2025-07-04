/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	uuid "github.com/google/uuid"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

// var x = 0
var m sync.Mutex
var createdTenants string

func resourceBigipAs3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAs3Create,
		ReadContext:   resourceBigipAs3Read,
		UpdateContext: resourceBigipAs3Update,
		DeleteContext: resourceBigipAs3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// d.Id() here is the last argument passed to the `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_ID` command
				// Here we use a function to parse the import ID (like the example above) to simplify our logic

				_ = d.Set("tenant_list", d.Id())
				_ = d.Set("tenant_filter", d.Id())

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"as3_json": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Full AS3 declaration as a JSON string to deploy on BIG-IP. **Mutually exclusive with `delete_apps`**: only one of `as3_json` or `delete_apps` can be set in a resource block.",
				ConflictsWith: []string{"delete_apps"},
				StateFunc: func(v interface{}) string {
					jsonString, _ := structure.NormalizeJsonString(v)
					return jsonString
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldResp := []byte(old)
					newResp := []byte(new)
					oldJsonref := make(map[string]interface{})
					newJsonref := make(map[string]interface{})
					_ = json.Unmarshal(oldResp, &oldJsonref)
					_ = json.Unmarshal(newResp, &newJsonref)
					delete(oldJsonref, "$schema")
					delete(newJsonref, "$schema")
					jsonEqualityBefore := reflect.DeepEqual(oldJsonref, newJsonref)
					if jsonEqualityBefore {
						return true
					}
					for key, value := range oldJsonref {
						if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
							for range rec {
								delete(rec, "updateMode")
								delete(rec, "schemaVersion")
								delete(rec, "id")
								delete(rec, "label")
								delete(rec, "remark")
								delete(rec, "Common")
							}
						}
						if key == "persist" {
							delete(oldJsonref, "persist")
						}
					}
					for key, value := range newJsonref {
						if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
							for range rec {
								delete(rec, "updateMode")
								delete(rec, "schemaVersion")
								delete(rec, "id")
								delete(rec, "label")
								delete(rec, "remark")
								delete(rec, "Common")
							}
						}
						if key == "persist" {
							delete(newJsonref, "persist")
						}
					}
					ignoreMetadata := d.Get("ignore_metadata").(bool)
					jsonEqualityAfter := reflect.DeepEqual(oldJsonref, newJsonref)
					if ignoreMetadata {
						if jsonEqualityAfter {
							return true
						} else {
							return false
						}

					} else {
						if !jsonEqualityBefore {
							return false
						}
					}
					return true
				},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if _, err := structure.NormalizeJsonString(v); err != nil {
						errors = append(errors, fmt.Errorf("%q contains an invalid JSON: %s", k, err))
					}
					as3json := v.(string)
					resp := []byte(as3json)
					jsonRef := make(map[string]interface{})
					_ = json.Unmarshal(resp, &jsonRef)
					for key, value := range jsonRef {
						if key == "class" && value != "AS3" {
							errors = append(errors, fmt.Errorf("JSON must have AS3 class"))
						}
						if rec, ok := value.(map[string]interface{}); ok && key == "declaration" {
							for k, v := range rec {
								if k == "class" && v != "ADC" {
									errors = append(errors, fmt.Errorf("JSON must have ADC class"))
								}
							}
						}
					}
					return
				},
			},
			"controls": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:         true,
				Description:      "Controls parameters for AS3, you can use the following parameters, dry_run, trace, trace_response, log_level, user_agent",
				ValidateDiagFunc: validateControlsParam,
			},
			"ignore_metadata": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set True if you want to ignore metadata update",
				Default:     false,
			},
			"tenant_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of Tenant. This name is used only in the case of Per-Application Deployment. If it is not provided, then a random name would be generated.",
			},
			"tenant_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of Tenant",
			},
			"tenant_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Tenant",
			},
			"application_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Application deployed through AS3 Declaration",
			},
			"task_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "ID of AS3 post declaration async task",
			},
			"per_app_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Will define Perapp mode enabled on BIG-IP or not",
			},
			"delete_apps": {
				Type:          schema.TypeList,
				MaxItems:      1,    // Ensures only one delete_apps block is allowed
				Optional:      true, // The block is optional in the configuration
				Description:   "Block for specifying tenant name and applications to delete from BIG-IP. **Mutually exclusive with `as3_json`**: only one of `delete_apps` or `as3_json` can be set in a resource block.",
				ConflictsWith: []string{"as3_json"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the tenant for application deletion.",
						},
						"apps": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "List of applications to delete from the specified tenant.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func validateControlsParam(val interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	controls, ok := val.(map[string]interface{})
	if !ok {
		return diag.Diagnostics{{
			Severity: diag.Error,
			Summary:  "Invalid type",
			Detail:   "The controls parameter must be a map.",
		}}
	}

	allowedKeys := map[string]bool{
		"dry_run":        true,
		"trace":          true,
		"trace_response": true,
		"log_level":      true,
		"user_agent":     true,
	}

	for k, v := range controls {
		value := fmt.Sprintf("%v", v)
		if !allowedKeys[k] {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid key",
				Detail:   fmt.Sprintf("The key %s is not allowed in the 'controls' attribute", k),
			})
			continue
		}

		switch k {
		case "dry_run", "trace", "trace_response":
			if value != "yes" && value != "no" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid value",
					Detail:   fmt.Sprintf("The value for key %s must be yes or no", k),
				})
			}
		case "log_level":
			if value != "emergency" &&
				value != "alert" &&
				value != "critical" &&
				value != "error" &&
				value != "warning" &&
				value != "notice" &&
				value != "info" &&
				value != "debug" {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Invalid value",
					Detail:   fmt.Sprintf("The value for key %s must be one of emergency, alert, critical, error, warning, notice, info, debug", k),
				})
			}
		case "user_agent":
			// No specific validation for user_agent, just ensure the key is valid
		}
	}

	return diags
}

func controlsQueraString(d *schema.ResourceData) string {
	controls := d.Get("controls").(map[string]interface{})
	query := ""
	if dryRun, ok := controls["dry_run"]; ok {
		if dryRun.(string) == "yes" {
			query += "&controls.dryRun=true"
		} else if dryRun.(string) == "no" {
			query += "&controls.dryRun=false"
		}
	}
	if trace, ok := controls["trace"]; ok {
		if trace.(string) == "yes" {
			query += "&controls.trace=true"
		} else if trace.(string) == "no" {
			query += "&controls.trace=false"
		}
	}
	if traceResponse, ok := controls["trace_response"]; ok {
		if traceResponse.(string) == "yes" {
			query += "&controls.traceResponse=true"
		} else if traceResponse.(string) == "no" {
			query += "&controls.traceResponse=false"
		}
	}
	if logLevel, ok := controls["log_level"]; ok {
		query += fmt.Sprintf("&controls.logLevel=%s", logLevel.(string))
	}

	return query
}

func resourceBigipAs3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	as3Json := d.Get("as3_json").(string)
	tenantFilter := d.Get("tenant_filter").(string)
	deleteAppsBlocks := d.Get("delete_apps").([]interface{})

	// Check if delete_apps is set, call the delete_apps handler
	// can you please properly fix the below code with proper logging?
	// saying the delete_apps and as3_json is mutual exclusive

	if len(deleteAppsBlocks) > 0 {
		log.Printf("[INFO] Detected delete_apps block. Redirecting to deletion logic.")
		return handleDeleteApps(ctx, d, client)
	}

	var tenantCount []string
	perApplication, err := client.CheckSetting()
	if err != nil {
		return diag.FromErr(err)
	}
	tenantList, _, applicationList := client.GetTenantList(as3Json)

	var controlsQuerParam string
	if _, controls := d.GetOk("controls"); controls {
		controlsQuerParam = controlsQueraString(d)
	}

	log.Printf("[DEBUG] perApplication:%+v", perApplication)

	if perApplication && len(tenantList) == 0 {
		log.Printf("[INFO] Creating As3 config perApplication : tenant name :%+v", d.Get("tenant_name").(string))
		var tenant string
		if d.Get("tenant_name").(string) != "" {
			tenant = d.Get("tenant_name").(string)
		} else {
			tenant, err = GenerateRandomString(10)
			if err != nil {
				return diag.FromErr(fmt.Errorf("could not generate random tenant name"))
			}
		}
		log.Printf("[DEBUG] tenant name :%+v", tenant)

		applicationList := client.GetAppsList(as3Json)
		err, taskID := client.PostPerAppBigIp(as3Json, tenant, controlsQuerParam)
		log.Printf("[DEBUG] task Id from deployment :%+v", taskID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("posting as3 config failed for tenants:(%s) with error: %v", tenantFilter, err))
		}
		tenantCount = append(tenantCount, tenant)
		_ = d.Set("tenant_filter", tenant)
		_ = d.Set("tenant_name", tenant)
		_ = d.Set("tenant_list", tenant)
		_ = d.Set("task_id", taskID)
		_ = d.Set("application_list", applicationList)
		_ = d.Set("per_app_mode", true)
	} else {
		log.Printf("[INFO] Creating AS3 config traditionally for tenants:%+v", tenantList)
		tenantCount := strings.Split(tenantList, ",")
		if tenantFilter != "" {
			log.Printf("[DEBUG] tenantFilter:%+v", tenantFilter)
			if !contains(tenantCount, tenantFilter) {
				return diag.FromErr(fmt.Errorf("tenant_filter: (%s) not exist in as3_json provided ", tenantFilter))
			}
			tenantList = tenantFilter
		}
		_ = d.Set("tenant_list", tenantList)
		_ = d.Set("application_list", applicationList)

		strTrimSpace, err := client.AddTeemAgent(as3Json)
		if err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] Creating as3 config in bigip:%s", strTrimSpace)
		err, successfulTenants, taskID := client.PostAs3Bigip(strTrimSpace, tenantList, controlsQuerParam)
		log.Printf("[DEBUG] successfulTenants :%+v", successfulTenants)
		if err != nil {
			if successfulTenants == "" {
				return diag.FromErr(fmt.Errorf("posting as3 config failed for tenants:(%s) with error: %v", tenantList, err))
			}
			_ = d.Set("tenant_list", successfulTenants)
			if len(successfulTenants) != len(tenantList) {
				return diag.FromErr(err)
			}
		}

		log.Printf("[DEBUG] ID for resource :%+v", d.Get("tenant_list").(string))
		_ = d.Set("task_id", taskID)
		_ = d.Set("per_app_mode", false)
		_ = d.Set("tenant_name", tenantList)
	}

	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Number_of_tenants": len(tenantCount),
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_as3", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	if d.Get("tenant_list").(string) != "" {
		d.SetId(d.Get("tenant_list").(string))
	} else {
		d.SetId("Common")
	}
	createdTenants = d.Get("tenant_list").(string)
	return resourceBigipAs3Read(ctx, d, meta)
}
func resourceBigipAs3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading AS3 config")
	var name string
	var tList string
	as3Json := d.Get("as3_json").(string)
	perappMode := d.Get("per_app_mode").(bool)
	log.Printf("[INFO] AS3 config:%+v", as3Json)
	if d.Get("as3_json") != nil && !perappMode && d.Get("tenant_filter") == "" {
		tList, _, _ = client.GetTenantList(as3Json)
		if createdTenants != "" && createdTenants != tList {
			tList = createdTenants
		}
	}
	if d.Id() != "" && tList != "" {
		name = tList
	} else {
		name = d.Id()
	}
	applicationList := d.Get("application_list").(string)
	log.Printf("[DEBUG] Tenants in AS3 get call : %s", name)
	log.Printf("[DEBUG] Applications in AS3 get call : %s", applicationList)
	if name != "" {
		as3Resp, err := client.GetAs3(name, applicationList, d.Get("per_app_mode").(bool))

		if err != nil {
			log.Printf("[ERROR] Unable to retrieve json ")
			if err.Error() == "unexpected end of JSON input" {
				log.Printf("[ERROR] %v", err)
				return nil
			}
			d.SetId("")
			return diag.FromErr(err)
		}
		if as3Resp == "" {
			log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
			_ = d.Set("as3_json", "")
			// d.SetId("")
			return nil
		}

		if d.Get("per_app_mode").(bool) {
			as3Json := make(map[string]interface{})
			filteredAs3Json := make(map[string]interface{})
			_ = json.Unmarshal([]byte(as3Resp), &as3Json)
			for _, value := range strings.Split(applicationList, ",") {
				log.Printf("[DEBUG] Fetching  AS3 get for Application : %s", value)
				filteredAs3Json[value] = as3Json[value]
			}
			filteredAs3Json["schemaVersion"] = as3Json["schemaVersion"]
			out, _ := json.Marshal(filteredAs3Json)
			filteredAs3String := string(out)
			log.Printf("[DEBUG] AS3 GET call in Read function : %s", filteredAs3Json)
			_ = d.Set("as3_json", filteredAs3String)
		} else {
			_ = d.Set("as3_json", as3Resp)
		}

		_ = d.Set("tenant_list", name)
	} else if d.Get("task_id") != nil {
		taskResponse, err := client.Getas3TaskResponse(d.Get("task_id").(string))
		if err != nil {
			d.SetId("")
			return nil
		}
		_ = d.Set("as3_json", taskResponse)
		_ = d.Set("tenant_list", name)
	}
	return nil
}

func resourceBigipAs3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	as3Json := d.Get("as3_json").(string)
	deleteAppsBlocks := d.Get("delete_apps").([]interface{})

	// Handle specific application deletions if delete_apps is set
	if len(deleteAppsBlocks) > 0 {
		log.Printf("[INFO] Detected delete_apps block. Redirecting to deletion-specific logic.")
		return handleDeleteApps(ctx, d, client)
	}
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	oldApplicationList := d.Get("application_list").(string)
	tenantList, _, applicationList := client.GetTenantList(as3Json)

	var controlsQuerParam string
	if _, controls := d.GetOk("controls"); controls {
		controlsQuerParam = controlsQueraString(d)
	}

	_ = d.Set("application_list", applicationList)
	perApplication, err := client.CheckSetting()
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] perApplication:%+v", perApplication)
	if d.Get("per_app_mode").(bool) {
		if perApplication && len(tenantList) == 0 {
			oldTenantList := d.Id()
			log.Printf("[INFO] oldApplicationList :%s", oldApplicationList)
			curApplicationList := client.GetAppsList(as3Json)
			log.Printf("[INFO] curApplicationList :%s", curApplicationList)
			for _, appName := range strings.Split(oldApplicationList, ",") {
				if !strings.Contains(curApplicationList, appName) {
					log.Printf("[INFO] Deleting As3 Config for Application:%s in Tenant:%v", appName, oldTenantList)
					err := client.DeletePerApplicationAs3Bigip(oldTenantList, appName)
					if err != nil {
						log.Printf("[ERROR] Unable to DeleteContext: %v :", err)
						return diag.FromErr(err)
					}
				}
			}

			log.Printf("[INFO] Updating As3 Config for tenant:%s with Per-Application Mode:%v", oldTenantList, perApplication)
			err, task_id := client.PostPerAppBigIp(as3Json, oldTenantList, controlsQuerParam)
			log.Printf("[DEBUG] task_id from PostPerAppBigIp:%+v", task_id)
			if err != nil {
				return diag.FromErr(fmt.Errorf("posting as3 config failed for tenant:(%s) with error: %v", oldTenantList, err))
			}
			// tenantCount = append(tenantCount, tenant)
			_ = d.Set("tenant_list", oldTenantList)
			_ = d.Set("task_id", task_id)
			_ = d.Set("tenant_filter", oldTenantList)
			_ = d.Set("application_list", curApplicationList)

		} else {
			if !perApplication {
				return diag.FromErr(fmt.Errorf("Per-Application should be true in Big-IP Setting"))
			} else {
				return diag.FromErr(fmt.Errorf("declartion not valid for Per-Application deployment"))
			}
		}
	} else {
		log.Printf("[INFO] Updating As3 Config Traditionally for tenants:%s", tenantList)
		oldTenantList := d.Get("tenant_list").(string)
		tenantFilter := d.Get("tenant_filter").(string)
		if tenantFilter == "" {
			if tenantList != oldTenantList {
				_ = d.Set("tenant_list", tenantList)
				newList := strings.Split(tenantList, ",")
				oldList := strings.Split(oldTenantList, ",")
				deletedTenants := client.TenantDifference(oldList, newList)
				if deletedTenants != "" {
					err, _ := client.DeleteAs3Bigip(deletedTenants)
					if err != nil {
						log.Printf("[ERROR] Unable to Delete removed tenants: %v :", err)
						return diag.FromErr(err)
					}
				}
			}
		} else {
			if !contains(strings.Split(tenantList, ","), tenantFilter) {
				log.Printf("[WARNING]tenant_filter: (%s) not exist in as3_json provided ", tenantFilter)
			} else {
				tenantList = tenantFilter
			}
		}
		strTrimSpace, err := client.AddTeemAgent(as3Json)
		if err != nil {
			return diag.FromErr(err)
		}
		err, successfulTenants, taskID := client.PostAs3Bigip(strTrimSpace, tenantList, controlsQuerParam)
		log.Printf("[DEBUG] successfulTenants :%+v", successfulTenants)
		if err != nil {
			if successfulTenants == "" {
				return diag.FromErr(fmt.Errorf("error updating json  %s: %v", tenantList, err))
			}
			_ = d.Set("tenant_list", successfulTenants)
			if len(successfulTenants) != len(tenantList) {
				return diag.FromErr(err)
			}
		}
		_ = d.Set("task_id", taskID)
		_ = d.Set("tenant_name", tenantList)
	}
	if d.Get("tenant_filter").(string) != "" {
		createdTenants = d.Get("tenant_filter").(string)
	} else {
		createdTenants = d.Get("tenant_list").(string)
	}
	return resourceBigipAs3Read(ctx, d, meta)
}

func resourceBigipAs3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	var name string
	var tList string
	as3Json := d.Get("as3_json").(string)
	deleteAppsBlocks := d.Get("delete_apps").([]interface{})

	// Handle specific application deletions if delete_apps is set
	if len(deleteAppsBlocks) > 0 {
		log.Printf("[INFO] Detected delete_apps block. Redirecting to deletion-specific logic.")
		return handleDeleteApps(ctx, d, client)
	}

	if c_attr, c_ok := d.GetOk("controls"); c_ok {
		controls := c_attr.(map[string]interface{})
		if dryRun, ok := controls["dry_run"]; ok && dryRun == "yes" {
			d.SetId("")
			return nil
		}
	}

	if d.Get("as3_json") != nil {
		tList, _, _ = client.GetTenantList(as3Json)
	}

	if d.Id() != "" && tList != "" && d.Get("tenant_filter") == "" {
		name = tList
	} else {
		name = d.Id()
	}
	log.Printf("[INFO] Deleting As3 config for tenants:%+v", name)
	if d.Get("per_app_mode").(bool) {
		applicationList := d.Get("application_list").(string)
		log.Printf("[INFO] Deleting As3 config for Applications:%+v", applicationList)
		for _, appName := range strings.Split(applicationList, ",") {
			log.Printf("[INFO] Deleting AS3 for Application : %s", appName)
			err := client.DeletePerApplicationAs3Bigip(name, appName)
			if err != nil {
				log.Printf("[ERROR] Unable to DeleteContext: %v :", err)
				return diag.FromErr(err)
			}
		}
	} else {
		err, failedTenants := client.DeleteAs3Bigip(name)
		if err != nil {
			log.Printf("[ERROR] Unable to DeleteContext: %v :", err)
			return diag.FromErr(err)
		}
		if failedTenants != "" {
			_ = d.Set("tenant_list", name)
			return resourceBigipAs3Read(ctx, d, meta)
		}
	}
	d.SetId("")
	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func GenerateRandomString(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	charLen := 0
	randomString := make([]byte, length)
	for i := range randomString {
		if i == 0 {
			charLen = 52
		} else {
			charLen = len(charset)
		}
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(charLen)))
		if err != nil {
			return "", err
		}
		randomString[i] = charset[randomIndex.Int64()]
	}
	return string(randomString), nil
}

func handleDeleteApps(ctx context.Context, d *schema.ResourceData, client *bigip.BigIP) diag.Diagnostics {
	deleteAppsBlocks := d.Get("delete_apps").([]interface{})
	tenantName := d.Get("tenant_name").(string)

	for _, block := range deleteAppsBlocks {
		blockData := block.(map[string]interface{})
		tenant := blockData["tenant_name"].(string)
		appsToDeleteRaw := blockData["apps"].([]interface{})
		var appsToDelete []string
		for _, app := range appsToDeleteRaw {
			appsToDelete = append(appsToDelete, app.(string))
		}

		log.Printf("[INFO] Deleting applications %v under tenant '%s'", appsToDelete, tenant)

		// Check if tenant exists
		as3Resp, err := client.GetAs3(tenant, "", false)
		if err != nil || len(as3Resp) == 0 {
			log.Printf("[WARN] Skipping deletion: Tenant '%s' not found or empty: %v", tenant, err)
			continue // Do not fail – just skip this block
		}

		for _, app := range appsToDelete {
			log.Printf("[INFO] Attempting to delete application '%s' in tenant '%s'", app, tenant)

			err := client.DeletePerApplicationAs3Bigip(tenant, app)
			if err != nil {
				log.Printf("[ERROR] Failed to delete application '%s' in tenant '%s': %v", app, tenant, err)
				return diag.FromErr(fmt.Errorf("failed to delete app '%s': %v", app, err))
			}
			log.Printf("[INFO] Successfully deleted application '%s' in tenant '%s'", app, tenant)
		}
	}
	d.SetId(fmt.Sprintf("deleted-%s-%d", tenantName, time.Now().Unix()))

	return nil
}
