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

func resourceBigipGtmMonitorPostgresql() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipGtmMonitorPostgresqlCreate,
		ReadContext:   resourceBigipGtmMonitorPostgresqlRead,
		UpdateContext: resourceBigipGtmMonitorPostgresqlUpdate,
		DeleteContext: resourceBigipGtmMonitorPostgresqlDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Name of the GTM PostgreSQL monitor",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"defaults_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/postgresql",
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
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the database in which the user is created",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user name if the monitored target requires authentication",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Specifies the password if the monitored target requires authentication",
			},
			"probe_count": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of monitor probes after which the system times out",
			},
			"debug": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the monitor sends error messages and additional information to a log file created and labeled specifically for this monitor",
				Default:     "no",
			},
			"receive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the text string that the monitor looks for in the returned resource",
			},
		},
	}
}

func resourceBigipGtmMonitorPostgresqlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Printf("[INFO] Creating GTM PostgreSQL Monitor: %s", name)

	monitor := &bigip.Gtmmonitor{
		Name:                 name,
		Defaults_from:        d.Get("defaults_from").(string),
		Destination:          d.Get("destination").(string),
		Interval:             d.Get("interval").(int),
		Timeout:              d.Get("timeout").(int),
		Probe_timeout:        d.Get("probe_timeout").(int),
		Ignore_down_response: d.Get("ignore_down_response").(string),
		Database:             d.Get("database").(string),
		Username:             d.Get("username").(string),
		Password:             d.Get("password").(string),
		Count:                d.Get("probe_count").(string),
		Debug:                d.Get("debug").(string),
		Recv:                 d.Get("receive").(string),
	}

	err := client.CreateGtmMonitor(monitor, "postgresql")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating GTM PostgreSQL Monitor %s: %v", name, err))
	}

	d.SetId(name)

	return resourceBigipGtmMonitorPostgresqlRead(ctx, d, meta)
}

func resourceBigipGtmMonitorPostgresqlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Reading GTM PostgreSQL Monitor: %s", name)

	monitor, err := client.GetGtmMonitor(name, "postgresql")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM PostgreSQL Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading GTM PostgreSQL Monitor %s: %v", name, err))
	}

	if monitor == nil {
		log.Printf("[WARN] GTM PostgreSQL Monitor %s not found, removing from state", name)
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
	d.Set("database", monitor.Database)
	d.Set("username", monitor.Username)
	// Password is sensitive, don't set it back
	d.Set("probe_count", monitor.Count)
	d.Set("debug", monitor.Debug)
	d.Set("receive", monitor.Recv)

	return nil
}

func resourceBigipGtmMonitorPostgresqlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Updating GTM PostgreSQL Monitor: %s", name)

	monitor := &bigip.Gtmmonitor{
		Name:                 name,
		Defaults_from:        d.Get("defaults_from").(string),
		Destination:          d.Get("destination").(string),
		Interval:             d.Get("interval").(int),
		Timeout:              d.Get("timeout").(int),
		Probe_timeout:        d.Get("probe_timeout").(int),
		Ignore_down_response: d.Get("ignore_down_response").(string),
		Database:             d.Get("database").(string),
		Username:             d.Get("username").(string),
		Password:             d.Get("password").(string),
		Count:                d.Get("probe_count").(string),
		Debug:                d.Get("debug").(string),
		Recv:                 d.Get("receive").(string),
	}

	err := client.ModifyGtmMonitor(name, monitor, "postgresql")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating GTM PostgreSQL Monitor %s: %v", name, err))
	}

	return resourceBigipGtmMonitorPostgresqlRead(ctx, d, meta)
}

func resourceBigipGtmMonitorPostgresqlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Printf("[INFO] Deleting GTM PostgreSQL Monitor: %s", name)

	err := client.DeleteGtmMonitor(name, "postgresql")
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Printf("[WARN] GTM PostgreSQL Monitor %s not found, removing from state", name)
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting GTM PostgreSQL Monitor %s: %v", name, err))
	}

	d.SetId("")
	return nil
}
