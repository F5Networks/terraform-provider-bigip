/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBigipFastUdpApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipFastUdpAppCreate,
		ReadContext:   resourceBigipFastUdpAppRead,
		UpdateContext: resourceBigipFastUdpAppUpdate,
		DeleteContext: resourceBigipFastUdpAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the UDP FAST application",
			},
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the UDP FAST application tenant",
			},
			"virtual_server": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "This IP address, combined with the port you specify below, becomes the BIG-IP virtual server address and port, which clients use to access the application.",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Port for the virtual server.",
						},
					},
				},
			},
			"enable_fastl4": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables use of FastL4 profiles on the virtual server.",
			},
			"existing_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of an existing BIG-IP FastL4 or UDP profile.",
			},
			"existing_snat_pool": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Name of an existing BIG-IP SNAT pool.",
				ConflictsWith: []string{"snat_pool_address"},
			},
			"snat_pool_address": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"existing_snat_pool"},
			},
			"persistence_profile": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Name of an existing BIG-IP persistence profile to be used.",
				ConflictsWith: []string{"persistence_type"},
			},
			"persistence_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of persistence profile to be created.",
				ValidateFunc: validation.StringInSlice([]string{
					"destination-address",
					"source-address"}, false),
				ConflictsWith: []string{"persistence_profile"},
			},
			"fallback_persistence": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of fallback persistence record to be created for each new client connection.",
				ValidateFunc: validation.StringInSlice([]string{
					"destination-address",
					"source-address"}, false),
			},
			"existing_pool": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Select an existing BIG-IP Pool.",
				ConflictsWith: []string{"pool_members"},
			},
			"pool_members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  80,
						},
						"connection_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"priority_group": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"share_nodes": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				ConflictsWith: []string{"existing_pool"},
			},
			"load_balancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "none",
				ValidateFunc: validation.StringInSlice([]string{
					"dynamic-ratio-member",
					"dynamic-ratio-node",
					"fastest-app-response",
					"fastest-node",
					"least-connections-member",
					"least-connections-node",
					"least-sessions",
					"observed-member",
					"observed-node",
					"predictive-member",
					"predictive-node",
					"ratio-least-connections-member",
					"ratio-least-connections-node",
					"ratio-member",
					"ratio-node",
					"ratio-session",
					"round-robin",
					"weighted-least-connections-member",
					"weighted-least-connections-node"}, false),
			},
			"slow_ramp_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Slow ramp temporarily throttles the number of connections to a new pool member.",
			},
			"existing_monitor": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Select an existing BIG-IP HTTPS pool monitor. Monitors are used to determine the health of the application on each server",
				ConflictsWith: []string{"existing_pool", "monitor"},
			},
			"monitor": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Use a FAST generated pool monitor.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Set the time between health checks, in seconds.",
						},
						"send_string": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Optional data to be sent during each health check.",
						},
						"expected_response": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The presence of this optional string is required in the response, if specified to confirms availability.",
						},
					},
				},
				ConflictsWith: []string{"existing_monitor", "existing_pool"},
			},
			"irules": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Irules to attach to Virtual Server.",
			},
			"vlans_allowed": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "Names of existing VLANs to allow.",
				ConflictsWith: []string{"vlans_rejected"},
			},
			"vlans_rejected": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				Description:   "Names of existing VLANs to reject.",
				ConflictsWith: []string{"vlans_allowed"},
			},
			"security_log_profiles": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Existing security log profiles to enable.",
			},
			"fast_udp_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Json payload for FAST UDP application.",
			},
		},
	}
}

func resourceBigipFastUdpAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	const templateName string = "bigip-fast-templates/udp"
	m.Lock()
	defer m.Unlock()

	log.Printf("[INFO] Creating FAST UDP Application")

	userAgent := fmt.Sprintf("?userAgent=%s/%s", client.UserAgent, templateName)
	cfg, err := getParamsConfigMapUdp(d)
	if err != nil {
		return nil
	}
	tenant, app, err := client.PostFastAppBigip(cfg, templateName, userAgent)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("application", app)
	_ = d.Set("tenant", tenant)
	log.Printf("[DEBUG] ID for resource :%+v", app)
	d.SetId(app)

	if !client.Teem {
		id := uuid.New()
		uniqueID := id.String()
		assetInfo := f5teem.AssetInfo{
			Name:    "Terraform-provider-bigip",
			Version: client.UserAgent,
			Id:      uniqueID,
		}
		apiKey := os.Getenv("TEEM_API_KEY")
		teemDevice := f5teem.AnonymousClient(assetInfo, apiKey)
		f := map[string]interface{}{
			"Terraform Version": client.UserAgent,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_fast_application", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}

	return resourceBigipFastUdpAppRead(ctx, d, meta)
}

func resourceBigipFastUdpAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	var fastUdp bigip.FastUDPJson
	log.Printf("[INFO] Reading FastApp config")
	tenant := d.Get("tenant").(string)
	appName := d.Id()

	log.Printf("[INFO] Reading FAST UDP Application config")
	fastJson, err := client.GetFastApp(tenant, appName)
	log.Printf("[DEBUG] FAST json retreived from the GET call in Read function : %s", fastJson)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if fastJson == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("fast_json", fastJson)
	err = json.Unmarshal([]byte(fastJson), &fastUdp)
	if err != nil {
		return diag.FromErr(err)
	}
	err = setFastUdpData(d, fastUdp)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceBigipFastUdpAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()

	name := d.Get("application").(string)
	tenant := d.Get("tenant").(string)
	cfg, err := getParamsConfigMapUdp(d)
	log.Printf("[INFO] Updating FastApp Config :%v", cfg)
	if err != nil {
		return nil
	}
	err = client.ModifyFastAppBigip(cfg, tenant, name)

	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipFastUdpAppRead(ctx, d, meta)
}

func resourceBigipFastUdpAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.DeleteFastAppBigip(tenant, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return resourceBigipFastUdpAppRead(ctx, d, meta)
}

func setFastUdpData(d *schema.ResourceData, data bigip.FastUDPJson) error {
	vsdata := make(map[string]interface{})
	vsdata["ip"] = data.VirtualAddress
	vsdata["port"] = data.VirtualPort
	_ = d.Set("virtual_server", []interface{}{vsdata})
	_ = d.Set("enable_fastl4", data.Fastl4Enable)
	if data.Fastl4Enable {
		_ = d.Set("existing_profile", data.Fastl4ProfileName)
		_ = d.Set("persistence_profile", data.Fastl4PersistenceProfile)
		_ = d.Set("persistence_type", data.Fastl4PersistenceType)
	} else {
		_ = d.Set("existing_profile", data.UdpProfileName)
		_ = d.Set("persistence_profile", data.UdpPersistenceProfile)
		_ = d.Set("persistence_type", data.UdpPersistenceType)
	}
	_ = d.Set("existing_snat_pool", data.SnatPoolName)
	_ = d.Set("snat_pool_address", data.SnatAddresses)
	_ = d.Set("existing_pool", data.PoolName)
	members := flattenFastPoolMembers(data.PoolMembers)
	_ = d.Set("pool_members", members)
	_ = d.Set("slow_ramp_time", data.SlowRampTime)
	_ = d.Set("existing_monitor", data.UdpMonitor)
	monitorData := make(map[string]interface{})
	monitorData["send_string"] = data.MonitorSendString
	monitorData["expected_response"] = data.MonitorExpectedResponse
	monitorData["interval"] = data.MonitorInterval
	if _, ok := d.GetOk("monitor"); ok {
		// _ = d.Set("monitor", []interface{}{monitorData})
		if err := d.Set("monitor", []interface{}{monitorData}); err != nil {
			return fmt.Errorf("error setting monitor: %w", err)
		}
	}
	_ = d.Set("irules", data.IruleNames)
	_ = d.Set("fallback_persistence", data.FallbackPersistenceType)
	if data.VlansAllow && data.VlansEnable {
		_ = d.Set("vlans_allowed", data.Vlans)
	}
	if !data.VlansAllow && data.VlansEnable {
		_ = d.Set("vlans_rejected", data.Vlans)
	}
	_ = d.Set("security_log_profiles", data.LogProfileNames)
	_ = d.Set("load_balancing_mode", data.LoadBalancingMode)
	return nil
}

