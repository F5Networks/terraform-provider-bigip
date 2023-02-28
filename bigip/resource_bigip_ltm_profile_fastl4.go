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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipLtmProfileFastl4() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipProfileLtmFastl4Create,
		UpdateContext: resourceBigipLtmProfileFastl4Update,
		ReadContext:   resourceBigipLtmProfileFastl4Read,
		DeleteContext: resourceBigipLtmProfileFastl4Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the Fastl4 Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent Fastl4 profile",
			},
			"late_binding": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "Enables intelligent selection of a back-end server or pool, using an iRule to make the selection. The default is disabled",
			},
			"client_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds allowed for a client to transmit enough data to select a server when you have late binding enabled. Value -1 means indefinite (not recommended)",
			},
			"explicitflow_migration": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "Specifies whether a qualified late-binding connection requires an explicit iRule command to migrate down to ePVA hardware. The default is disabled",
			},
			"hardware_syncookie": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "Use the parent Fastl4 profile",
			},
			"idle_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300; may not be 0) connection may remain idle before it becomes eligible for deletion. Value -1 (not recommended) means infinite",
			},
			"iptos_toclient": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"iptos_toserver": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Use the parent Fastl4 profile",
			},
			"tcp_handshake_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the acceptable duration for a TCP handshake, that is, the maximum idle time between a client synchronization (SYN) and a client acknowledgment (ACK).The default is 5 seconds",
			},
			"keepalive_interval": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the keep-alive probe interval, in seconds. The default is Disabled",
			},
			"loose_initiation": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies, when checked (enabled), that the system initializes a connection when it receives any TCP packet, rather that requiring a SYN packet for connection initiation. The default is disabled. We recommend that if you enable the Loose Initiation option, you also enable the Loose Close option.",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"loose_close": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Specifies, when checked (enabled), that the system closes a loosely-initiated connection when the system receives the first FIN packet from either the client or the server. The default is disabled.",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"receive_windowsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the amount of data the BIG-IP system can accept without acknowledging the server. The default is 0 (zero)",
			},
		},
	}
}

func resourceBigipProfileLtmFastl4Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)

	log.Println("[INFO] Creating Fastl4 profile")
	configFastl4 := &bigip.Fastl4{
		Name: name,
	}
	if d.Get("explicitflow_migration").(string) == "enabled" && d.Get("late_binding").(string) != "enabled" {
		return diag.FromErr(fmt.Errorf("explicitflow_migration can be enabled only if late_binding set to enabled"))
	}
	fastL4ProfileConfig := getFastL4ProfileConfig(d, configFastl4)

	err := client.CreateFastl4(fastL4ProfileConfig)

	if err != nil {
		log.Printf("[ERROR] Unable to Create FastL4  (%s) (%v) ", name, err)
		return diag.FromErr(fmt.Errorf("Error retrieving profile fastl4 (%s): %s ", name, err))
	}

	d.SetId(name)
	return resourceBigipLtmProfileFastl4Read(ctx, d, meta)
}

func resourceBigipLtmProfileFastl4Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	configFastl4 := &bigip.Fastl4{
		Name: name,
	}
	if d.Get("explicitflow_migration").(string) == "enabled" && d.Get("late_binding").(string) != "enabled" {
		return diag.FromErr(fmt.Errorf("explicitflow_migration can be enabled only if late_binding set to enabled"))
	}
	log.Println("[INFO] Updating Fastl4 profile")
	fastL4ProfileConfig := getFastL4ProfileConfig(d, configFastl4)

	err := client.ModifyFastl4(name, fastL4ProfileConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Modify FastL4  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	return resourceBigipLtmProfileFastl4Read(ctx, d, meta)
}

func resourceBigipLtmProfileFastl4Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	obj, err := client.GetFastl4(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve FastL4  (%s) (%v) ", name, err)
		return diag.FromErr(err)
	}
	if obj == nil {
		log.Printf("[WARN] Fastl4 profile  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	_ = d.Set("defaults_from", obj.DefaultsFrom)

	if _, ok := d.GetOk("client_timeout"); ok {
		_ = d.Set("client_timeout", obj.ClientTimeout)
	}
	if _, ok := d.GetOk("explicitflow_migration"); ok {
		_ = d.Set("explicitflow_migration", obj.ExplicitFlowMigration)
	}
	if _, ok := d.GetOk("iptos_toclient"); ok {
		_ = d.Set("iptos_toclient", obj.IpTosToClient)
	}
	if _, ok := d.GetOk("iptos_toserver"); ok {
		_ = d.Set("iptos_toserver", obj.IpTosToServer)
	}
	if _, ok := d.GetOk("hardware_syncookie"); ok {
		_ = d.Set("hardware_syncookie", obj.HardwareSynCookie)
	}
	if _, ok := d.GetOk("idle_timeout"); ok {
		_ = d.Set("idle_timeout", obj.IdleTimeout)
	}
	if _, ok := d.GetOk("keepalive_interval"); ok {
		_ = d.Set("keepalive_interval", obj.KeepAliveInterval)
	}
	if _, ok := d.GetOk("tcp_handshake_timeout"); ok {
		_ = d.Set("tcp_handshake_timeout", obj.TCPHandshakeTimeout)
	}
	if _, ok := d.GetOk("loose_initiation"); ok {
		_ = d.Set("loose_initiation", obj.LooseInitialization)
	}
	if _, ok := d.GetOk("loose_close"); ok {
		_ = d.Set("loose_close", obj.LooseClose)
	}
	if _, ok := d.GetOk("late_binding"); ok {
		_ = d.Set("late_binding", obj.LateBinding)
	}
	if _, ok := d.GetOk("receive_windowsize"); ok {
		_ = d.Set("receive_windowsize", obj.ReceiveWindowSize)
	}
	return nil
}

func resourceBigipLtmProfileFastl4Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Fastl4 Profile " + name)

	err := client.DeleteFastl4(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve node (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getFastL4ProfileConfig(d *schema.ResourceData, config *bigip.Fastl4) *bigip.Fastl4 {
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.ClientTimeout = d.Get("client_timeout").(int)
	config.LateBinding = d.Get("late_binding").(string)
	config.ExplicitFlowMigration = d.Get("explicitflow_migration").(string)
	config.HardwareSynCookie = d.Get("hardware_syncookie").(string)
	config.IdleTimeout = d.Get("idle_timeout").(string)
	config.IpTosToClient = d.Get("iptos_toclient").(string)
	config.IpTosToServer = d.Get("iptos_toserver").(string)
	config.KeepAliveInterval = d.Get("keepalive_interval").(string)
	config.TCPHandshakeTimeout = d.Get("tcp_handshake_timeout").(string)
	config.LooseInitialization = d.Get("loose_initiation").(string)
	config.LooseClose = d.Get("loose_close").(string)
	config.ReceiveWindowSize = d.Get("receive_windowsize").(int)
	return config
}
