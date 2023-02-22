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

func resourceBigipLtmProfileTcp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmProfileTcpCreate,
		UpdateContext: resourceBigipLtmProfileTcpUpdate,
		ReadContext:   resourceBigipLtmProfileTcpRead,
		DeleteContext: resourceBigipLtmProfileTcpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
				Description:  "Name of the TCP Profile",
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of partition",
			},
			"defaults_from": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateF5Name,
				Description:  "Use the parent tcp profile",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300; may not be 0) connection may remain idle before it becomes eligible for deletion. Value -1 (not recommended) means infinite",
			},
			"close_wait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 5) connection will remain in LAST-ACK state before exiting. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"finwait_2timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 300) connection will remain in LAST-ACK state before closing. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"finwait_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 5) connection will remain in FIN-WAIT-1 or closing state before exiting. Value -1 means indefinite, limited by maximum retransmission timeout",
			},
			"keepalive_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of seconds (default 1800) between keep-alive probes",
			},
			"congestion_control": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies the algorithm to use to share network resources among competing users to reduce congestion. The default is High Speed.",
				ValidateFunc: validation.StringInSlice([]string{"none", "high-speed", "bbr", "cdg"}, false),
			},
			"initial_congestion_windowsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the initial congestion window size for connections to this destination. Actual window size is this value multiplied by the MSS (Maximum Segment Size) for the same connection. The default is 10. Valid values range from 0 to 64",
			},
			"delayed_acks": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies, when checked (enabled), that the system can send fewer than one ACK (acknowledgment) segment per data segment received. By default, this setting is enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"nagle": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies whether the system applies Nagle's algorithm to reduce the number of short segments on the network.If you select Auto, the system determines whether to use Nagle's algorithm based on network conditions. By default, this setting is disabled.",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled", "auto"}, false),
			},
			"early_retransmit": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Enabling this setting allows TCP to assume a packet is lost after fewer than the standard number of duplicate ACKs, if there is no way to send new data and generate more duplicate ACKs",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"tailloss_probe": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Enabling this setting allows TCP to send a probe segment to trigger fast recovery instead of recovering a loss via a retransmission timeout,By default, this setting is enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"timewait_recycle": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Using this setting enabled, the system can recycle a wait-state connection immediately upon receipt of a new connection request instead of having to wait until the connection times out of the wait state. By default, this setting is enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
			"proxybuffer_high": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the proxy buffer level, in bytes, at which the receive window is closed.",
			},
			"receive_windowsize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum advertised RECEIVE window size. This value represents the maximum number of bytes to which the RECEIVE window can scale. The default is 65535 bytes",
			},
			"send_buffersize": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the SEND window size. The default is 131072 bytes",
			},
			"zerowindow_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the timeout in milliseconds for terminating a connection with an effective zero length TCP transmit window",
			},
			"deferred_accept": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "If enabled, ADC will defer allocating resources to a connection until some payload data has arrived from the client (default false). This may help minimize the impact of certain DoS attacks but adds undesirable latency under normal conditions. Note: ‘deferredAccept’ is incompatible with server-speaks-first application protocols,Default : disabled",
			},
			"fast_open": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "If enabled (default), the system can use the TCP Fast Open protocol extension to reduce latency by sending payload data with initial SYN",
			},
			"verified_accept": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Specifies, when checked (enabled), that the system can actually communicate with the server before establishing a client connection. To determine this, the system sends the server a SYN packet before responding to the client's SYN with a SYN-ACK. When unchecked, the system accepts the client connection before selecting a server to talk to. By default, this setting is disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
			},
		},
	}
}

func resourceBigipLtmProfileTcpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	tcpConfig := &bigip.Tcp{
		Name: name,
	}
	tcpProfileConfig := getTCPProfileConfig(d, tcpConfig)
	log.Println("[INFO] Creating TCP profile")
	err := client.CreateTcp(tcpProfileConfig)
	if err != nil {
		log.Printf("[ERROR] Unable to Create tcp Profile  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
	return resourceBigipLtmProfileTcpRead(ctx, d, meta)
}

func resourceBigipLtmProfileTcpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating TCP Profile " + name)
	tcpConfig := &bigip.Tcp{
		Name: name,
	}
	tcpProfileConfig := getTCPProfileConfig(d, tcpConfig)
	err := client.ModifyTcp(name, tcpProfileConfig)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error create profile tcp (%s): %s", name, err))
	}
	return resourceBigipLtmProfileTcpRead(ctx, d, meta)
}

func resourceBigipLtmProfileTcpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Reading TCP Profile  " + name)
	obj, err := client.GetTcp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve tcp Profile  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Reading TCP Object:%+v ", obj)
	if obj == nil {
		log.Printf("[WARN] tcp  Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("name", name)
	if _, ok := d.GetOk("defaults_from"); ok {
		_ = d.Set("defaults_from", obj.DefaultsFrom)
	}
	if _, ok := d.GetOk("idle_timeout"); ok {
		_ = d.Set("idle_timeout", obj.IdleTimeout)
	}
	if _, ok := d.GetOk("close_wait_timeout"); ok {
		_ = d.Set("close_wait_timeout", obj.CloseWaitTimeout)
	}
	if _, ok := d.GetOk("finwait_2timeout"); ok {
		_ = d.Set("finwait_2timeout", obj.FinWait_2Timeout)
	}
	if _, ok := d.GetOk("finwait_timeout"); ok {
		_ = d.Set("finwait_timeout", obj.FinWaitTimeout)
	}
	if _, ok := d.GetOk("congestion_control"); ok {
		_ = d.Set("congestion_control", obj.CongestionControl)
	}
	if _, ok := d.GetOk("delayed_acks"); ok {
		_ = d.Set("delayed_acks", obj.DelayedAcks)
	}
	if _, ok := d.GetOk("nagle"); ok {
		_ = d.Set("nagle", obj.Nagle)
	}
	if _, ok := d.GetOk("early_retransmit"); ok {
		_ = d.Set("early_retransmit", obj.EarlyRetransmit)
	}
	if _, ok := d.GetOk("tailloss_probe"); ok {
		_ = d.Set("tailloss_probe", obj.TailLossProbe)
	}
	if _, ok := d.GetOk("initial_congestion_windowsize"); ok {
		_ = d.Set("initial_congestion_windowsize", obj.InitCwnd)
	}
	if _, ok := d.GetOk("zerowindow_timeout"); ok {
		_ = d.Set("zerowindow_timeout", obj.ZeroWindowTimeout)
	}
	if _, ok := d.GetOk("send_buffersize"); ok {
		_ = d.Set("send_buffersize", obj.SendBufferSize)
	}
	if _, ok := d.GetOk("receive_windowsize"); ok {
		_ = d.Set("receive_windowsize", obj.ReceiveWindowSize)
	}
	if _, ok := d.GetOk("proxybuffer_high"); ok {
		_ = d.Set("proxybuffer_high", obj.ProxyBufferHigh)
	}
	if _, ok := d.GetOk("timewait_recycle"); ok {
		_ = d.Set("timewait_recycle", obj.TimeWaitRecycle)
	}
	if _, ok := d.GetOk("verified_accept"); ok {
		_ = d.Set("verified_accept", obj.VerifiedAccept)
	}
	if _, ok := d.GetOk("keepalive_interval"); ok {
		_ = d.Set("keepalive_interval", obj.KeepAliveInterval)
	}
	if _, ok := d.GetOk("deferred_accept"); ok {
		_ = d.Set("deferred_accept", obj.DeferredAccept)
	}
	if _, ok := d.GetOk("fast_open"); ok {
		_ = d.Set("fast_open", obj.FastOpen)
	}
	return nil
}

func resourceBigipLtmProfileTcpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting Tcp Profile " + name)

	err := client.DeleteTcp(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete tcp Profile (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getTCPProfileConfig(d *schema.ResourceData, config *bigip.Tcp) *bigip.Tcp {
	config.Partition = d.Get("partition").(string)
	config.DefaultsFrom = d.Get("defaults_from").(string)
	config.IdleTimeout = d.Get("idle_timeout").(int)
	config.CloseWaitTimeout = d.Get("close_wait_timeout").(int)
	config.FinWait_2Timeout = d.Get("finwait_2timeout").(int)
	config.FinWaitTimeout = d.Get("finwait_timeout").(int)
	config.SendBufferSize = d.Get("send_buffersize").(int)
	config.ReceiveWindowSize = d.Get("receive_windowsize").(int)
	config.ProxyBufferHigh = d.Get("proxybuffer_high").(int)
	config.ZeroWindowTimeout = d.Get("zerowindow_timeout").(int)
	config.KeepAliveInterval = d.Get("keepalive_interval").(int)
	config.CongestionControl = d.Get("congestion_control").(string)
	config.InitCwnd = d.Get("initial_congestion_windowsize").(int)
	config.DelayedAcks = d.Get("delayed_acks").(string)
	config.Nagle = d.Get("nagle").(string)
	config.EarlyRetransmit = d.Get("early_retransmit").(string)
	config.TailLossProbe = d.Get("tailloss_probe").(string)
	config.TimeWaitRecycle = d.Get("timewait_recycle").(string)
	config.VerifiedAccept = d.Get("verified_accept").(string)
	config.DeferredAccept = d.Get("deferred_accept").(string)
	config.FastOpen = d.Get("fast_open").(string)
	return config
}
