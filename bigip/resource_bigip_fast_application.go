/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"log"
	"reflect"
)

func resourceBigipFastApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastAppCreate,
		Read:   resourceBigipFastAppRead,
		Update: resourceBigipFastAppUpdate,
		Delete: resourceBigipFastAppDelete,
		Exists: resourceBigipFastAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
					old_resp := []byte(old)
					new_resp := []byte(new)
					old_jsonRef := make(map[string]interface{})
					new_jsonRef := make(map[string]interface{})
					json.Unmarshal(old_resp, &old_jsonRef)
					json.Unmarshal(new_resp, &new_jsonRef)
					json_equality_before := reflect.DeepEqual(old_jsonRef, new_jsonRef)
					if json_equality_before == true {
						return true
					}
					iterate := make(map[string]interface{})

					for k, v := range old_jsonRef {
						iterate[k] = v
					}
					for k1 := range iterate {
						_, ok := new_jsonRef[k1]
						if !ok {
							delete(old_jsonRef, k1)
						}
					}
					json_equality_after := reflect.DeepEqual(old_jsonRef, new_jsonRef)
					if json_equality_after == true {
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

func resourceBigipFastAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastTmpl := d.Get("template").(string)
	fastJson := d.Get("fast_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Creating FastApp config")
	tenant, app, err := client.PostFastAppBigip(fastJson, fastTmpl)
	if err != nil {
		return err
	}
	_ = d.Set("tenant", tenant)
	_ = d.Set("application", app)
	log.Printf("[DEBUG] ID for resource :%+v", app)
	d.SetId(d.Get("application").(string))
	return resourceBigipFastAppRead(d, meta)
}
func resourceBigipFastAppRead(d *schema.ResourceData, meta interface{}) error {
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
		return err
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

func resourceBigipFastAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastJson := d.Get("fast_json").(string)
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Updating FastApp Config :%s", fastJson)
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.ModifyFastAppBigip(fastJson, tenant, name)
	if err != nil {
		return err
	}
	return resourceBigipFastAppRead(d, meta)
}

func resourceBigipFastAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.DeleteFastAppBigip(tenant, name)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
