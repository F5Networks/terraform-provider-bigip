/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipSysProvision() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysProvisionCreate,
		UpdateContext: resourceBigipSysProvisionUpdate,
		ReadContext:   resourceBigipSysProvisionRead,
		DeleteContext: resourceBigipSysProvisionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of module to provision in BIG-IP.",
				ValidateFunc: validation.StringInSlice([]string{"afm", "am", "apm", "asm", "avr", "cgnat", "dos", "fps", "gtm", "ilx", "lc", "ltm", "pem", "sslo", "swg", "urldb"}, false),
			},
			"full_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				//ValidateFunc: validation.StringInSlice([]string{"afm", "am", "apm","asm","avr","dos","fps","gtm","ilx","lc","ltm","pem", "sslo" ,"swg","urldb"}, false),
			},
			"cpu_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
			"disk_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
			"level": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Sets the provisioning level for the requested modules. Changing the level for one module may require modifying the level of another module. For example, changing one module to dedicated requires setting all others to none. Setting the level of a module to none means the module is not activated.",
				Default:      "nominal",
				ValidateFunc: validation.StringInSlice([]string{"nominal", "none", "minimum", "dedicated"}, false),
			},
			"memory_ratio": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Use this option only when the level option is set to custom.F5 Networks recommends that you do not modify this option. The default value is none",
			},
		},
	}
}

func resourceBigipSysProvisionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Printf("[INFO] Provisioning for %v module", name)

	pss := &bigip.Provision{
		Name: name,
	}
	config := getsysProvisionConfig(d, pss)

	err := client.ProvisionModule(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Provision  (%s) ", err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	if provisionModuleHitsDevice(name) {
		waitForProvisionReady(client, name)
	}
	return resourceBigipSysProvisionRead(ctx, d, meta)
}

func resourceBigipSysProvisionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Provisioning for :%v module", name)

	pss := &bigip.Provision{
		Name: name,
	}
	config := getsysProvisionConfig(d, pss)

	err := client.ProvisionModule(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Update Provision (%v) ", err)
		return diag.FromErr(err)
	}
	if provisionModuleHitsDevice(name) {
		waitForProvisionReady(client, name)
	}
	return resourceBigipSysProvisionRead(ctx, d, meta)
}

// provisionModuleHitsDevice reports whether ProvisionModule actually
// issues an HTTP request to the device for the given module name --
// go-bigip's ProvisionModule only does so for a fixed set of module
// names (asm, afm, gtm, apm, avr, ilx); any other name is a complete
// no-op. waitForProvisionReady is only worth calling for the former: a
// no-op provisioning call never triggers a device-side restart, so there
// is nothing to wait out.
func provisionModuleHitsDevice(name string) bool {
	switch name {
	case "asm", "afm", "gtm", "apm", "avr", "ilx":
		return true
	default:
		return false
	}
}

// waitForProvisionReady polls the device after a provisioning change
// (ProvisionModule) until it's responding to API calls normally again, or
// a bounded timeout elapses. Provisioning afm/asm in particular restarts
// mcpd and briefly the mgmt layer itself (confirmed directly: the device
// can return a 503 "Configuration Utility restarting..." HTML page, or
// even a harder "Public URI path not registered" error if restjavad's own
// routing table hasn't finished reloading yet, well past go-bigip's own
// single-APICall retry budget) -- without this, a Create/Update's
// immediately-following Read has a real chance of hitting that window and
// failing the apply outright, even though the provisioning change itself
// succeeded. This is deliberately tolerant of (and ignores) errors while
// polling: any error here just means "not ready yet", and the subsequent
// real Read call is what actually surfaces a genuine, non-transient
// failure to the user if the timeout is reached.
//
// Polls both client.Provisions(name) (an already-authenticated API call
// on the existing session/token) AND a brand-new, independent
// /mgmt/shared/authn/login attempt. These were confirmed directly against
// a live device to not always recover at the same time after a
// provisioning-triggered restart: an already-authenticated session's own
// calls can start succeeding again before a *fresh* login does. Acceptance
// tests in particular are affected by this gap specifically because the
// SDK's test harness shells out to a brand-new `terraform plan`/`destroy`
// CLI invocation for its own post-apply/post-test verification steps, each
// establishing a fresh provider configuration (and, since this provider's
// token_auth schema attribute defaults to true when BIGIP_TOKEN_AUTH is
// unset -- as it is in this project's own CI -- a fresh
// /mgmt/shared/authn/login token request) rather than reusing the
// already-authenticated client instance that ran Create/Update/Read
// itself, so only polling the latter was previously missing exactly the
// failure this surfaced: "URI path /mgmt/shared/authn/login not
// registered" from the SDK's own separate verification pass, even though
// waitForProvisionReady itself had already returned successfully.
func waitForProvisionReady(client *bigip.BigIP, name string) {
	const (
		pollInterval            = 5 * time.Second
		pollTimeout             = 10 * time.Minute
		requiredStableSuccesses = 3
	)
	deadline := time.Now().Add(pollTimeout)
	consecutiveSuccesses := 0
	for time.Now().Before(deadline) {
		if _, err := client.Provisions(name); err == nil && freshLoginSucceeds(client) {
			consecutiveSuccesses++
			if consecutiveSuccesses >= requiredStableSuccesses {
				return
			}
		} else {
			consecutiveSuccesses = 0
		}
		time.Sleep(pollInterval)
	}
	log.Printf("[WARN] timed out after %s waiting for the device to settle after provisioning %s; proceeding anyway", pollTimeout, name)
}

// freshLoginSucceeds attempts a brand-new /mgmt/shared/authn/login POST
// (bypassing the given client's own, potentially already-valid session/
// token) and reports whether it succeeds. Used by waitForProvisionReady
// to detect the specific case where an existing authenticated session's
// calls are already succeeding again post-restart, but a new client
// establishing its own fresh session (as every separate Terraform CLI
// invocation in an acceptance test does) still cannot.
func freshLoginSucceeds(client *bigip.BigIP) bool {
	body, err := json.Marshal(map[string]string{
		"username":          client.User,
		"password":          client.Password,
		"loginProviderName": "tmos",
	})
	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPost, client.Host+"/mgmt/shared/authn/login", bytes.NewReader(body))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Transport: client.Transport, Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func resourceBigipSysProvisionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Provisions " + name)
	p, err := client.Provisions(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Provision (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if p == nil {
		log.Printf("[WARN] Provision (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("full_path", p.FullPath)
	_ = d.Set("cpu_ratio", p.CpuRatio)
	_ = d.Set("disk_ratio", p.DiskRatio)
	_ = d.Set("level", p.Level)
	_ = d.Set("memory_ratio", p.MemoryRatio)

	return nil
}

func resourceBigipSysProvisionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// API is not supported for Deleting
	return nil
}

func getsysProvisionConfig(d *schema.ResourceData, config *bigip.Provision) *bigip.Provision {
	config.FullPath = d.Get("full_path").(string)
	config.CpuRatio = d.Get("cpu_ratio").(int)
	config.DiskRatio = d.Get("disk_ratio").(int)
	config.Level = d.Get("level").(string)
	config.MemoryRatio = d.Get("memory_ratio").(int)
	return config
}
