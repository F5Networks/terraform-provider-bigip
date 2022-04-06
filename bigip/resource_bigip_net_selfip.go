/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"regexp"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipNetSelfIP() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetSelfIPCreate,
		Read:   resourceBigipNetSelfIPRead,
		Update: resourceBigipNetSelfIPUpdate,
		Delete: resourceBigipNetSelfIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceBigipNetSelfIPCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	log.Printf("[DEBUG] Creating SelfIP %s", name)

	err := client.CreateSelfIP(config)

	if err != nil {
		return fmt.Errorf("Error creating SelfIP %s: %v", name, err)
	}

	d.SetId(name)

	return resourceBigipNetSelfIPRead(d, meta)
}

func resourceBigipNetSelfIPRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[DEBUG] Reading SelfIP %s", name)

	selfIP, err := client.SelfIP(name)
	if err != nil {
		return fmt.Errorf("Error retrieving SelfIP %s: %v", name, err)
	}
	if selfIP == nil {
		log.Printf("[DEBUG] SelfIP %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", selfIP.FullPath)
	d.Set("vlan", selfIP.Vlan)

	// Extract Self IP address from "(selfip_address)[%route_domain](/mask)" groups 1 + 2
	regex := regexp.MustCompile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?(\/\d+)`)
	selfipAddress := regex.FindStringSubmatch(selfIP.Address)
	log.Printf("[DEBUG] value of selfipAddress: %v", selfipAddress)
	parsedSelfipAddress := selfipAddress[1] + selfipAddress[2]
	d.Set("ip", parsedSelfipAddress)

	// Extract Traffic Group name from the full path (ignoring /Common/ prefix)
	regex = regexp.MustCompile(`\/Common\/(.+)`)
	trafficGroup := regex.FindStringSubmatch(selfIP.TrafficGroup)
	d.Set("traffic_group", trafficGroup[1])

	return nil
}

func resourceBigipNetSelfIPUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating SelfIP %s", name)

	pss := &bigip.SelfIP{
		Name: name,
	}
	config := getNetSelfIPConfig(d, pss)

	err := client.ModifySelfIP(name, config)
	if err != nil {
		return fmt.Errorf("Error modifying SelfIP %s: %v", name, err)
	}

	return resourceBigipNetSelfIPRead(d, meta)

}

func resourceBigipNetSelfIPDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[DEBUG] Deleting SelfIP %s", name)

	err := client.DeleteSelfIP(name)
	if err != nil {
		return fmt.Errorf("Error deleting SelfIP %s: %v", name, err)
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
