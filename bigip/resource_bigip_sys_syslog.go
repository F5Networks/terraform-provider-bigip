/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"log"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// sysSyslogRemoteServer mirrors one entry of sys/syslog's remoteServers
// sub-collection. Defined locally (rather than reusing go-bigip's
// RemoteServer) so its JSON tags stay independent of the vendored type.
type sysSyslogRemoteServer struct {
	Name       string `json:"name"`
	Host       string `json:"host"`
	LocalIP    string `json:"localIp,omitempty"`
	RemotePort int    `json:"remotePort,omitempty"`
}

// sysSyslogConfig is the request/response body for sys/syslog. go-bigip's
// own Syslog type cannot be used here: its MarshalJSON builds an empty,
// never-populated DTO (a vendor bug), so client.CreateSyslog/ModifySyslog
// always send "{}" regardless of the configured value. This resource talks
// to sys/syslog directly via restGet/restPatch instead (see
// bigip/rest_helpers.go).
type sysSyslogConfig struct {
	AuthPrivFrom         string                  `json:"authPrivFrom,omitempty"`
	AuthPrivTo           string                  `json:"authPrivTo,omitempty"`
	ClusteredHostSlot    string                  `json:"clusteredHostSlot,omitempty"`
	ClusteredMessageSlot string                  `json:"clusteredMessageSlot,omitempty"`
	ConsoleLog           string                  `json:"consoleLog,omitempty"`
	CronFrom             string                  `json:"cronFrom,omitempty"`
	CronTo               string                  `json:"cronTo,omitempty"`
	DaemonFrom           string                  `json:"daemonFrom,omitempty"`
	DaemonTo             string                  `json:"daemonTo,omitempty"`
	IsoDate              string                  `json:"isoDate,omitempty"`
	KernFrom             string                  `json:"kernFrom,omitempty"`
	KernTo               string                  `json:"kernTo,omitempty"`
	Local6From           string                  `json:"local6From,omitempty"`
	Local6To             string                  `json:"local6To,omitempty"`
	MailFrom             string                  `json:"mailFrom,omitempty"`
	MailTo               string                  `json:"mailTo,omitempty"`
	MessagesFrom         string                  `json:"messagesFrom,omitempty"`
	MessagesTo           string                  `json:"messagesTo,omitempty"`
	UserLogFrom          string                  `json:"userLogFrom,omitempty"`
	UserLogTo            string                  `json:"userLogTo,omitempty"`
	RemoteServers        []sysSyslogRemoteServer `json:"remoteServers"`
}

func resourceBigipSysSyslog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipSysSyslogCreate,
		ReadContext:   resourceBigipSysSyslogRead,
		UpdateContext: resourceBigipSysSyslogUpdate,
		DeleteContext: resourceBigipSysSyslogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Configures the BIG-IP system's syslog settings (log levels and remote syslog servers). This is a singleton resource: there is exactly one syslog configuration per BIG-IP, so Delete resets it to BIG-IP's factory defaults rather than removing an object.",
		Schema: map[string]*schema.Schema{
			"auth_priv_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the auth/priv facility to include in the log.",
			},
			"auth_priv_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the auth/priv facility to include in the log.",
			},
			"console_log": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether logs are provided to the console. The default is enabled.",
			},
			"cron_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the cron facility to include in the log.",
			},
			"cron_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the cron facility to include in the log.",
			},
			"daemon_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the daemon facility to include in the log.",
			},
			"daemon_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the daemon facility to include in the log.",
			},
			"iso_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies whether the date included in a log message conforms to the ISO 8601 standard.",
			},
			"kern_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the kernel facility to include in the log.",
			},
			"kern_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the kernel facility to include in the log.",
			},
			"local6_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the local6 facility to include in the log.",
			},
			"local6_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the local6 facility to include in the log.",
			},
			"mail_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the mail facility to include in the log.",
			},
			"mail_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the mail facility to include in the log.",
			},
			"messages_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the user/messages facility to include in the log.",
			},
			"messages_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the user/messages facility to include in the log.",
			},
			"user_log_from": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the lowest level of messages from the user-log facility to include in the log.",
			},
			"user_log_to": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the highest level of messages from the user-log facility to include in the log.",
			},
			"remote_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of remote syslog servers to which the BIG-IP system sends log messages.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the remote syslog server entry.",
						},
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP address or hostname of the remote syslog server.",
						},
						"local_ip": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Local IP address the BIG-IP system uses when sending messages to this remote syslog server.",
						},
						"remote_port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     514,
							Description: "Port on the remote syslog server to which the BIG-IP system sends log messages.",
						},
					},
				},
			},
		},
	}
}

