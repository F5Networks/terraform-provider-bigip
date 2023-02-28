/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var p = 0
var q sync.Mutex

func resourceBigiqAs3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigiqAs3Create,
		ReadContext:   resourceBigiqAs3Read,
		UpdateContext: resourceBigiqAs3Update,
		DeleteContext: resourceBigiqAs3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bigiq_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The registration key pool to use",
			},
			"bigiq_user": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The registration key pool to use",
			},
			"bigiq_token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Sensitive:   true,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_TOKEN_AUTH", true),
			},
			"bigiq_login_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Login reference for token authentication (see BIG-IQ REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIQ_LOGIN_REF", "local"),
			},
			"as3_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "AS3 json",
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldResp := []byte(old)
					newResp := []byte(new)
					oldJsonref := make(map[string]interface{})
					newJsonref := make(map[string]interface{})
					_ = json.Unmarshal(oldResp, &oldJsonref)
					_ = json.Unmarshal(newResp, &newJsonref)
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
							}
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
							}
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
				ValidateFunc: validation.StringIsJSON,
			},
			"ignore_metadata": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set True if you want to ignore metadata update",
				Default:     true,
			},
			"tenant_list": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Name of Tenant",
			},
		},
	}
}

func resourceBigiqAs3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return diag.FromErr(err)
	}
	q.Lock()
	defer q.Unlock()
	as3Json := d.Get("as3_json").(string)
	tenantList, _, _ := bigiqRef.GetTenantList(as3Json)
	targetInfo := bigiqRef.GetTarget(as3Json)
	_ = d.Set("tenant_list", tenantList)
	err, successfulTenants := bigiqRef.PostAs3Bigiq(as3Json)
	if err != nil {
		if successfulTenants == "" {
			return diag.FromErr(fmt.Errorf("error creating json  %s: %v", tenantList, err))
		}
		_ = d.Set("tenant_list", successfulTenants)
	}
	as3ID := fmt.Sprintf("%s_%s", targetInfo, successfulTenants)
	d.SetId(as3ID)
	p++
	log.Printf("[TRACE] %+v", p)
	return resourceBigiqAs3Read(ctx, d, meta)
}

func resourceBigiqAs3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return diag.FromErr(err)
	}
	tenantRef := d.Id()
	log.Println("[INFO] Reading As3 config")
	targetRef := strings.Split(tenantRef, "_")[0]
	name := strings.Split(tenantRef, "_")[1]
	if name != d.Get("tenant_list").(string) {
		as3Resp, err := bigiqRef.GetAs3Bigiq(targetRef, d.Get("tenant_list").(string))
		if err != nil {
			log.Printf("[ERROR] Unable to retrieve json ")
			return diag.FromErr(err)
		}
		if as3Resp == "" {
			log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		_ = d.Set("as3_json", as3Resp)
	} else {
		as3Resp, err := bigiqRef.GetAs3Bigiq(targetRef, name)
		if err != nil {
			log.Printf("[ERROR] Unable to retrieve json ")
			return diag.FromErr(err)
		}
		if as3Resp == "" {
			log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		_ = d.Set("as3_json", as3Resp)

	}
	return nil
}

func resourceBigiqAs3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return diag.FromErr(err)
	}
	as3Json := d.Get("as3_json").(string)
	q.Lock()
	defer q.Unlock()
	log.Printf("[INFO] Updating As3 Config :%s", as3Json)
	name := d.Get("tenant_list").(string)
	tenantList, _, _ := bigiqRef.GetTenantList(as3Json)
	if tenantList != name {
		_ = d.Set("tenant_list", tenantList)
	}
	err, successfulTenants := bigiqRef.PostAs3Bigiq(as3Json)
	if err != nil {
		if successfulTenants == "" {
			return diag.FromErr(fmt.Errorf("Error creating json  %s: %v", tenantList, err))
		}
		_ = d.Set("tenant_list", successfulTenants)
	}
	p++
	return resourceBigiqAs3Read(ctx, d, meta)
}

func resourceBigiqAs3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	time.Sleep(20 * time.Second)
	bigiqRef, err := connectBigIq(d)
	if err != nil {
		log.Printf("Connection to BIGIQ Failed with :%v", err)
		return diag.FromErr(err)
	}
	q.Lock()
	defer q.Unlock()
	log.Printf("[INFO] Deleting As3 config")
	name := d.Get("tenant_list").(string)
	as3Json := d.Get("as3_json").(string)
	err, failedTenants := bigiqRef.DeleteAs3Bigiq(as3Json, name)
	if err != nil {
		log.Printf("[ERROR] Unable to DeleteContext: %v :", err)
		return diag.FromErr(err)
	}
	if failedTenants != "" {
		_ = d.Set("tenant_list", name)
		return resourceBigipAs3Read(ctx, d, meta)
	}
	p++
	d.SetId("")
	return nil
}
