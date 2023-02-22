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
	"regexp"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var parentMonitors = map[string]bool{
	"/Common/udp":           true,
	"/Common/postgresql":    true,
	"/Common/mysql":         true,
	"/Common/mssql":         true,
	"/Common/http":          true,
	"/Common/https":         true,
	"/Common/icmp":          true,
	"/Common/gateway_icmp":  true,
	"/Common/tcp":           true,
	"/Common/tcp_half_open": true,
	"/Common/ftp":           true,
	"/Common/ldap":          true,
	"/Common/smtp":          true,
}

func resourceBigipLtmMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmMonitorCreate,
		ReadContext:   resourceBigipLtmMonitorRead,
		UpdateContext: resourceBigipLtmMonitorUpdate,
		DeleteContext: resourceBigipLtmMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the monitor",
				ForceNew:     true,
				ValidateFunc: validateF5Name,
			},
			"parent": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateParent,
				ForceNew:     true,
				Description:  "Existing monitor to inherit from. Must be one of /Common/http, /Common/https, /Common/icmp, /Common/gateway_icmp or /Common/tcp_half_open or /Common/smtp.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies, in seconds, the frequency at which the system issues the monitor check when either the resource is down or the status of the resource is unknown. The default is 5",
				Computed:    true,
			},
			"up_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the interval for the system to use to perform the health check when a resource is up. The default is 0 (Disabled)",
				Computed:    true,
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of seconds the target has in which to respond to the monitor request. The default is 16 seconds",
				Computed:    true,
			},
			"send": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the text string that the monitor sends to the target object.",
				StateFunc: func(s interface{}) string {
					return strings.ReplaceAll(s.(string), "\r\n", "\\r\\n")
				},
			},
			"receive": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the regular expression representing the text string that the monitor looks for in the returned resource.",
			},

			"receive_disable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The system marks the node or pool member disabled when its response matches Receive Disable String but not Receive String.",
			},

			"reverse": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instructs the system to mark the target resource down when the test is successful.",
				Computed:    true,
			},

			"transparent": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the monitor operates in transparent mode.",
				Computed:    true,
			},

			"manual_resume": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies whether the system automatically changes the status of a resource to Enabled at the next successful monitor check",
				Computed:    true,
			},

			"ip_dscp": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Displays the differentiated services code point (DSCP).The default is 0 (zero)",
				Computed:    true,
			},

			"time_until_up": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the number of seconds to wait after a resource first responds correctly to the monitor before setting the resource to up.",
			},
			"destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Alias for the destination",
			},
			"compatibility": {
				Type:     schema.TypeString,
				Optional: true,
				//Computed:     true,
				Default:      "enabled",
				Description:  "Specifies, when enabled, that the SSL options setting (in OpenSSL) is set to ALL. The default value is enabled.",
				ValidateFunc: validateEnabledDisabled,
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the full path and file name of the file that the system attempts to download. The health check is successful if the system can download the file.",
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the data transfer process (DTP) mode. The default value is passive.",
			},
			"adaptive": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether adaptive response time monitoring is enabled for this monitor. The default is Disabled",
			},
			"adaptive_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the absolute number of milliseconds that may not be exceeded by a monitor probe, regardless of Allowed Divergence",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Specifies the password if the monitored target requires authentication",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the user name if the monitored target requires authentication",
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the database in which your user is created",
			},

			"ssl_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the ssl profile",
			},
		},
	}
}

func resourceBigipLtmMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	parent := monitorParent(d.Get("parent").(string))

	log.Println("[INFO] Creating LTM Monitor " + name + " :: " + parent)
	pss := &bigip.Monitor{
		Name: name,
	}
	config := getLtmMonitorConfig(d, pss)

	if strings.Contains(parent, "gateway") {
		parent = "gateway-icmp"
	}
	if strings.Contains(parent, "tcp_half_open") {
		parent = "tcp-half-open"
	}

	err := client.CreateMonitor(config, parent)

	if err != nil {
		log.Printf("[ERROR] Unable to Create Monitor (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}

	d.SetId(name)
	return resourceBigipLtmMonitorRead(ctx, d, meta)
}

func resourceBigipLtmMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Printf("[INFO] Reading LTM Monitor: %+v", name)

	parentMonitor := d.Get("parent").(string)
	re := regexp.MustCompile("/.*/https$")
	matchresult := re.MatchString(parentMonitor)

	monitors, err := client.Monitors()
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve Monitor (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if monitors == nil {
		log.Printf("[WARN] Monitor (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	for _, m := range monitors {
		if m.FullPath == name {
			_ = d.Set("interval", m.Interval)
			_ = d.Set("up_interval", m.UpInterval)
			_ = d.Set("timeout", m.Timeout)
			if err := d.Set("send", m.SendString); err != nil {
				return diag.FromErr(fmt.Errorf("[DEBUG] Error saving SendString to state for Monitor (%s): %s", d.Id(), err))
			}
			if err := d.Set("receive", m.ReceiveString); err != nil {
				return diag.FromErr(fmt.Errorf("[DEBUG] Error saving ReceiveString to state for Monitor (%s): %s", d.Id(), err))
			}
			_ = d.Set("receive_disable", m.ReceiveDisable)
			_ = d.Set("reverse", m.Reverse)
			_ = d.Set("transparent", m.Transparent)
			_ = d.Set("ip_dscp", m.IPDSCP)
			_ = d.Set("parent", m.ParentMonitor)
			_ = d.Set("time_until_up", m.TimeUntilUp)
			_ = d.Set("manual_resume", m.ManualResume)
			_ = d.Set("destination", m.Destination)
			if matchresult {
				_ = d.Set("compatibility", m.Compatibility)
				_ = d.Set("ssl_profile", m.SSLProfile)
			} else {
				_ = d.Set("compatibility", d.Get("compatibility").(string))
			}
			_ = d.Set("filename", m.Filename)
			_ = d.Set("mode", m.Mode)
			_ = d.Set("adaptive", m.Adaptive)
			_ = d.Set("adaptive_limit", m.AdaptiveLimit)
			_ = d.Set("username", m.Username)
			_ = d.Set("password", m.Password)
			_ = d.Set("name", name)
			_ = d.Set("database", m.Database)
			return nil
		}
	}
	return diag.FromErr(fmt.Errorf("Couldn't find LTM Monitor %s ", name))
}

func resourceBigipLtmMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	pss := &bigip.Monitor{
		Name: name,
	}
	config := getLtmMonitorConfig(d, pss)

	parent := monitorParent(d.Get("parent").(string))

	if strings.Contains(parent, "gateway") {
		parent = "gateway-icmp"
	}
	if strings.Contains(parent, "tcp_half_open") {
		parent = "tcp-half-open"
	}

	err := client.ModifyMonitor(name, parent, config)
	if err != nil {
		log.Printf("[ERROR] Unable to Update Monitor (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}

	return resourceBigipLtmMonitorRead(ctx, d, meta)
}

func resourceBigipLtmMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	parent := monitorParent(d.Get("parent").(string))
	log.Println("[INFO] Deleting monitor " + name + "::" + parent)

	if strings.Contains(parent, "gateway") {
		parent = "gateway-icmp"
	}
	if strings.Contains(parent, "tcp_half_open") {
		parent = "tcp-half-open"
	}

	err := client.DeleteMonitor(name, parent)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Monitor (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func validateParent(v interface{}, k string) ([]string, []error) {
	p := v.(string)
	if parentMonitors[p] {
		return nil, nil
	}

	return nil, []error{fmt.Errorf("parent must be one of /Common/udp, /Common/postgresql, /Common/mysql,/Common/mssql, /Common/http, /Common/https, /Common/icmp, /Common/gateway_icmp, /Common/tcp_half_open, /Common/tcp, /Common/ftp. /Common/smtp")}
}

func monitorParent(s string) string {
	return strings.TrimPrefix(s, "/Common/")
}

func getLtmMonitorConfig(d *schema.ResourceData, config *bigip.Monitor) *bigip.Monitor {
	config.ParentMonitor = d.Get("parent").(string)
	config.Adaptive = d.Get("adaptive").(string)
	config.AdaptiveLimit = d.Get("adaptive_limit").(int)
	config.Compatibility = d.Get("compatibility").(string)
	config.Database = d.Get("database").(string)
	config.Destination = d.Get("destination").(string)
	config.Interval = d.Get("interval").(int)
	config.IPDSCP = d.Get("ip_dscp").(int)
	config.Mode = d.Get("mode").(string)
	config.Filename = d.Get("filename").(string)
	config.ReceiveDisable = d.Get("receive_disable").(string)
	config.ReceiveString = d.Get("receive").(string)
	config.Reverse = d.Get("reverse").(string)
	config.SendString = d.Get("send").(string)
	config.Timeout = d.Get("timeout").(int)
	config.TimeUntilUp = d.Get("time_until_up").(int)
	config.ManualResume = d.Get("manual_resume").(string)
	config.Transparent = d.Get("transparent").(string)
	config.Username = d.Get("username").(string)
	config.Password = d.Get("password").(string)
	config.UpInterval = d.Get("up_interval").(int)
	config.SSLProfile = d.Get("ssl_profile").(string)
	return config
}