func getParamsConfigMapUdp(d *schema.ResourceData) (string, error) {
	udpJson := &bigip.FastUDPJson{
		Tenant:      d.Get("tenant").(string),
		Application: d.Get("application").(string),
	}

	if v, ok := d.GetOk("virtual_server"); ok {
		vL := v.([]interface{})
		for _, v := range vL {
			udpJson.VirtualAddress = v.(map[string]interface{})["ip"].(string)
			udpJson.VirtualPort = v.(map[string]interface{})["port"]
		}
	}

	udpJson.Fastl4Enable = false
	udpJson.MakeFastl4Profile = false
	if v, ok := d.GetOk("enable_fastl4"); ok {
		if v.(bool) {
			udpJson.Fastl4Enable = v.(bool)
			if v2, ok2 := d.GetOk("existing_profile"); ok2 {
				udpJson.Fastl4ProfileName = v2.(string)
			} else {
				udpJson.MakeFastl4Profile = true
			}
		}
	}
	if v, ok := d.GetOk("existing_profile"); ok {
		udpJson.UdpProfileName = v.(string)
	}

	udpJson.SnatEnable = true
	udpJson.SnatAutomap = true
	udpJson.MakeSnatPool = false
	if v, ok := d.GetOk("existing_snat_pool"); ok {
		udpJson.SnatPoolName = v.(string)
		udpJson.SnatAutomap = false
		udpJson.MakeSnatPool = false
	}
	if s, ok := d.GetOk("snat_pool_address"); ok {
		udpJson.SnatAutomap = false
		udpJson.MakeSnatPool = true
		var snatAdd []string
		for _, addr := range s.([]interface{}) {
			snatAdd = append(snatAdd, addr.(string))
		}
		udpJson.SnatAddresses = snatAdd
	}

	udpJson.EnablePersistence = false
	udpJson.UseExistingPersistence = false
	udpJson.EnableFallbackPersistence = false
	if v, ok := d.GetOk("persistence_profile"); ok {
		udpJson.EnablePersistence = true
		udpJson.UseExistingPersistence = true
		if udpJson.Fastl4Enable {
			udpJson.Fastl4PersistenceProfile = v.(string)
		} else {
			udpJson.UdpPersistenceProfile = v.(string)
		}
	}
	if v, ok := d.GetOk("persistence_type"); ok {
		udpJson.EnablePersistence = true
		if udpJson.Fastl4Enable {
			udpJson.Fastl4PersistenceType = v.(string)
		} else {
			udpJson.UdpPersistenceType = v.(string)
		}
	}
	if v, ok := d.GetOk("fallback_persistence"); ok {
		udpJson.EnableFallbackPersistence = true
		udpJson.FallbackPersistenceType = v.(string)
	}

	udpJson.PoolEnable = false
	if v, ok := d.GetOk("existing_pool"); ok {
		udpJson.PoolEnable = true
		udpJson.PoolName = v.(string)
		udpJson.MakePool = false
	}
	if p, ok := d.GetOk("pool_members"); ok {
		udpJson.PoolEnable = true
		udpJson.MakePool = true
		log.Printf("[DEBUG] Adding Pool Members:%+v", p)
		var members []bigip.FastHttpPool
		for _, r := range p.(*schema.Set).List() {
			log.Printf("[DEBUG] Pool Members:%+v and Type :%T", r, r)
			memberConfig := bigip.FastHttpPool{}
			var serAdd []string
			for _, addr := range r.(map[string]interface{})["addresses"].([]interface{}) {
				serAdd = append(serAdd, addr.(string))
			}
			memberConfig.ServerAddresses = serAdd
			memberConfig.ServicePort = r.(map[string]interface{})["port"].(int)
			if s, ok := r.(map[string]interface{})["connection_limit"].(int); ok {
				memberConfig.ConnectionLimit = s
			}
			if s, ok := r.(map[string]interface{})["priority_group"].(int); ok {
				memberConfig.PriorityGroup = s
			}
			if s, ok := r.(map[string]interface{})["share_nodes"].(bool); ok {
				memberConfig.ShareNodes = s
			}
			members = append(members, memberConfig)
		}
		udpJson.PoolMembers = members
	}

	udpJson.MonitorEnable = false
	udpJson.MakeMonitor = false
	if v, ok := d.GetOk("existing_monitor"); ok {
		udpJson.MonitorEnable = true
		udpJson.UdpMonitor = v.(string)
	}
	if s, ok := d.GetOk("monitor"); ok {
		udpJson.MonitorEnable = true
		udpJson.MakeMonitor = true
		v := s.([]interface{})[0].(map[string]interface{})
		interval := v["interval"]
		sendString := v["send_string"]
		response := v["expected_response"]
		udpJson.MonitorInterval = interval.(int)
		udpJson.MonitorSendString = sendString.(string)
		udpJson.MonitorExpectedResponse = response.(string)
	}
	if v, ok := d.GetOk("load_balancing_mode"); ok {
		udpJson.LoadBalancingMode = v.(string)
	}
	if v, ok := d.GetOk("slow_ramp_time"); ok {
		udpJson.SlowRampTime = v.(int)
	}

	if s, ok := d.GetOk("irules"); ok {
		var irules []string
		for _, rule := range s.([]interface{}) {
			irules = append(irules, rule.(string))
		}
		udpJson.IruleNames = irules
	}

	udpJson.VlansEnable = false
	udpJson.VlansAllow = false
	if s, ok := d.GetOk("vlans_allowed"); ok {
		udpJson.VlansEnable = true
		udpJson.VlansAllow = true
		var vlans []string
		for _, vlan := range s.([]interface{}) {
			vlans = append(vlans, vlan.(string))
		}
		udpJson.Vlans = vlans
	}
	if s, ok := d.GetOk("vlans_rejected"); ok {
		udpJson.VlansEnable = true
		var vlans []string
		for _, vlan := range s.([]interface{}) {
			vlans = append(vlans, vlan.(string))
		}
		udpJson.Vlans = vlans
	}

	udpJson.EnableAsmLogging = false
	if s, ok := d.GetOk("security_log_profiles"); ok {
		udpJson.EnableAsmLogging = true
		var asmlogs []string
		for _, asmlog := range s.([]interface{}) {
			asmlogs = append(asmlogs, asmlog.(string))
		}
		udpJson.LogProfileNames = asmlogs
	}

	data, err := json.Marshal(udpJson)
	if err != nil {
		return "", err
	}
	_ = d.Set("fast_udp_json", string(data))
	return string(data), nil
}
