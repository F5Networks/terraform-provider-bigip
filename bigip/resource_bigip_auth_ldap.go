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

// authLdapConfig is the request/response body for auth/ldap/~Common~system-auth.
// go-bigip has no LDAP auth support at all, so this resource talks to the
// endpoint directly via restGet/restPut/restPost (see bigip/rest_helpers.go).
type authLdapConfig struct {
	Name                 string   `json:"name,omitempty"`
	Servers              []string `json:"servers,omitempty"`
	Port                 int      `json:"port,omitempty"`
	BindDn               string   `json:"bindDn,omitempty"`
	BindPw               string   `json:"bindPw,omitempty"`
	BindTimeout          int      `json:"bindTimeout,omitempty"`
	SearchBaseDn         string   `json:"searchBaseDn,omitempty"`
	LoginAttribute       string   `json:"loginAttribute,omitempty"`
	UserTemplate         string   `json:"userTemplate,omitempty"`
	Filter               string   `json:"filter,omitempty"`
	GroupDn              string   `json:"groupDn,omitempty"`
	GroupMemberAttribute string   `json:"groupMemberAttribute,omitempty"`
	Ssl                  string   `json:"ssl,omitempty"`
	SslCaCertFile        string   `json:"sslCaCertFile,omitempty"`
	SslCheckPeer         string   `json:"sslCheckPeer,omitempty"`
	Version              int      `json:"version,omitempty"`
}

// authLdapFullPath is the fixed full path of the BIG-IP's single system
// authentication LDAP configuration object.
const authLdapFullPath = "/Common/system-auth"

func resourceBigipAuthLdap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAuthLdapCreate,
		ReadContext:   resourceBigipAuthLdapRead,
		UpdateContext: resourceBigipAuthLdapUpdate,
		DeleteContext: resourceBigipAuthLdapDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures the BIG-IP system's LDAP remote authentication settings (auth/ldap system-auth). This is a singleton resource: BIG-IP supports exactly one system LDAP configuration. Note this resource only configures the LDAP backend itself; use bigip_command or a separate resource to set auth/source (\"type\": \"ldap\") to activate it as the active authentication source.",
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of LDAP server IP addresses or hostnames.",
			},
			"port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     389,
				Description: "Port the BIG-IP system uses to communicate with the LDAP servers.",
			},
			"bind_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Distinguished name the BIG-IP system uses to bind to the LDAP server.",
			},
			"bind_pw": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for bind_dn. BIG-IP never returns this value on read, so it will show as an unresolvable diff source if not also tracked out of band; set once and use lifecycle.ignore_changes if BIG-IP-side rotation is expected.",
			},
			"bind_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     30,
				Description: "Number of seconds the BIG-IP system waits for a response when binding to an LDAP server.",
			},
			"search_base_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Base distinguished name to use as the starting point for user searches.",
			},
			"login_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP attribute that contains the user login name.",
			},
			"user_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template the BIG-IP system uses to construct a user's distinguished name, e.g. \"uid=%s,ou=people,dc=example,dc=com\".",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP search filter used when searching for a user.",
			},
			"group_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Distinguished name of the group used for group-based role assignment.",
			},
			"group_member_attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "LDAP attribute the BIG-IP system uses to determine group membership.",
			},
			"ssl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Whether the BIG-IP system uses SSL/TLS to communicate with the LDAP servers. Valid values: enabled, disabled, start-tls.",
			},
			"ssl_ca_cert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Full path on the BIG-IP system of the CA certificate file used to verify the LDAP server's certificate.",
			},
			"ssl_check_peer": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "disabled",
				Description: "Whether the BIG-IP system verifies the LDAP server's SSL/TLS certificate.",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3,
				Description: "LDAP protocol version.",
			},
		},
	}
}

func getAuthLdapConfig(d *schema.ResourceData) *authLdapConfig {
	return &authLdapConfig{
		Name:                 authLdapFullPath,
		Servers:              listToStringSlice(d.Get("servers").([]interface{})),
		Port:                 d.Get("port").(int),
		BindDn:               d.Get("bind_dn").(string),
		BindPw:               d.Get("bind_pw").(string),
		BindTimeout:          d.Get("bind_timeout").(int),
		SearchBaseDn:         d.Get("search_base_dn").(string),
		LoginAttribute:       d.Get("login_attribute").(string),
		UserTemplate:         d.Get("user_template").(string),
		Filter:               d.Get("filter").(string),
		GroupDn:              d.Get("group_dn").(string),
		GroupMemberAttribute: d.Get("group_member_attribute").(string),
		Ssl:                  d.Get("ssl").(string),
		SslCaCertFile:        d.Get("ssl_ca_cert_file").(string),
		SslCheckPeer:         d.Get("ssl_check_peer").(string),
		Version:              d.Get("version").(int),
	}
}

func resourceBigipAuthLdapCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Creating Auth LDAP config " + authLdapFullPath)

	if err := restPost(client, "auth/ldap", getAuthLdapConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to create Auth LDAP config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId(authLdapFullPath)
	return resourceBigipAuthLdapRead(ctx, d, meta)
}

func resourceBigipAuthLdapUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Updating Auth LDAP config " + d.Id())

	// PATCH, not PUT: verified against a live device that PUT here performs
	// a full replace -- omitting bind_pw (e.g. because it's Optional and
	// not tracked in state, since BIG-IP never returns it on read) silently
	// clears both bind_pw *and* bind_dn from the live config, not just the
	// field left out of the request body.
	path := "auth/ldap/" + mangleFullPath(d.Id())
	if err := restPatch(client, path, getAuthLdapConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update Auth LDAP config (%v)", err)
		return diag.FromErr(err)
	}
	return resourceBigipAuthLdapRead(ctx, d, meta)
}

func resourceBigipAuthLdapRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Reading Auth LDAP config " + d.Id())

	var cfg authLdapConfig
	path := "auth/ldap/" + mangleFullPath(d.Id())
	found, err := restGet(client, path, &cfg)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Auth LDAP config (%v)", err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] Auth LDAP config (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("servers", cfg.Servers)
	_ = d.Set("port", cfg.Port)
	_ = d.Set("bind_dn", cfg.BindDn)
	_ = d.Set("bind_timeout", cfg.BindTimeout)
	_ = d.Set("search_base_dn", cfg.SearchBaseDn)
	_ = d.Set("login_attribute", cfg.LoginAttribute)
	_ = d.Set("user_template", cfg.UserTemplate)
	_ = d.Set("filter", cfg.Filter)
	_ = d.Set("group_dn", cfg.GroupDn)
	_ = d.Set("group_member_attribute", cfg.GroupMemberAttribute)
	_ = d.Set("ssl", cfg.Ssl)
	_ = d.Set("ssl_ca_cert_file", cfg.SslCaCertFile)
	_ = d.Set("ssl_check_peer", cfg.SslCheckPeer)
	_ = d.Set("version", cfg.Version)
	// bind_pw is intentionally not set here: BIG-IP never returns it on
	// read, so leave whatever value is already in state/config alone.

	return nil
}

func resourceBigipAuthLdapDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Deleting Auth LDAP config " + d.Id())

	path := "auth/ldap/" + mangleFullPath(d.Id())
	if err := restDelete(client, path); err != nil {
		log.Printf("[ERROR] Unable to delete Auth LDAP config (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
