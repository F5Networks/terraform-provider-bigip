/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
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
)

// normalizeGtmBigipDestination ensures the destination is in a consistent format
// for the BIG-IP "bigip" monitor type. The API returns destination as "*:*" by default.
// This function normalizes various input formats to a consistent "address:port" format.
func normalizeGtmBigipDestination(destination string) string {
	if destination == "" {
		return "*:*"
	}
	// If destination doesn't contain a colon, it's just an address - add wildcard port
	if !strings.Contains(destination, ":") {
		return destination + ":*"
	}
	return destination
}

// validateGtmBigipDestination validates that the destination follows BIG-IP API rules:
// - "*:*" is valid (all wildcards)
// - "*:port" is valid (wildcard IP + specific port)
// - "ip:port" is valid (specific IP + specific port)
// - "ip:*" is INVALID (specific IP + wildcard port) - BIG-IP requires a port when an address is specified
func validateGtmBigipDestination(destination string) error {
	if destination == "" || destination == "*:*" {
		return nil
	}
	parts := strings.SplitN(destination, ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid destination format %q: must be in 'address:port' format", destination)
	}
	address := parts[0]
	port := parts[1]
	// If a specific IP address is provided (not wildcard), a specific port is required
	if address != "*" && port == "*" {
		return fmt.Errorf("invalid destination %q: when a specific IP address is provided, a specific port is required (BIG-IP API constraint). Use '*:*' for all wildcards, or specify both IP and port like '%s:80'", destination, address)
	}
	return nil
}

func resourceBigipGtmMonitorBigip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmMonitorBigipCreate,
		ReadContext:   resourceBigipGtmMonitorBigipRead,
		UpdateContext: resourceBigipGtmMonitorBigipUpdate,
		DeleteContext: resourceBigipGtmMonitorBigipDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the GTM BIG-IP monitor",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/bigip",
				Description: "Inherit properties from this monitor",
			},
			"destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the IP address and service port of the resource that is the destination of this monitor. Format: ip:port. Default is \"*:*\"",
				Default:     "*:*",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies, in seconds, the frequency at which the system issues the monitor check",
				Default:     30,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of seconds the target has in which to respond to the monitor request",
				Default:     90,
			},
			"ignore_down_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the monitor ignores a down response from the system it is monitoring",
				Default:     "disabled",
			},
			"aggregation_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies how the system combines monitor information for a monitored pool. The default is none",
				Default:     "none",
			},
		},
	}
}

func resourceBigipGtmMonitorBigipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating GTM BIG-IP Monitor: %s", name)

	destination := normalizeGtmBigipDestination(d.Get("destination").(string))
	if err := validateGtmBigipDestination(destination); err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM BIG-IP Monitor %s: %v", name, err))
	}

	monitor := &bigip.Gtmmonitor{
		Name:                     name,
		Defaults_from:            d.Get("defaults_from").(string),
		Destination:              destination,
		Interval:                 d.Get("interval").(int),
		Timeout:                  d.Get("timeout").(int),
		Ignore_down_response:     d.Get("ignore_down_response").(string),
		Aggregate_dynamic_ratios: d.Get("aggregation_type").(string),
	}

	err := client.CreateGtmMonitor(monitor, "bigip")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM BIG-IP Monitor %s: %v", name, err))
	}

	d.SetId(name)

	return resourceBigipGtmMonitorBigipRead(ctx, d, meta)
}

func resourceBigipGtmMonitorBigipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading GTM BIG-IP Monitor: %s", name)

	monitor, err := client.GetGtmMonitor(name, "bigip")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM BIG-IP Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading GTM BIG-IP Monitor %s: %v", name, err))
	}

	if monitor == nil {
		log.Printf("[WARN] GTM BIG-IP Monitor %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", monitor.FullPath)
	d.Set("defaults_from", monitor.Defaults_from)
	d.Set("destination", normalizeGtmBigipDestination(monitor.Destination))
	d.Set("interval", monitor.Interval)
	d.Set("timeout", monitor.Timeout)
	d.Set("ignore_down_response", monitor.Ignore_down_response)
	d.Set("aggregation_type", monitor.Aggregate_dynamic_ratios)

	return nil
}

func resourceBigipGtmMonitorBigipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating GTM BIG-IP Monitor: %s", name)

	destination := normalizeGtmBigipDestination(d.Get("destination").(string))
	if err := validateGtmBigipDestination(destination); err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM BIG-IP Monitor %s: %v", name, err))
	}

	monitor := &bigip.Gtmmonitor{
		Name:                     name,
		Defaults_from:            d.Get("defaults_from").(string),
		Destination:              destination,
		Interval:                 d.Get("interval").(int),
		Timeout:                  d.Get("timeout").(int),
		Ignore_down_response:     d.Get("ignore_down_response").(string),
		Aggregate_dynamic_ratios: d.Get("aggregation_type").(string),
	}

	err := client.ModifyGtmMonitor(name, monitor, "bigip")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM BIG-IP Monitor %s: %v", name, err))
	}

	return resourceBigipGtmMonitorBigipRead(ctx, d, meta)
}

func resourceBigipGtmMonitorBigipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting GTM BIG-IP Monitor: %s", name)

	err := client.DeleteGtmMonitor(name, "bigip")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM BIG-IP Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting GTM BIG-IP Monitor %s: %v", name, err))
	}

	d.SetId("")
	return nil
}
