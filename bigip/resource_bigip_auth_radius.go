/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// authRadiusConfig is the request/response body for
// auth/radius/~Common~system-auth: the system-wide RADIUS authentication
// configuration, which references one or more auth/radius-server objects by
// name. go-bigip has no RADIUS auth support at all, so this resource talks
// to the endpoint directly via restGet/restPut/restPost (see
// bigip/rest_helpers.go).
type authRadiusConfig struct {
	Name        string   `json:"name,omitempty"`
	ServiceType string   `json:"serviceType,omitempty"`
	Retries     int      `json:"retries,omitempty"`
	Servers     []string `json:"servers,omitempty"`
}

// authRadiusFullPath is the fixed full path of the BIG-IP's single system
// authentication RADIUS configuration object.
const authRadiusFullPath = "/Common/system-auth"

func resourceBigipAuthRadius() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAuthRadiusCreate,
		ReadContext:   resourceBigipAuthRadiusRead,
		UpdateContext: resourceBigipAuthRadiusUpdate,
		DeleteContext: resourceBigipAuthRadiusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures the BIG-IP system's system-wide RADIUS remote authentication settings (auth/radius system-auth), referencing one or more bigip_auth_radius_server resources by name. This is a singleton resource: BIG-IP supports exactly one system RADIUS configuration. Note this resource only configures the RADIUS backend itself; use a separate resource to set auth/source (\"type\": \"radius\") to activate it as the active authentication source.",
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of the auth/radius-server (bigip_auth_radius_server) objects to use, in priority order.",
			},
			"service_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "authenticate-only",
				Description: "RADIUS service type BIG-IP requests, e.g. authenticate-only, login, framed, callback-login, callback-framed, outbound, administrative, nas-prompt, call-check, callback-nas-prompt.",
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Number of times the BIG-IP system retries a RADIUS request before failing over to the next server.",
			},
		},
	}
}

func getAuthRadiusConfig(d *schema.ResourceData) *authRadiusConfig {
	return &authRadiusConfig{
		Name:        authRadiusFullPath,
		Servers:     listToStringSlice(d.Get("servers").([]interface{})),
		ServiceType: d.Get("service_type").(string),
		Retries:     d.Get("retries").(int),
	}
}

func resourceBigipAuthRadiusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Creating Auth RADIUS config " + authRadiusFullPath)

	if err := restPost(client, "auth/radius", getAuthRadiusConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to create Auth RADIUS config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId(authRadiusFullPath)
	return resourceBigipAuthRadiusRead(ctx, d, meta)
}

func resourceBigipAuthRadiusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Updating Auth RADIUS config " + d.Id())

	path := "auth/radius/" + mangleFullPath(d.Id())
	if err := restPut(client, path, getAuthRadiusConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update Auth RADIUS config (%v)", err)
		return diag.FromErr(err)
	}
	return resourceBigipAuthRadiusRead(ctx, d, meta)
}

func resourceBigipAuthRadiusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Reading Auth RADIUS config " + d.Id())

	var cfg authRadiusConfig
	path := "auth/radius/" + mangleFullPath(d.Id())
	found, err := restGet(client, path, &cfg)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Auth RADIUS config (%v)", err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] Auth RADIUS config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	servers := make([]string, len(cfg.Servers))
	for i, s := range cfg.Servers {
		servers[i] = stripPartitionPrefix(s)
	}
	_ = d.Set("servers", servers)
	_ = d.Set("service_type", cfg.ServiceType)
	_ = d.Set("retries", cfg.Retries)

	return nil
}

func resourceBigipAuthRadiusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Deleting Auth RADIUS config " + d.Id())

	path := "auth/radius/" + mangleFullPath(d.Id())
	if err := restDelete(client, path); err != nil {
		log.Printf("[ERROR] Unable to delete Auth RADIUS config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
