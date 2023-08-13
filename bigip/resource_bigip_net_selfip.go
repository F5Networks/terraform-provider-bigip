/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBigipNetSelfIP() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceBigipNetSelfIPCreate,
		ReadContext:   resourceBigipNetSelfIPRead,
		UpdateContext: resourceBigipNetSelfIPUpdate,
		DeleteContext: resourceBigipNetSelfIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the SelfIP",
			},

			"ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "SelfIP IP address",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					old = strings.Replace(old, "%0", "", 1)
					new = strings.Replace(new, "%0", "", 1)
					return old == new
				},
			},

			"vlan": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the vlan",
			},

			"traffic_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the traffic group, defaults to traffic-group-local-only if not specified",
				Default:     "traffic-group-local-only",
			},

			"port_lockdown": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "port lockdown",
			},
		},
	}
}

func resourceBigipNetSelfIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	log.Printf("[INFO] Creating SelfIP %s", name)

	err := client.CreateSelfIP(config)

	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating SelfIP %s: %v ", name, err))
	}

	d.SetId(name)

	return resourceBigipNetSelfIPRead(ctx, d, meta)
}

func resourceBigipNetSelfIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading SelfIP %s", name)

	selfIP, err := client.SelfIP(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving SelfIP %s: %v ", name, err))
	}
	if selfIP == nil {
		log.Printf("[DEBUG] SelfIP %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", selfIP.FullPath)
	_ = d.Set("vlan", selfIP.Vlan)
	_ = d.Set("ip", selfIP.Address)

	// Extract Traffic Group name from the full path (ignoring /Common/ prefix)
	regex := regexp.MustCompile(`\/Common\/(.+)`)
	_ = d.Set("traffic_group", selfIP.TrafficGroup)
	trafficGroup := regex.FindStringSubmatch(selfIP.TrafficGroup)
	if len(trafficGroup) > 0 {
		_ = d.Set("traffic_group", trafficGroup[1])
	}
	if selfIP.AllowService == nil {
		_ = d.Set("port_lockdown", []string{"none"})
	} else {
		_ = d.Set("port_lockdown", selfIP.AllowService)
	}
	return nil
}

func resourceBigipNetSelfIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating SelfIP %s", name)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	err := client.ModifySelfIP(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error modifying SelfIP %s: %v ", name, err))
	}

	return resourceBigipNetSelfIPRead(ctx, d, meta)

}

func resourceBigipNetSelfIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting SelfIP %s", name)

	err := client.DeleteSelfIP(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error deleting SelfIP %s: %v ", name, err))
	}

	d.SetId("")
	return nil
}

func getNetSelfIPConfig(d *schema.ResourceData, config *bigip.SelfIP) *bigip.SelfIP {
	var portLockdown interface{}
	p := d.Get("port_lockdown").([]interface{})

	if len(p) > 0 {
		switch p[0] {
		case "all":
			portLockdown = "all"
		case "none":
			portLockdown = nil
		default:
			portLockdown = p
		}
	}

	config.Address = d.Get("ip").(string)
	config.Vlan = d.Get("vlan").(string)
	config.TrafficGroup = d.Get("traffic_group").(string)
	config.AllowService = portLockdown

	return config
}
