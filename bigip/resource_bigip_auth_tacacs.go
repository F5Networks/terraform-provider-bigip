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

// authTacacsConfig is the request/response body for
// auth/tacacs/~Common~system-auth. go-bigip has no TACACS+ auth support at
// all, so this resource talks to the endpoint directly via
// restGet/restPut/restPost (see bigip/rest_helpers.go).
type authTacacsConfig struct {
	Name              string   `json:"name,omitempty"`
	Servers           []string `json:"servers,omitempty"`
	Secret            string   `json:"secret,omitempty"`
	Service           string   `json:"service,omitempty"`
	Protocol          string   `json:"protocol,omitempty"`
	Encryption        string   `json:"encryption,omitempty"`
	Accounting        string   `json:"accounting,omitempty"`
	AuthenticationVia string   `json:"authenticationVia,omitempty"`
}

// authTacacsFullPath is the fixed full path of the BIG-IP's single system
// authentication TACACS+ configuration object.
const authTacacsFullPath = "/Common/system-auth"

func resourceBigipAuthTacacs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAuthTacacsCreate,
		ReadContext:   resourceBigipAuthTacacsRead,
		UpdateContext: resourceBigipAuthTacacsUpdate,
		DeleteContext: resourceBigipAuthTacacsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures the BIG-IP system's TACACS+ remote authentication settings (auth/tacacs system-auth). This is a singleton resource: BIG-IP supports exactly one system TACACS+ configuration. Note this resource only configures the TACACS+ backend itself; use a separate resource to set auth/source (\"type\": \"tacacs\") to activate it as the active authentication source.",
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of TACACS+ server IP addresses or hostnames, in priority order.",
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Shared secret used to authenticate to the TACACS+ servers. BIG-IP never returns this value on read.",
			},
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ppp",
				Description: "TACACS+ service requested of the daemon, e.g. ppp, shell, slip, system.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ip",
				Description: "TACACS+ protocol name associated with the value of service, e.g. ip, lcp, vpdn.",
			},
			"encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "Whether the BIG-IP system encrypts TACACS+ traffic.",
			},
			"accounting": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "send-to-first-server",
				Description: "Specifies whether the system sends TACACS+ accounting messages to the first available server or to all available servers. Valid values: send-to-first-server, send-to-all-servers.",
			},
			"authentication_via": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether the system authenticates via PAP or CHAP.",
			},
		},
	}
}

func getAuthTacacsConfig(d *schema.ResourceData) *authTacacsConfig {
	return &authTacacsConfig{
		Name:              authTacacsFullPath,
		Servers:           listToStringSlice(d.Get("servers").([]interface{})),
		Secret:            d.Get("secret").(string),
		Service:           d.Get("service").(string),
		Protocol:          d.Get("protocol").(string),
		Encryption:        d.Get("encryption").(string),
		Accounting:        d.Get("accounting").(string),
		AuthenticationVia: d.Get("authentication_via").(string),
	}
}

func resourceBigipAuthTacacsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Creating Auth TACACS+ config " + authTacacsFullPath)

	if err := restPost(client, "auth/tacacs", getAuthTacacsConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to create Auth TACACS+ config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId(authTacacsFullPath)
	return resourceBigipAuthTacacsRead(ctx, d, meta)
}

func resourceBigipAuthTacacsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Updating Auth TACACS+ config " + d.Id())

	// PATCH, not PUT: verified against a live device that PUT here rejects
	// the request outright ("\"secret\" is a required property and may not
	// be set to \"none\" or an empty value") whenever secret is omitted --
	// which it will be if it's Optional and not tracked in state, since
	// BIG-IP never returns it on read. That would make every update fail,
	// not just a credential-clearing one, for any tacacs config extracted
	// via import.
	path := "auth/tacacs/" + mangleFullPath(d.Id())
	if err := restPatch(client, path, getAuthTacacsConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update Auth TACACS+ config (%v)", err)
		return diag.FromErr(err)
	}
	return resourceBigipAuthTacacsRead(ctx, d, meta)
}

func resourceBigipAuthTacacsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Reading Auth TACACS+ config " + d.Id())

	var cfg authTacacsConfig
	path := "auth/tacacs/" + mangleFullPath(d.Id())
	found, err := restGet(client, path, &cfg)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Auth TACACS+ config (%v)", err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] Auth TACACS+ config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("servers", cfg.Servers)
	_ = d.Set("service", cfg.Service)
	_ = d.Set("protocol", cfg.Protocol)
	_ = d.Set("encryption", cfg.Encryption)
	_ = d.Set("accounting", cfg.Accounting)
	_ = d.Set("authentication_via", cfg.AuthenticationVia)
	// secret is intentionally not set here: BIG-IP never returns it on read.

	return nil
}

func resourceBigipAuthTacacsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Deleting Auth TACACS+ config " + d.Id())

	path := "auth/tacacs/" + mangleFullPath(d.Id())
	if err := restDelete(client, path); err != nil {
		log.Printf("[ERROR] Unable to delete Auth TACACS+ config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
