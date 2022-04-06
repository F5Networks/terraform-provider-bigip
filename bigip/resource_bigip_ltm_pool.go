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
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
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
				ValidateFunc: validateF5NameWithDirectory,
			},
			"monitors": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Computed:    true,
				Optional:    true,
				Description: "Specifies an association between a health or performance monitor and an entire pool, rather than with individual pool members",
			},
			"allow_nat": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether NATs are automatically enabled or disabled for any connections using this pool.",
			},
			"allow_snat": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether SNATs are automatically enabled or disabled for any connections using this pool.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies descriptive text that identifies the pool.",
			},
			"load_balancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the load balancing method. The default is Round Robin.Possible values: round-robin, ...",
			},
			"minimum_active_members": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the system load balances traffic according to the priority number assigned to the pool member,Default Value is 0(disabled)",
			},
			"slow_ramp_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the duration during which the system sends less traffic to a newly-enabled pool member.",
			},
			"service_down_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies how the system should respond when the target pool member becomes unavailable. The default is None, Possible values: [none, reset, reselect, drop]",
			},
			"reselect_tries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of times the system tries to contact a new pool member after a passive failure.",
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
		return fmt.Errorf("Error retrieving pool (%s): %s ", name, err)
	}
	err = resourceBigipLtmPoolUpdate(d, meta)
	if err != nil {
		_ = client.DeletePool(name)
		return err
	}
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
		err = teemDevice.Report(f, "bigip_ltm_pool", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmPoolRead(d, meta)
}

func resourceBigipLtmPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	_ = d.Set("name", name)
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
	if err := d.Set("minimum_active_members", pool.MinActiveMembers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving SlowRampTime to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("service_down_action", pool.ServiceDownAction); err != nil {
		return fmt.Errorf("[DEBUG] Error saving ServiceDownAction to state for Pool  (%s): %s", d.Id(), err)
	}
	if err := d.Set("reselect_tries", pool.ReselectTries); err != nil {
		return fmt.Errorf("[DEBUG] ERror saving ReselectTries to state for Pool  (%s): %s", d.Id(), err)
	}
	_ = d.Set("description", pool.Description)
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
		MinActiveMembers:  d.Get("minimum_active_members").(int),
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
