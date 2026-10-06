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

// authRadiusServerConfig is the request/response body for one
// auth/radius-server item. go-bigip has no RADIUS auth support at all, so
// this resource talks to the endpoint directly via restGet/restPut/restPost
// (see bigip/rest_helpers.go).
type authRadiusServerConfig struct {
	Name    string `json:"name,omitempty"`
	Server  string `json:"server,omitempty"`
	Secret  string `json:"secret,omitempty"`
	Port    int    `json:"port,omitempty"`
	Timeout int    `json:"timeout,omitempty"`
}

func resourceBigipAuthRadiusServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAuthRadiusServerCreate,
		ReadContext:   resourceBigipAuthRadiusServerRead,
		UpdateContext: resourceBigipAuthRadiusServerUpdate,
		DeleteContext: resourceBigipAuthRadiusServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures a single RADIUS server (auth/radius-server) BIG-IP can use for remote authentication. Reference one or more of these by name from a bigip_auth_radius resource's servers list.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the RADIUS server configuration object.",
			},
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "IP address or hostname of the RADIUS server.",
			},
			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Shared secret used to authenticate to the RADIUS server. BIG-IP never returns this value on read.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1812,
				Description: "Port used to communicate with the RADIUS server.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "Number of seconds the BIG-IP system waits for a response from the RADIUS server.",
			},
		},
	}
}

func getAuthRadiusServerConfig(d *schema.ResourceData) *authRadiusServerConfig {
	return &authRadiusServerConfig{
		Name:    d.Get("name").(string),
		Server:  d.Get("server").(string),
		Secret:  d.Get("secret").(string),
		Port:    d.Get("port").(int),
		Timeout: d.Get("timeout").(int),
	}
}

func resourceBigipAuthRadiusServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Auth RADIUS server " + name)

	if err := restPost(client, "auth/radius-server", getAuthRadiusServerConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to create Auth RADIUS server (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipAuthRadiusServerRead(ctx, d, meta)
}

func resourceBigipAuthRadiusServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Auth RADIUS server " + name)

	// PATCH, not PUT: verified against a live device that PUT on the
	// sibling auth/tacacs endpoint rejects the request outright when its
	// required secret-like field is omitted; auth/radius-server's secret
	// is Required in this schema so d.Get should never actually be empty
	// in practice, but PATCH is used here anyway for consistency with the
	// other auth/* resources and to avoid relying on full-replace-vs-merge
	// semantics that differ across these endpoints.
	path := "auth/radius-server/" + mangleFullPath(name)
	if err := restPatch(client, path, getAuthRadiusServerConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update Auth RADIUS server (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipAuthRadiusServerRead(ctx, d, meta)
}

func resourceBigipAuthRadiusServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Auth RADIUS server " + name)

	var cfg authRadiusServerConfig
	path := "auth/radius-server/" + mangleFullPath(name)
	found, err := restGet(client, path, &cfg)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Auth RADIUS server (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] Auth RADIUS server (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", cfg.Name)
	_ = d.Set("server", cfg.Server)
	_ = d.Set("port", cfg.Port)
	_ = d.Set("timeout", cfg.Timeout)
	// secret is intentionally not set here: BIG-IP never returns it on read.

	return nil
}

func resourceBigipAuthRadiusServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Auth RADIUS server " + name)

	path := "auth/radius-server/" + mangleFullPath(name)
	if err := restDelete(client, path); err != nil {
		log.Printf("[ERROR] Unable to delete Auth RADIUS server (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
