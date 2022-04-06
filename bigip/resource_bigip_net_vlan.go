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

func resourceBigipNetVlan() *schema.Resource {

	return &schema.Resource{
		Create: resourceBigipNetVlanCreate,
		Read:   resourceBigipNetVlanRead,
		Update: resourceBigipNetVlanUpdate,
		Delete: resourceBigipNetVlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}

}

func resourceBigipNetVlanCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	tag := d.Get("tag").(int)

	log.Printf("[DEBUG] Creating VLAN %s", name)

	d.Partial(true)

	err := client.CreateVlan(
		name,
		tag,
	)

	if err != nil {
		return fmt.Errorf("Error creating VLAN %s: %v", name, err)
	}

	d.SetId(name)
	d.SetPartial("name")
	d.SetPartial("tag")

	ifaceCount := d.Get("interfaces.#").(int)
	for i := 0; i < ifaceCount; i++ {
		prefix := fmt.Sprintf("interfaces.%d", i)
		iface := d.Get(prefix + ".vlanport").(string)
		tagged := d.Get(prefix + ".tagged").(bool)

		err = client.AddInterfaceToVlan(name, iface, tagged)
		if err != nil {
			return fmt.Errorf("Error adding Interface %s to VLAN %s: %v", iface, name, err)
		}
	}
	d.SetPartial("interfaces")

	d.Partial(false)

	return resourceBigipNetVlanRead(d, meta)
}

func resourceBigipNetVlanRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Reading VLAN %s", name)

	vlan, err := client.Vlan(name)
	if err != nil {
		return fmt.Errorf("Error retrieving VLAN %s: %v", name, err)
	}
	if vlan == nil {
		log.Printf("[DEBUG] VLAN %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", vlan.FullPath)
	d.Set("tag", vlan.Tag)

	log.Printf("[DEBUG] Reading VLAN %s Interfaces", name)

	vlanInterfaces, err := client.GetVlanInterfaces(name)
	if err != nil {
		return fmt.Errorf("Error retrieving VLAN %s Interfaces: %v", name, err)
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
		return fmt.Errorf("Error updating Interfaces in state for VLAN %s: %v", name, err)
	}

	return nil
}

func resourceBigipNetVlanUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Updating VLAN %s", name)

	r := &bigip.Vlan{
		Name: name,
		Tag:  d.Get("tag").(int),
	}

	err := client.ModifyVlan(name, r)
	if err != nil {
		return fmt.Errorf("Error modifying VLAN %s: %v", name, err)
	}

	return resourceBigipNetVlanRead(d, meta)
}

func resourceBigipNetVlanDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[DEBUG] Deleting VLAN %s", name)

	err := client.DeleteVlan(name)
	if err != nil {
		return fmt.Errorf("Error Deleting Vlan : %s", err)
	}

	d.SetId("")
	return nil
}
