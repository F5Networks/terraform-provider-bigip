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
)

func resourceBigipLtmSnat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmSnatCreate,
		UpdateContext: resourceBigipLtmSnatUpdate,
		ReadContext:   resourceBigipLtmSnatRead,
		DeleteContext: resourceBigipLtmSnatDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the SNAT",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Partition or path to which the SNAT belongs",
			},
			"full_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Fullpath ",
			},
			"autolasthop": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether to automatically map last hop for pools or not. The default is to use next level's defaul",
			},
			"mirror": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enables or disables mirroring of SNAT connections.",
			},
			"sourceport": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "preserve",
				Description: "Specifies how the SNAT object handles the client's source port. The default is Preserve.",
			},
			"translation": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Specifies a particular IP address that you want the SNAT to use as a translation address. When you select IP Address, you also type the IP address in the associated text box",
				ConflictsWith: []string{"snatpool"},
			},
			"snatpool": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Specifies an existing SNAT pool to which you want to map the client IP address. When you select SNAT Pool, you also select an existing SNAT pool from the associated list.",
				ConflictsWith: []string{"translation"},
			},
			"vlansdisabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Disables the SNAT on all VLANs.",
			},
			"vlans": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Specifies the VLANs or tunnels for which the SNAT is enabled or disabled. The default is All",
			},
			"origins": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies, for each SNAT that you create, the origin addresses that are to be members of that SNAT. Specify origin addresses by their IP addresses and service ports",
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

func resourceBigipLtmSnatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Snat: " + name)

	p := dataToSnat(name, d)
	err := client.CreateSnat(&p)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Snat  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipLtmSnatRead(ctx, d, meta)
}

func resourceBigipLtmSnatRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Fetching Ltm Snat:%+v", name)
	p, err := client.GetSnat(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Snat  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if p == nil {
		log.Printf("[WARN] Snat  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", p.FullPath)
	_ = d.Set("autolasthop", p.AutoLasthop)
	_ = d.Set("mirror", p.Mirror)
	_ = d.Set("sourceport", p.SourcePort)

	_ = d.Set("translation", p.Translation)

	_ = d.Set("snatpool", p.Snatpool)

	return SnatToData(ctx, p, d)
}

func resourceBigipLtmSnatUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Updating Ltm Snat:%+v", name)
	p := dataToSnat(name, d)
	err := client.UpdateSnat(name, &p)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Snat  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipLtmSnatRead(ctx, d, meta)
}

func resourceBigipLtmSnatDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Printf("[INFO] Deleting Ltm Snat:%+v", name)
	err := client.DeleteSnat(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Snat  (%s) (%v) ", name, err)
		return diag.FromErr(err)
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
	if d.Get("vlansdisabled").(bool) {
		p.VlansDisabled = d.Get("vlansdisabled").(bool)
	} else {
		p.VlansEnabled = true
	}
	p.Vlans = setToStringSlice(d.Get("vlans").(*schema.Set))
	originsCount := d.Get("origins.#").(int)
	p.Origins = make([]bigip.Originsrecord, 0, originsCount)
	for i := 0; i < originsCount; i++ {
		var r bigip.Originsrecord
		prefix := fmt.Sprintf("origins.%d", i)
		r.Name = d.Get(prefix + ".name").(string)
		p.Origins = append(p.Origins, r)
	}
	return p
}

func SnatToData(ctx context.Context, p *bigip.Snat, d *schema.ResourceData) diag.Diagnostics {
	_ = d.Set("autolasthop", p.AutoLasthop)
	_ = d.Set("mirror", p.Mirror)
	_ = d.Set("sourceport", p.SourcePort)
	_ = d.Set("translation", p.Translation)
	_ = d.Set("snatpool", p.Snatpool)
	if p.VlansDisabled {
		_ = d.Set("vlansdisabled", p.VlansDisabled)
	}
	if p.VlansEnabled {
		_ = d.Set("vlansdisabled", false)
	}
	_ = d.Set("vlans", p.Vlans)
	for i, r := range p.Origins {
		origins := fmt.Sprintf("origins.%d", i)
		_ = d.Set(fmt.Sprintf("%s.name", origins), r.Name)
	}
	return nil
}
