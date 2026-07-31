/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipNetRouteDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipNetRouteDomainCreate,
		ReadContext:   resourceBigipNetRouteDomainRead,
		UpdateContext: resourceBigipNetRouteDomainUpdate,
		DeleteContext: resourceBigipNetRouteDomainDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5Name,
				Description:  "Name of the route domain",
			},
			"id_": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  "Unique numeric ID of the route domain (0–65534)",
			},
			"strict": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enforce strict isolation between route domains",
			},
			"vlans": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "VLANs associated with this route domain",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateF5Name,
				},
			},
		},
	}
}

// getRouteDomain fetches a single route domain by full path or name without
// requiring changes to the upstream go-bigip library.
func getRouteDomain(client *bigip.BigIP, name string) (*bigip.RouteDomain, error) {
	all, err := client.RouteDomains()
	if err != nil {
		return nil, err
	}
	for i := range all.RouteDomains {
		rd := &all.RouteDomains[i]
		if rd.FullPath == name || rd.Name == name {
			return rd, nil
		}
	}
	return nil, nil
}

func resourceBigipNetRouteDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	id := d.Get("id_").(int)
	strict := d.Get("strict").(bool)

	vlansRaw := d.Get("vlans").(*schema.Set).List()
	vlanParts := make([]string, 0, len(vlansRaw))
	for _, v := range vlansRaw {
		vlanParts = append(vlanParts, v.(string))
	}
	vlansStr := strings.Join(vlanParts, ",")

	log.Printf("[INFO] Creating Route Domain %s", name)
	err := client.CreateRouteDomain(name, id, strict, vlansStr)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating Route Domain %s: %v", name, err))
	}

	d.SetId(name)
	return resourceBigipNetRouteDomainRead(ctx, d, meta)
}

func resourceBigipNetRouteDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading Route Domain %s", name)
	rd, err := getRouteDomain(client, name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading Route Domain %s: %v", name, err))
	}
	if rd == nil {
		log.Printf("[WARN] Route Domain %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	_ = d.Set("name", rd.FullPath)
	_ = d.Set("id_", rd.ID)
	_ = d.Set("strict", rd.Strict == "enabled")
	_ = d.Set("vlans", rd.Vlans)

	return nil
}

func resourceBigipNetRouteDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating Route Domain %s", name)

	strict := "disabled"
	if d.Get("strict").(bool) {
		strict = "enabled"
	}

	vlansRaw := d.Get("vlans").(*schema.Set).List()
	vlansList := make([]string, 0, len(vlansRaw))
	for _, v := range vlansRaw {
		vlansList = append(vlansList, v.(string))
	}

	config := &bigip.RouteDomain{
		ID:     d.Get("id_").(int),
		Strict: strict,
		Vlans:  vlansList,
	}

	err := client.ModifyRouteDomain(name, config)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating Route Domain %s: %v", name, err))
	}

	return resourceBigipNetRouteDomainRead(ctx, d, meta)
}

func resourceBigipNetRouteDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting Route Domain %s", name)
	err := client.DeleteRouteDomain(name)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting Route Domain %s: %v", name, err))
	}

	d.SetId("")
	return nil
}
