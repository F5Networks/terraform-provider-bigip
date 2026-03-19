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

func resourceBigipGtmMonitorHttp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmMonitorHttpCreate,
		ReadContext:   resourceBigipGtmMonitorHttpRead,
		UpdateContext: resourceBigipGtmMonitorHttpUpdate,
		DeleteContext: resourceBigipGtmMonitorHttpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the GTM HTTP monitor",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/http",
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
				Default:     120,
			},
			"probe_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of seconds after which the BIG-IP system times out the probe request to the BIG-IP system",
				Default:     5,
			},
			"ignore_down_response": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the monitor ignores a down response from the system it is monitoring",
				Default:     "disabled",
			},
			"transparent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the monitor operates in transparent mode",
				Default:     "disabled",
			},
			"reverse": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instructs the system to mark the target resource down when the test is successful",
				Default:     "disabled",
			},
			"send": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the text string that the monitor sends to the target object",
				Default:     "GET /\\r\\n",
			},
			"receive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the text string that the monitor looks for in the returned resource",
			},
		},
	}
}

func resourceBigipGtmMonitorHttpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating GTM HTTP Monitor: %s", name)

	monitor := &bigip.Gtmmonitor{
		Name:                 name,
		Defaults_from:        d.Get("defaults_from").(string),
		Destination:          d.Get("destination").(string),
		Interval:             d.Get("interval").(int),
		Timeout:              d.Get("timeout").(int),
		Probe_timeout:        d.Get("probe_timeout").(int),
		Ignore_down_response: d.Get("ignore_down_response").(string),
		Transparent:          d.Get("transparent").(string),
		Reverse:              d.Get("reverse").(string),
		Send:                 d.Get("send").(string),
		Recv:                 d.Get("receive").(string),
	}

	err := client.CreateGtmMonitor(monitor, "http")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM HTTP Monitor %s: %v", name, err))
	}

	d.SetId(name)

	return resourceBigipGtmMonitorHttpRead(ctx, d, meta)
}

func resourceBigipGtmMonitorHttpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading GTM HTTP Monitor: %s", name)

	monitor, err := client.GetGtmMonitor(name, "http")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM HTTP Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading GTM HTTP Monitor %s: %v", name, err))
	}

	if monitor == nil {
		log.Printf("[WARN] GTM HTTP Monitor %s not found, removing from state", name)
		d.SetId("")
		return nil
	}

	d.Set("name", monitor.FullPath)
	d.Set("defaults_from", monitor.Defaults_from)
	d.Set("destination", monitor.Destination)
	d.Set("interval", monitor.Interval)
	d.Set("timeout", monitor.Timeout)
	d.Set("probe_timeout", monitor.Probe_timeout)
	d.Set("ignore_down_response", monitor.Ignore_down_response)
	d.Set("transparent", monitor.Transparent)
	d.Set("reverse", monitor.Reverse)
	d.Set("send", monitor.Send)
	d.Set("receive", monitor.Recv)

	return nil
}

func resourceBigipGtmMonitorHttpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating GTM HTTP Monitor: %s", name)

	monitor := &bigip.Gtmmonitor{
		Name:                 name,
		Defaults_from:        d.Get("defaults_from").(string),
		Destination:          d.Get("destination").(string),
		Interval:             d.Get("interval").(int),
		Timeout:              d.Get("timeout").(int),
		Probe_timeout:        d.Get("probe_timeout").(int),
		Ignore_down_response: d.Get("ignore_down_response").(string),
		Transparent:          d.Get("transparent").(string),
		Reverse:              d.Get("reverse").(string),
		Send:                 d.Get("send").(string),
		Recv:                 d.Get("receive").(string),
	}

	err := client.ModifyGtmMonitor(name, monitor, "http")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM HTTP Monitor %s: %v", name, err))
	}

	return resourceBigipGtmMonitorHttpRead(ctx, d, meta)
}

func resourceBigipGtmMonitorHttpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting GTM HTTP Monitor: %s", name)

	err := client.DeleteGtmMonitor(name, "http")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM HTTP Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting GTM HTTP Monitor %s: %v", name, err))
	}

	d.SetId("")
	return nil
}
