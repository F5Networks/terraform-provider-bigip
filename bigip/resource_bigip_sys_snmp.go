/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// this module does not have DELETE function as there is no API for Delete.
func resourceBigipSysSnmp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysSnmpCreate,
		UpdateContext: resourceBigipSysSnmpUpdate,
		ReadContext:   resourceBigipSysSnmpRead,
		DeleteContext: resourceBigipSysSnmpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"sys_contact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Contact Person email",
			},
			"sys_location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Location of the F5 ",
			},
			"allowedaddresses": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of SNMP addresses",
			},
		},
	}

}

func resourceBigipSysSnmpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	sysContact := d.Get("sys_contact").(string)
	sysLocation := d.Get("sys_location").(string)
	allowedAddresses := setToStringSlice(d.Get("allowedaddresses").(*schema.Set))

	log.Println("[INFO] Creating Snmp ")

	err := client.CreateSNMP(
		sysContact,
		sysLocation,
		allowedAddresses,
	)

	if err != nil {
		log.Printf("[ERROR] Unable to Configure SNMP  (%v) ", err)
		return diag.FromErr(err)
	}
	d.SetId(sysContact)
	return resourceBigipSysSnmpRead(ctx, d, meta)
}

func resourceBigipSysSnmpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Updating SNMP " + sysContact)

	r := &bigip.SNMP{
		SysContact:       sysContact,
		SysLocation:      d.Get("sys_location").(string),
		AllowedAddresses: setToStringSlice(d.Get("allowedaddresses").(*schema.Set)),
	}

	err := client.ModifySNMP(r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify SNMP (%s) (%v) ", sysContact, err)
		return diag.FromErr(err)
	}
	return resourceBigipSysSnmpRead(ctx, d, meta)
}

func resourceBigipSysSnmpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	sysContact := d.Id()

	log.Println("[INFO] Reading SNMP " + sysContact)

	snmp, err := client.SNMPs()
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve SNMP  (%v) ", err)
		return diag.FromErr(err)
	}
	if snmp == nil {
		log.Printf("[WARN] SNMP (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("sys_contact", snmp.SysContact)
	_ = d.Set("sys_location", snmp.SysLocation)
	_ = d.Set("allowedaddresses", snmp.AllowedAddresses)

	return nil
}

func resourceBigipSysSnmpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No API support for Delete
	return nil
}
