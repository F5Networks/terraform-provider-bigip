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

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipNetVlan() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceBigipNetVlanCreate,
		ReadContext:   resourceBigipNetVlanRead,
		UpdateContext: resourceBigipNetVlanUpdate,
		DeleteContext: resourceBigipNetVlanDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the VLAN",
			},

			"tag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "VLAN ID (tag)",
			},

			"interfaces": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Interface(s) attached to the VLAN",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vlanport": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Vlan name",
						},

						"tagged": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Interface tagged",
						},
					},
				},
			},
			"cmp_hash": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"default", "src-ip", "dst-ip"}, false),
				Description:  "Specifies how the traffic on the VLAN will be disaggregated. The value selected determines the traffic disaggregation method",
			},
		},
	}

}

func resourceBigipNetVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	tag := d.Get("tag").(int)

	log.Printf("[INFO] Creating VLAN %s", name)

	d.Partial(true)

	r := &bigip.Vlan{
		Name:    name,
		Tag:     tag,
		CMPHash: d.Get("cmp_hash").(string),
	}

	err := client.CreateVlan(r)

	if err != nil {
		return diag.FromErr(fmt.Errorf("Error creating VLAN %s: %v ", name, err))
	}

	d.SetId(name)

	ifaceCount := d.Get("interfaces.#").(int)
	for i := 0; i < ifaceCount; i++ {
		prefix := fmt.Sprintf("interfaces.%d", i)
		iface := d.Get(prefix + ".vlanport").(string)
		tagged := d.Get(prefix + ".tagged").(bool)

		err = client.AddInterfaceToVlan(name, iface, tagged)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error adding Interface %s to VLAN %s: %v", iface, name, err))
		}
	}

	d.Partial(false)

	return resourceBigipNetVlanRead(ctx, d, meta)
}

func resourceBigipNetVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Reading VLAN %s", name)

	vlan, err := client.Vlan(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving VLAN %s: %v", name, err))
	}
	if vlan == nil {
		log.Printf("[DEBUG] VLAN %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", vlan.FullPath)
	_ = d.Set("tag", vlan.Tag)
	_ = d.Set("cmp_hash", vlan.CMPHash)

	log.Printf("[DEBUG] Reading VLAN %s Interfaces", name)

	vlanInterfaces, err := client.GetVlanInterfaces(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving VLAN %s Interfaces: %v", name, err))
	}

	var interfaces []map[string]interface{}
	var ifaceTagged bool
	for _, iface := range vlanInterfaces.VlanInterfaces {
		if iface.Tagged {
			ifaceTagged = true
		} else {
			ifaceTagged = false
		}
		log.Printf("[DEBUG] Retrieved VLAN Interface %s, tagging is set to %t", iface.Name, ifaceTagged)

		vlanIface := map[string]interface{}{
			"vlanport": iface.Name,
			"tagged":   ifaceTagged,
		}

		interfaces = append(interfaces, vlanIface)
	}

	if err := d.Set("interfaces", interfaces); err != nil {
		return diag.FromErr(fmt.Errorf("error updating Interfaces in state for VLAN %s: %v", name, err))
	}

	return nil
}

func resourceBigipNetVlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Updating VLAN %s", name)

	r := &bigip.Vlan{
		Name:    name,
		Tag:     d.Get("tag").(int),
		CMPHash: d.Get("cmp_hash").(string),
	}

	err := client.ModifyVlan(name, r)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error modifying VLAN %s: %v", name, err))
	}

	return resourceBigipNetVlanRead(ctx, d, meta)
}

func resourceBigipNetVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Deleting VLAN %s", name)

	err := client.DeleteVlan(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error Deleting Vlan : %s", err))
	}
	d.SetId("")
	return nil
}
