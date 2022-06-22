/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBigipFastTcpApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastTcpAppCreate,
		Read:   resourceBigipFastTcpAppRead,
		Update: resourceBigipFastTcpAppUpdate,
		Delete: resourceBigipFastTcpAppDelete,
		Exists: resourceBigipFastTcpAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the TCP FAST application",
			},
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the TCP FAST application tenant",
			},
			"virtual_server": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Virtual Server IP. This address combined with virtual port becomes the address to access the application.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Virtual Server Port.",
						},
					},
				},
			},
			"fastl4": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Determine whether to use fastL4 Protocol Profiles.",
						},
						"generate_fastl4_profile": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Determine whether to use FAST-Generated fastL4 Protocol Profile",
						},
						"fastl4_profile_name": {
							Type:          schema.TypeString,
							Optional:      true,
							Default:       "/Common/fastl4",
							ConflictsWith: []string{"generate_fastl4_profile"},
							ValidateFunc: validation.StringInSlice([]string{
								"/Common/fastL4",
								"/Common/apm-forwarding-fastL4",
								"/Common/full-acceleration",
								"/Common/security-fastL4",
							}, false),
							Description: "Select an existing BIG-IP fastL4 profile.",
						},
					},
				},
			},
			"fast_tcp_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Json payload for FAST TCP application.",
			},
		},
	}
}

func resourceBigipFastTcpAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	const templateName string = "bigip-fast-templates/tcp"
	m.Lock()
	defer m.Unlock()

	log.Printf("[INFO] Creating FAST TCP Application")
	cfg := getParamsConfigMap(d)

	userAgent := fmt.Sprintf("?userAgent=%s/%s", client.UserAgent, templateName)
	payload, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}
	tenant, app, err := client.PostFastAppBigip(string(payload), templateName, userAgent)
	if err != nil {
		return err
	}
	d.Set("application", app)
	d.Set("tenant", tenant)
	d.SetId(app)

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

	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	tenant := d.Get("tenant").(string)
	app_name := d.Get("application").(string)

	log.Printf("[INFO] Reading FAST TCP Application config")
	resp, err := client.GetFastApp(tenant, app_name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return nil
		}
		return err
	}
	if resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("fast_tcp_json", resp)
	return nil
}

func resourceBigipFastTcpAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()

	cfg := getParamsConfigMap(d)
	log.Printf("[INFO] Updating FastApp Config :%v", cfg)
	name := d.Get("application").(string)
	tenant := d.Get("tenant").(string)
	payload, err := json.Marshal(cfg)
	if err != nil {
		return nil
	}
	err = client.ModifyFastAppBigip(string(payload), tenant, name)

	if err != nil {
		return err
	}
	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppDelete(d *schema.ResourceData, meta interface{}) error {
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
	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	tenant := d.Get("tenant").(string)
	app_name := d.Get("application").(string)

	log.Printf("[INFO] Reading FAST TCP Application config")
	resp, err := client.GetFastApp(tenant, app_name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return false, nil
		}
		return false, err
	}
	if resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	return true, nil
}

func getParamsConfigMap(d *schema.ResourceData) *map[string]interface{} {
	paramConfig := make(map[string]interface{})

	paramConfig["app_name"] = d.Get("application")
	paramConfig["tenant_name"] = d.Get("tenant")

	if v, ok := d.GetOk("virtual_server"); ok {
		virtual_server := v.(map[string]interface{})
		paramConfig["virtual_address"] = virtual_server["ip"].(string)
		paramConfig["virtual_port"], _ = strconv.Atoi(virtual_server["port"].(string))
	}
	if f4, ok := d.GetOk("fastl4"); ok {
		fastl4_prof := f4.(map[string]interface{})
		paramConfig["fastl4"], _ = strconv.ParseBool(fastl4_prof["enable"].(string))
		if val, ok := fastl4_prof["generate_fastl4_profile"]; ok {
			paramConfig["make_fastl4_profile"], _ = strconv.ParseBool(val.(string))
			if ok, _ := strconv.ParseBool(val.(string)); !ok {
				paramConfig["fastl4_profile_name"] = fastl4_prof["fastl4_profile_name"]
			}
		}
	}

	return &paramConfig
}
