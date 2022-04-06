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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipSysDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysDnsCreate,
		Update: resourceBigipSysDnsUpdate,
		Read:   resourceBigipSysDnsRead,
		Delete: resourceBigipSysDnsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Dns Servers",
				ValidateFunc: validateF5Name,
			},

			"name_servers": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

			"number_of_dots": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "how many DNS Servers",
			},

			"search": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers search domain",
			},
		},
	}

}

func resourceBigipSysDnsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	nameservers := setToStringSlice(d.Get("name_servers").(*schema.Set))
	numberofdots := d.Get("number_of_dots").(int)
	search := setToStringSlice(d.Get("search").(*schema.Set))

	log.Println("[INFO] Creating Dns ")

	err := client.CreateDNS(
		description,
		nameservers,
		numberofdots,
		search,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Create DNS (%s) (%v) ", description, err)
		return err
	}
	d.SetId(description)

	return resourceBigipSysDnsRead(d, meta)
}

func resourceBigipSysDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating DNS " + description)

	r := &bigip.DNS{
		Description:  description,
		NameServers:  setToStringSlice(d.Get("name_servers").(*schema.Set)),
		NumberOfDots: d.Get("number_of_dots").(int),
		Search:       setToStringSlice(d.Get("search").(*schema.Set)),
	}

	err := client.ModifyDNS(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify DNS (%s) (%v) ", description, err)
		return err
	}
	return resourceBigipSysDnsRead(d, meta)
}

func resourceBigipSysDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading DNS " + description)

	dns, err := client.DNSs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve DNS (%s) (%v) ", description, err)
		return err
	}
	if dns == nil {
		log.Printf("[WARN] DNS (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("description", dns.Description)

	if err := d.Set("name_servers", dns.NameServers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Name Servers to state for DNS (%s): %s", d.Id(), err)
	}

	d.Set("number_of_dots", dns.NumberOfDots)

	if err := d.Set("search", dns.Search); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Search  to state for DNS (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBigipSysDnsDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no Delete API for this operation

	return nil
}
