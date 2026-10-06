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

// authUserPartitionAccess mirrors one entry of auth/user's partitionAccess
// sub-list.
type authUserPartitionAccess struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// authUserConfig is the request/response body for one auth/user item.
// go-bigip has no user-management support at all, so this resource talks to
// the endpoint directly via restGet/restPut/restPost (see
// bigip/rest_helpers.go).
type authUserConfig struct {
	Name            string                    `json:"name,omitempty"`
	Description     string                    `json:"description,omitempty"`
	Password        string                    `json:"password,omitempty"`
	Shell           string                    `json:"shell,omitempty"`
	PartitionAccess []authUserPartitionAccess `json:"partitionAccess,omitempty"`
}

func resourceBigipAuthUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipAuthUserCreate,
		ReadContext:   resourceBigipAuthUserRead,
		UpdateContext: resourceBigipAuthUserUpdate,
		DeleteContext: resourceBigipAuthUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures a BIG-IP local user account and its partition/role assignments (auth/user).",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Login name of the user account.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "User-defined description, commonly used for the account's display/full name. If not set, BIG-IP defaults this to the account's name.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for the account. BIG-IP never returns this value on read (only encryptedPassword, which this resource does not track); set once and manage rotation out of band, e.g. via lifecycle.ignore_changes.",
			},
			"shell": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "none",
				Description: "Login shell assigned to the user, e.g. bash, tmsh, or none.",
			},
			"partition_access": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "One or more partition/role assignments for this user.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Partition name the role applies to, or \"all-partitions\" for every partition.",
						},
						"role": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Role granted on partition, e.g. admin, resource-admin, manager, guest, operator, application-editor, auditor, no-access.",
						},
					},
				},
			},
		},
	}
}

func getAuthUserConfig(d *schema.ResourceData) *authUserConfig {
	config := &authUserConfig{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Password:    d.Get("password").(string),
		Shell:       d.Get("shell").(string),
	}
	access := d.Get("partition_access").([]interface{})
	config.PartitionAccess = make([]authUserPartitionAccess, 0, len(access))
	for _, a := range access {
		m := a.(map[string]interface{})
		config.PartitionAccess = append(config.PartitionAccess, authUserPartitionAccess{
			Name: m["partition"].(string),
			Role: m["role"].(string),
		})
	}
	return config
}

func resourceBigipAuthUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Auth User " + name)

	if err := restPost(client, "auth/user", getAuthUserConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to create Auth User (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipAuthUserRead(ctx, d, meta)
}

func resourceBigipAuthUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Auth User " + name)

	// PATCH, not PUT: password is Optional and BIG-IP never returns it on
	// read, so it won't be tracked in state for an imported user. Use
	// PATCH (partial update) rather than PUT so an update that doesn't
	// touch password can't affect the account's existing password,
	// regardless of full-replace-vs-merge semantics on any given BIG-IP
	// version.
	path := "auth/user/" + mangleFullPath(name)
	if err := restPatch(client, path, getAuthUserConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update Auth User (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipAuthUserRead(ctx, d, meta)
}

func resourceBigipAuthUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading Auth User " + name)

	var cfg authUserConfig
	path := "auth/user/" + mangleFullPath(name)
	found, err := restGet(client, path, &cfg)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Auth User (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] Auth User (%s) not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", cfg.Name)
	_ = d.Set("description", cfg.Description)
	_ = d.Set("shell", cfg.Shell)

	access := make([]map[string]interface{}, 0, len(cfg.PartitionAccess))
	for _, a := range cfg.PartitionAccess {
		access = append(access, map[string]interface{}{
			"partition": a.Name,
			"role":      a.Role,
		})
	}
	_ = d.Set("partition_access", access)
	// password is intentionally not set here: BIG-IP never returns it on read.

	return nil
}

func resourceBigipAuthUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Deleting Auth User " + name)

	path := "auth/user/" + mangleFullPath(name)
	if err := restDelete(client, path); err != nil {
		log.Printf("[ERROR] Unable to delete Auth User (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
