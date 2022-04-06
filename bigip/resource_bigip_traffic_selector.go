/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipTrafficselector() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipTrafficselectorCreate,
		Read:   resourceBigipTrafficselectorRead,
		Update: resourceBigipTrafficselectorUpdate,
		Delete: resourceBigipTrafficselectorDelete,
		Exists: resourceBigipTrafficselectorExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the name of the traffic selector",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the traffic selector.",
			},
			"destination_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the host or network IP address to which the application traffic is destined.When creating a new traffic selector, this parameter is required",
			},
			"destination_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the IP port used by the application. The default value is All Ports",
			},
			"source_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the host or network IP address from which the application traffic originates.When creating a new traffic selector, this parameter is required.",
			},
			"source_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the IP port used by the application. The default value is All Ports",
			},
			"direction": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the traffic selector applies to inbound or outbound traffic, or both. The default value is Both.",
			},
			"ipsec_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateF5Name,
				Computed:     true,
				Description:  "Specifies the IPsec policy that tells the BIG-IP system how to handle the packets.When creating a new traffic selector, if this parameter is not specified, the default is default-ipsec-policy.",
			},
			"order": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Specifies the order in which traffic is matched, if traffic can be matched to multiple traffic selectors.Traffic is matched to the traffic selector with the highest priority (lowest order number)." +
					"When creating a new traffic selector, if this parameter is not specified, the default is last.",
			},
			"ip_protocol": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the network protocol to use for this traffic. The default value is All Protocols (255)",
			},
		},
	}
}

func resourceBigipTrafficselectorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating IPSec traffic Selector " + name)

	pss := &bigip.TrafficSelector{
		Name: name,
	}
	selectorConfig := getTrafficSelectorConfig(d, pss)

	err := client.CreateTrafficSelector(selectorConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Create IPsec Traffic Selector (%s) (%v)", name, err)
		return err
	}
	d.SetId(name)
	return resourceBigipTrafficselectorRead(d, meta)
}

func resourceBigipTrafficselectorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Reading Traffic Selector :%+v", name)
	ts, err := client.GetTrafficselctor(name)
	if err != nil {
		return err
	}
	if ts == nil {
		d.SetId("")
		return fmt.Errorf("[ERROR] Traffic-selctor (%s) not found, removing from state", d.Id())
	}
	log.Printf("[DEBUG] Traffic Selector:%+v", ts)
	if err := d.Set("ip_protocol", ts.IPProtocol); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("destination_address", ts.DestinationAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("source_address", ts.SourceAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("ipsec_policy", ts.IpsecPolicy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("order", ts.Order); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("destination_port", ts.DestinationPort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("source_port", ts.SourcePort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	if err := d.Set("direction", ts.Direction); err != nil {
		return fmt.Errorf("[DEBUG] Error saving IPProtocol to state for Traffic selector (%s): %s", d.Id(), err)
	}
	_ = d.Set("description", ts.Description)
	_ = d.Set("name", name)
	return nil
}

func resourceBigipTrafficselectorExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Check existence of Traffic Selector: %+v ", name)
	ts, err := client.GetTrafficselctor(name)
	if err != nil {
		return false, err
	}
	if ts == nil {
		log.Printf("[WARN] Traffic-selctor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, fmt.Errorf("[ERROR] Traffic-selctor (%s) not found, removing from state", d.Id())
	}
	return true, nil
}

func resourceBigipTrafficselectorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Traffic Selector:%+v ", name)
	pss := &bigip.TrafficSelector{
		Name: name,
	}
	config := getTrafficSelectorConfig(d, pss)

	err := client.ModifyTrafficSelector(name, config)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify IPSec Traffic Selector   (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipTrafficselectorRead(d, meta)
}
func resourceBigipTrafficselectorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Deleting Traffic Selector :%+v ", name)
	err := client.DeleteTrafficSelector(name)
	if err != nil {
		return fmt.Errorf("[ERROR] Unable to Delete Traffic Selector (%s) (%v) ", name, err)
	}
	d.SetId("")
	return nil
}

func getTrafficSelectorConfig(d *schema.ResourceData, config *bigip.TrafficSelector) *bigip.TrafficSelector {
	config.DestinationAddress = d.Get("destination_address").(string)
	config.DestinationPort = d.Get("destination_port").(int)
	config.Direction = d.Get("direction").(string)
	config.IPProtocol = d.Get("ip_protocol").(int)
	config.IpsecPolicy = d.Get("ipsec_policy").(string)
	config.Order = d.Get("order").(int)
	config.SourceAddress = d.Get("source_address").(string)
	config.SourcePort = d.Get("source_port").(int)
	return config
}