func getSysSyslogConfig(d *schema.ResourceData) *sysSyslogConfig {
	config := &sysSyslogConfig{
		AuthPrivFrom: d.Get("auth_priv_from").(string),
		AuthPrivTo:   d.Get("auth_priv_to").(string),
		ConsoleLog:   d.Get("console_log").(string),
		CronFrom:     d.Get("cron_from").(string),
		CronTo:       d.Get("cron_to").(string),
		DaemonFrom:   d.Get("daemon_from").(string),
		DaemonTo:     d.Get("daemon_to").(string),
		IsoDate:      d.Get("iso_date").(string),
		KernFrom:     d.Get("kern_from").(string),
		KernTo:       d.Get("kern_to").(string),
		Local6From:   d.Get("local6_from").(string),
		Local6To:     d.Get("local6_to").(string),
		MailFrom:     d.Get("mail_from").(string),
		MailTo:       d.Get("mail_to").(string),
		MessagesFrom: d.Get("messages_from").(string),
		MessagesTo:   d.Get("messages_to").(string),
		UserLogFrom:  d.Get("user_log_from").(string),
		UserLogTo:    d.Get("user_log_to").(string),
		// ClusteredHostSlot/ClusteredMessageSlot are intentionally left
		// unset here (and therefore omitted from the PATCH body, via their
		// omitempty json tags): they aren't exposed in this resource's
		// schema, so hardcoding a value would silently override whatever
		// the platform's actual chassis-clustering setting is on every
		// create/update, rather than leaving it untouched.
	}
	remoteServers := d.Get("remote_servers").([]interface{})
	config.RemoteServers = make([]sysSyslogRemoteServer, 0, len(remoteServers))
	for _, rs := range remoteServers {
		m := rs.(map[string]interface{})
		config.RemoteServers = append(config.RemoteServers, sysSyslogRemoteServer{
			Name:       m["name"].(string),
			Host:       m["host"].(string),
			LocalIP:    m["local_ip"].(string),
			RemotePort: m["remote_port"].(int),
		})
	}
	return config
}

func resourceBigipSysSyslogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Configuring System Syslog")

	if err := restPatch(client, "sys/syslog", getSysSyslogConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to configure System Syslog (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId("syslog")
	return resourceBigipSysSyslogRead(ctx, d, meta)
}

func resourceBigipSysSyslogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Updating System Syslog")

	if err := restPatch(client, "sys/syslog", getSysSyslogConfig(d)); err != nil {
		log.Printf("[ERROR] Unable to update System Syslog (%v)", err)
		return diag.FromErr(err)
	}
	return resourceBigipSysSyslogRead(ctx, d, meta)
}

func resourceBigipSysSyslogRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Reading System Syslog")

	var syslog sysSyslogConfig
	found, err := restGet(client, "sys/syslog", &syslog)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve System Syslog (%v)", err)
		return diag.FromErr(err)
	}
	if !found {
		log.Printf("[WARN] System Syslog config not found, removing from state")
		d.SetId("")
		return nil
	}

	_ = d.Set("auth_priv_from", syslog.AuthPrivFrom)
	_ = d.Set("auth_priv_to", syslog.AuthPrivTo)
	_ = d.Set("console_log", syslog.ConsoleLog)
	_ = d.Set("cron_from", syslog.CronFrom)
	_ = d.Set("cron_to", syslog.CronTo)
	_ = d.Set("daemon_from", syslog.DaemonFrom)
	_ = d.Set("daemon_to", syslog.DaemonTo)
	_ = d.Set("iso_date", syslog.IsoDate)
	_ = d.Set("kern_from", syslog.KernFrom)
	_ = d.Set("kern_to", syslog.KernTo)
	_ = d.Set("local6_from", syslog.Local6From)
	_ = d.Set("local6_to", syslog.Local6To)
	_ = d.Set("mail_from", syslog.MailFrom)
	_ = d.Set("mail_to", syslog.MailTo)
	_ = d.Set("messages_from", syslog.MessagesFrom)
	_ = d.Set("messages_to", syslog.MessagesTo)
	_ = d.Set("user_log_from", syslog.UserLogFrom)
	_ = d.Set("user_log_to", syslog.UserLogTo)

	remoteServers := make([]map[string]interface{}, 0, len(syslog.RemoteServers))
	for _, rs := range syslog.RemoteServers {
		localIP := rs.LocalIP
		if localIP == "none" {
			// BIG-IP's placeholder for "not configured"; normalize to ""
			// to match the schema's zero value and avoid a perpetual diff
			// against a config that doesn't set local_ip.
			localIP = ""
		}
		remoteServers = append(remoteServers, map[string]interface{}{
			"name":        stripPartitionPrefix(rs.Name),
			"host":        rs.Host,
			"local_ip":    localIP,
			"remote_port": rs.RemotePort,
		})
	}
	_ = d.Set("remote_servers", remoteServers)

	return nil
}

func resourceBigipSysSyslogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	log.Println("[INFO] Resetting System Syslog to defaults")

	defaults := &sysSyslogConfig{
		AuthPrivFrom: "notice",
		AuthPrivTo:   "emerg",
		// See the equivalent comment in getSysSyslogConfig: ClusteredHostSlot/
		// ClusteredMessageSlot are intentionally left unset (and therefore
		// omitted from the PATCH body) since they aren't part of this
		// resource's schema.
		ConsoleLog:    "enabled",
		CronFrom:      "warning",
		CronTo:        "emerg",
		DaemonFrom:    "notice",
		DaemonTo:      "emerg",
		IsoDate:       "disabled",
		KernFrom:      "debug",
		KernTo:        "emerg",
		Local6From:    "notice",
		Local6To:      "emerg",
		MailFrom:      "notice",
		MailTo:        "emerg",
		MessagesFrom:  "notice",
		MessagesTo:    "warning",
		UserLogFrom:   "notice",
		UserLogTo:     "emerg",
		RemoteServers: []sysSyslogRemoteServer{},
	}
	if err := restPatch(client, "sys/syslog", defaults); err != nil {
		log.Printf("[ERROR] Unable to reset System Syslog (%v)", err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
