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

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmSnat() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmSnatCreate,
		Update: resourceBigipLtmSnatUpdate,
		Read:   resourceBigipLtmSnatRead,
		Delete: resourceBigipLtmSnatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Snat list Name",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which partition on BIG-IP",
			},

			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fullpath ",
			},

			"autolasthop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether to automatically map last hop for pools or not. The default is to use next level's defaul",
			},
			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enables or disables mirroring of SNAT connections.",
			},
			"sourceport": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the system preserves the source port of the connection. ",
			},
			"translation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of a translated IP address.",
			},
			"snatpool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of a SNAT pool. You can only use this option when automap and translation are not used",
			},
			"vlansdisabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disables the SNAT on all VLANs.",
			},
			"vlans": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Vlans or Vlan list",
			},

			"origins": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of origin",
						},
						"app_service": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "app service",
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmSnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Snat" + name)

	p := dataToSnat(name, d)
	d.SetId(name)
	err := client.CreateSnat(&p)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Snat  (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmSnatRead(d, meta)
}

func resourceBigipLtmSnatRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching Ltm Snat " + name)
	p, err := client.GetSnat(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Snat  (%s) (%v) ", name, err)
		return err
	}
	if p == nil {
		log.Printf("[WARN] Snat  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("partition", p.Partition)
	if err := d.Set("full_path", p.FullPath); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FullPath to state for Snat  (%s): %s", d.Id(), err)
	}
	if err := d.Set("autolasthop", p.AutoLasthop); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AutoLasthop to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("mirror", p.Mirror)
	if err := d.Set("sourceport", p.SourcePort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SourcePort to state for Snat  (%s): %s", d.Id(), err)
	}
	if err := d.Set("translation", p.Translation); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Translation to state for Snat  (%s): %s", d.Id(), err)
	}

	if err := d.Set("snatpool", p.Snatpool); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Snatpool to state for Snat  (%s): %s", d.Id(), err)
	}
	d.Set("vlansdisabled", p.VlansDisabled)

	if err != nil {
		return err
	}

	return SnatToData(p, d)
}

func resourceBigipLtmSnatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating LtmSnat " + name)
	p := dataToSnat(name, d)
	err := client.UpdateSnat(name, &p)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Snat  (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmSnatRead(d, meta)
}

func resourceBigipLtmSnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteSnat(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Snat  (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func dataToSnat(name string, d *schema.ResourceData) bigip.Snat {
	var p bigip.Snat

	p.Name = name
	p.Partition = d.Get("partition").(string)
	p.FullPath = d.Get("full_path").(string)
	p.AutoLasthop = d.Get("autolasthop").(string)
	p.Mirror = d.Get("mirror").(string)
	p.SourcePort = d.Get("sourceport").(string)
	p.Translation = d.Get("translation").(string)
	p.Snatpool = d.Get("snatpool").(string)
	p.VlansDisabled = d.Get("vlansdisabled").(bool)
	p.Vlans = setToStringSlice(d.Get("vlans").(*schema.Set))
	originsCount := d.Get("origins.#").(int)
	p.Origins = make([]bigip.Originsrecord, 0, originsCount)
	for i := 0; i < originsCount; i++ {
		var r bigip.Originsrecord
		prefix := fmt.Sprintf("origins.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		p.Origins = append(p.Origins, r)
	}

	log.Println("I am in DatatoSnat value of p                                                   ", p)

	return p
}

func SnatToData(p *bigip.Snat, d *schema.ResourceData) error {
	d.Set("partition", p.Partition)
	d.Set("full_path", p.FullPath)
	d.Set("autolasthop", p.AutoLasthop)
	d.Set("mirror", p.Mirror)
	d.Set("sourceport", p.SourcePort)
	d.Set("translation", p.Translation)
	d.Set("snatpool", p.Snatpool)
	d.Set("vlansdisabled", p.VlansDisabled)
	if err := d.Set("vlans", p.Vlans); err != nil {
		return fmt.Errorf("error setting Vlans for resource %s: %s", d.Id(), err)
	}
	for i, r := range p.Origins {
		origins := fmt.Sprintf("origins.%d", i)
		d.Set(fmt.Sprintf("%s.name", origins), r.Name)
	}
	return nil
}
