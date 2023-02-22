/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func resourceBigipFastApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipFastAppCreate,
		ReadContext:   resourceBigipFastAppRead,
		UpdateContext: resourceBigipFastAppUpdate,
		DeleteContext: resourceBigipFastAppDelete,
		Exists:        resourceBigipFastAppExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"fast_json": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "FAST application declaration.",
				StateFunc: func(v interface{}) string {
					fjson, _ := structure.NormalizeJsonString(v)

					return fjson
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
					iterate := make(map[string]interface{})

					for k, v := range oldJsonref {
						iterate[k] = v
					}
					for k1 := range iterate {
						_, ok := newJsonref[k1]
						if !ok {
							delete(oldJsonref, k1)
						}
					}
					jsonEqualityAfter := reflect.DeepEqual(oldJsonref, newJsonref)
					if jsonEqualityAfter {
						return true
					} else {
						return false
					}
				},
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of FAST application template.",
			},
			"tenant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of FAST application tenant.",
			},
			"application": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of FAST application.",
			},
		},
	}
}

func resourceBigipFastAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fastTmpl := d.Get("template").(string)
	fastJson := d.Get("fast_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Creating FastApp config")
	userAgent := fmt.Sprintf("?userAgent=%s/%s", client.UserAgent, fastTmpl)
	tenant, app, err := client.PostFastAppBigip(fastJson, fastTmpl, userAgent)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("tenant", tenant)
	_ = d.Set("application", app)
	log.Printf("[DEBUG] ID for resource :%+v", app)
	d.SetId(d.Get("application").(string))
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
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_fast_application", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipFastAppRead(ctx, d, meta)
}
func resourceBigipFastAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Reading FastApp config")
	name := d.Id()
	tenant := d.Get("tenant").(string)
	log.Printf("[DEBUG] FAST application get call : %s", name)
	fastJson, err := client.GetFastApp(tenant, name)
	log.Printf("[DEBUG] FAST json retreived from the GET call in Read function : %s", fastJson)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if fastJson == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("fast_json", fastJson)
	return nil
}

func resourceBigipFastAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if FastApp config exists in BIGIP")
	name := d.Id()
	tenant := d.Get("tenant").(string)
	fastJson, err := client.GetFastApp(tenant, name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return false, nil
		}
		return false, err
	}
	log.Printf("[INFO] FAST response Body:%+v", fastJson)
	if fastJson == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		return false, nil
	}
	return true, nil
}

func resourceBigipFastAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	fastJson := d.Get("fast_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Updating FastApp Config :%s", fastJson)
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.ModifyFastAppBigip(fastJson, tenant, name)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipFastAppRead(ctx, d, meta)
}

func resourceBigipFastAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.DeleteFastAppBigip(tenant, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
