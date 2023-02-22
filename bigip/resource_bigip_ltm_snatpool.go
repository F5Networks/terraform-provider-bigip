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

func resourceBigipLtmSnatpool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmSnatpoolCreate,
		UpdateContext: resourceBigipLtmSnatpoolUpdate,
		ReadContext:   resourceBigipLtmSnatpoolRead,
		DeleteContext: resourceBigipLtmSnatpoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "SNAT Pool list Name, format /partition/name. e.g. /Common/snat_pool",
				ValidateFunc: validateF5Name,
			},

			"members": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				MinItems:    1,
				Description: "Specifies a translation address to add to or delete from a SNAT pool, at least one address is required.",
			},
		},
	}
}

func resourceBigipLtmSnatpoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	members := setToStringSlice(d.Get("members").(*schema.Set))

	log.Println("[INFO] Creating SNAT Pool " + name)

	err := client.CreateSnatPool(name, members)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Snat Pool  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceBigipLtmSnatpoolRead(ctx, d, meta)
}

func resourceBigipLtmSnatpoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SNAT Pool " + name)

	r := &bigip.SnatPool{
		Name:    d.Get("name").(string),
		Members: setToStringSlice(d.Get("members").(*schema.Set)),
	}

	err := client.ModifySnatPool(name, r)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify Snat Pool  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}

	return resourceBigipLtmSnatpoolRead(ctx, d, meta)
}

func resourceBigipLtmSnatpoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching SNAT Pool " + name)

	snatpool, err := client.GetSnatPool(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Snat Pool  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if snatpool == nil {
		log.Printf("[WARN] SNAT Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("members", snatpool.Members)

	return nil

}

func resourceBigipLtmSnatpoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	err := client.DeleteSnatPool(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Snat Pool  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
