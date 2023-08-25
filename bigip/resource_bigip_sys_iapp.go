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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipSysIapp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysIappCreate,
		UpdateContext: resourceBigipSysIappUpdate,
		ReadContext:   resourceBigipSysIappRead,
		DeleteContext: resourceBigipSysIappDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"jsonfile": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Iapp which needs to be Iappensed",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address of the Iapp which needs to be Iappensed",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Common",
				Description: "Address of the Iapp which needs to be Iappensed",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address of the Iapp which needs to be Iappensed",
			},
			"devicegroup": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BIG-IP password",
			},
			"execute_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "BIG-IP password",
			},
			"inherited_devicegroup": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "true",
				Description: "BIG-IP password",
			},
			"inherited_traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "true",
				Description: "BIG-IP password",
			},
			"strict_updates": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "enabled",
				Description: "BIG-IP password",
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"template_modified": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "no",
				Description: "BIG-IP password",
			},
			"template_prerequisite_errors": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/traffic-group-1",
				Description: "BIG-IP password",
			},
			"lists": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "no",
							Description: "Name of origin",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},
					},
				},
			},

			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"persists": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "true",
							Description: "Name of origin",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},
					},
				},
			},

			"tables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},
						"column_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},

						"encrypted_columns": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},

						"rows": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"row": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},

			"variables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},

						"encrypted": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "no",
							Description: "Name of origin",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},
					},
				},
			},
		},
	}
}
func resourceBigipSysIappCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)

	log.Println("[INFO] Creating Iapp       " + name)
	p := dataToIapp(d)
	err := client.CreateIapp(&p)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Iapp  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipSysIappRead(ctx, d, meta)
}

func resourceBigipSysIappUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating Iapp " + name)
	p := dataToIapp(d)
	err := client.UpdateIapp(name, &p)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Iapp  (%s) ", err)
		return diag.FromErr(err)
	}
	return resourceBigipSysIappRead(ctx, d, meta)

}

func resourceBigipSysIappRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	partition := d.Get("partition").(string)

	log.Println("[INFO] Reading Iapp " + name)

	p, err := client.Iapp(name, partition)
	log.Printf("[INFO] Iapp Info:%+v", p)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Iapp  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	if p == nil {
		log.Printf("[WARN] IApp (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("partition", p.Partition)
	_ = d.Set("devicegroup", p.DeviceGroup)
	_ = d.Set("inherited_devicegroup", p.InheritedDevicegroup)
	_ = d.Set("inherited_traffic_group", p.InheritedTrafficGroup)
	_ = d.Set("strict_updates", p.StrictUpdates)
	_ = d.Set("template_modified", p.TemplateModified)
	_ = d.Set("template_prerequisite_errors", p.TemplatePrerequisiteErrors)
	_ = d.Set("traffic_group", p.TrafficGroup)
	_ = d.Set("lists", p.Lists)
	_ = d.Set("metadata", p.Metadata)
	return nil
}

func resourceBigipSysIappDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	partition := d.Get("partition").(string)
	err := client.DeleteIapp(name, partition)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Iapp  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func dataToIapp(d *schema.ResourceData) bigip.Iapp {
	var p bigip.Iapp
	jsonblob := []byte(d.Get("jsonfile").(string))
	err := json.Unmarshal(jsonblob, &p)
	if err != nil {
		fmt.Println("error", err)
	}
	if _, ok := d.GetOk("execute_action"); ok {
		p.ExecuteAction = d.Get("execute_action").(string)
	}
	if _, ok := d.GetOk("partition"); ok {
		p.Partition = d.Get("partition").(string)
	}
	if _, ok := d.GetOk("template_modified"); ok {
		p.TemplateModified = d.Get("template_modified").(string)
	}
	if _, ok := d.GetOk("strict_updates"); ok {
		p.StrictUpdates = d.Get("strict_updates").(string)
	}
	if _, ok := d.GetOk("description"); ok {
		p.Description = d.Get("description").(string)
	}
	if _, ok := d.GetOk("inherited_devicegroup"); ok {
		p.InheritedDevicegroup = d.Get("inherited_devicegroup").(string)
	}
	if _, ok := d.GetOk("inherited_traffic_group"); ok {
		p.InheritedTrafficGroup = d.Get("inherited_traffic_group").(string)
	}
	return p
}
