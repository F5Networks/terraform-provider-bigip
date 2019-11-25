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
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPoolCreate,
		Read:   resourceBigipLtmPoolRead,
		Update: resourceBigipLtmPoolUpdate,
		Delete: resourceBigipLtmPoolDelete,
		Exists: resourceBigipLtmPoolExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the pool",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"monitors": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Computed:    true,
				Optional:    true,
				Description: "Assign monitors to a pool.",
			},

			"allow_nat": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Allow NAT",
			},

			"allow_snat": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Allow SNAT",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"load_balancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Possible values: round-robin, ...",
			},

			"slow_ramp_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Slow ramp time for pool members",
			},

			"service_down_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Possible values: none, reset, reselect, drop",
			},

			"reselect_tries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of times the system tries to select a new pool member after a failure.",
			},
		},
	}
}

func resourceBigipLtmPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	d.SetId(name)
	log.Println("[INFO] Creating pool " + name)
	err := client.CreatePool(name)
	if err != nil {
		return fmt.Errorf("Error retrieving pool (%s): %s", name, err)
	}

	err = resourceBigipLtmPoolUpdate(d, meta)
	if err != nil {
		client.DeletePool(name)
		return err
	}

	return resourceBigipLtmPoolRead(d, meta)
}

func resourceBigipLtmPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	d.Set("name", name)
	log.Println("[INFO] Reading pool " + name)

	pool, err := client.GetPool(name)
	if err != nil {
		return err
	}
	if pool == nil {
		log.Printf("[WARN] Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err := d.Set("allow_nat", pool.AllowNAT); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AllowNAT to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("allow_snat", pool.AllowSNAT); err != nil {
		return fmt.Errorf("[DEBUG] Error saving AllowSNAT to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("load_balancing_mode", pool.LoadBalancingMode); err != nil {
		return fmt.Errorf("[DEBUG] Error saving LoadBalancingMode to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("slow_ramp_time", pool.SlowRampTime); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SlowRampTime to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("service_down_action", pool.ServiceDownAction); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ServiceDownAction to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("reselect_tries", pool.ReselectTries); err != nil {
		return fmt.Errorf("[DEBUG] ERror saving ReselectTries to state for Pool  (%s): %s", d.Id(), err)
	}
	d.Set("description", pool.Description)
	monitors := strings.Split(strings.TrimSpace(pool.Monitor), " and ")
	if err := d.Set("monitors", makeStringSet(&monitors)); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Monitors to state for Pool  (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBigipLtmPoolExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Checking pool " + name + " exists.")

	pool, err := client.GetPool(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Pool   (%s) (%v) ", name, err)
		return false, err
	}

	if pool == nil {
		log.Printf("[WARN] Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return pool != nil, nil
}

func resourceBigipLtmPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	//monitors
	var monitors []string
	if m, ok := d.GetOk("monitors"); ok {
		for _, monitor := range m.(*schema.Set).List() {
			monitors = append(monitors, monitor.(string))
		}
	}

	pool := &bigip.Pool{
		AllowNAT:          d.Get("allow_nat").(string),
		AllowSNAT:         d.Get("allow_snat").(string),
		LoadBalancingMode: d.Get("load_balancing_mode").(string),
		Description:       d.Get("description").(string),
		SlowRampTime:      d.Get("slow_ramp_time").(int),
		ServiceDownAction: d.Get("service_down_action").(string),
		ReselectTries:     d.Get("reselect_tries").(int),
		Monitor:           strings.Join(monitors, " and "),
	}
	err := client.ModifyPool(name, pool)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify Pool   (%s) (%v) ", name, err)
		return err
	}

	return resourceBigipLtmPoolRead(d, meta)
}

func resourceBigipLtmPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting pool " + name)

	err := client.DeletePool(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Pool   (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}
