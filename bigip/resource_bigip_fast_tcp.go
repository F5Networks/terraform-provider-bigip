/*
Copyright 2021 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceBigipFastTcpApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastTcpAppCreate,
		Read:   resourceBigipFastTcpAppRead,
		Update: resourceBigipFastTcpAppUpdate,
		Delete: resourceBigipFastTcpAppDelete,
		Exists: resourceBigipFastTcpAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCP FAST application",
			},
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the TCP FAST application tenant",
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
			"existing_snat_pool": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Name of an existing BIG-IP SNAT pool.",
				ConflictsWith: []string{"fast_create_snat_pool_address"},
			},
			"fast_create_snat_pool_address": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"existing_snat_pool"},
			},
			"exist_pool_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Select an existing BIG-IP Pool.",
				ConflictsWith: []string{"fast_create_pool_members"},
			},
			"fast_create_pool_members": {
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
				ConflictsWith: []string{"exist_pool_name"},
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
				ConflictsWith: []string{"exist_pool_name", "fast_create_monitor"},
			},
			"fast_create_monitor": {
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
					},
				},
				ConflictsWith: []string{"existing_monitor", "exist_pool_name"},
			},
			"fast_tcp_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Json payload for FAST TCP application.",
			},
		},
	}
}

func resourceBigipFastTcpAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	const templateName string = "bigip-fast-templates/tcp"
	m.Lock()
	defer m.Unlock()

	log.Printf("[INFO] Creating FAST TCP Application")

	userAgent := fmt.Sprintf("?userAgent=%s/%s", client.UserAgent, templateName)
	cfg, err := getParamsConfigMap(d)
	if err != nil {
		return nil
	}
	tenant, app, err := client.PostFastAppBigip(cfg, templateName, userAgent)
	if err != nil {
		return err
	}
	d.Set("application", app)
	d.Set("tenant", tenant)
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

	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var fastTcp bigip.FastTCPJson
	log.Printf("[INFO] Reading FastApp config")
	tenant := d.Get("tenant").(string)
	app_name := d.Id()

	log.Printf("[INFO] Reading FAST TCP Application config")
	fastJson, err := client.GetFastApp(tenant, app_name)
	log.Printf("[DEBUG] FAST json retreived from the GET call in Read function : %s", fastJson)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return nil
		}
		return err
	}
	if fastJson == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	_ = d.Set("fast_tcp_json", fastJson)
	err = json.Unmarshal([]byte(fastJson), &fastTcp)
	if err != nil {
		return err
	}
	err = setFastTcpData(d, fastTcp)
	if err != nil {
		return err
	}
	return nil
}

func resourceBigipFastTcpAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()

	name := d.Get("application").(string)
	tenant := d.Get("tenant").(string)
	cfg, err := getParamsConfigMap(d)
	log.Printf("[INFO] Updating FastApp Config :%v", cfg)
	if err != nil {
		return nil
	}
	err = client.ModifyFastAppBigip(cfg, tenant, name)

	if err != nil {
		return err
	}
	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	m.Lock()
	defer m.Unlock()
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.DeleteFastAppBigip(tenant, name)
	if err != nil {
		return err
	}
	d.SetId("")
	return resourceBigipFastTcpAppRead(d, meta)
}

func resourceBigipFastTcpAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	tenant := d.Get("tenant").(string)
	app_name := d.Get("application").(string)

	log.Printf("[INFO] Reading FAST TCP Application config")
	resp, err := client.GetFastApp(tenant, app_name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return false, nil
		}
		return false, err
	}
	if resp == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	return true, nil
}

func setFastTcpData(d *schema.ResourceData, data bigip.FastTCPJson) error {
	_ = d.Set("virtual_server.0.ip", data.VirtualAddress)
	_ = d.Set("virtual_server.0.port", data.VirtualPort)
	_ = d.Set("snat.enable", data.SnatEnable)
	_ = d.Set("snat.automap", data.SnatAutomap)
	_ = d.Set("snat.existing_snat_pool", data.SnatPoolName)
	_ = d.Set("snat.snat_addresses", data.SnatAddresses)
	_ = d.Set("pool.enable", data.PoolEnable)
	_ = d.Set("pool.existing_pool", data.PoolName)
	members := flattenFastPoolMembers(data.PoolMembers)
	_ = d.Set("pool.pool_members", members)
	_ = d.Set("load_balancing_mode", data.LoadBalancingMode)
	_ = d.Set("slow_ramp_time", data.SlowRampTime)
	_ = d.Set("monitor.enable", data.MonitorEnable)
	_ = d.Set("existing_monitor", data.TCPMonitor)
	_ = d.Set("fast_create_monitor.0.interval", data.MonitorInterval)

	return nil
}

func getParamsConfigMap(d *schema.ResourceData) (string, error) {
	// paramConfig := make(map[string]interface{})
	tcpJson := &bigip.FastTCPJson{
		Tenant:      d.Get("tenant").(string),
		Application: d.Get("application").(string),
	}

	if v, ok := d.GetOk("virtual_server"); ok {
		vL := v.([]interface{})
		for _, v := range vL {
			tcpJson.VirtualAddress = v.(map[string]interface{})["ip"].(string)
			tcpJson.VirtualPort = v.(map[string]interface{})["port"]
		}
	}

	tcpJson.SnatEnable = true
	tcpJson.SnatAutomap = true
	tcpJson.MakeSnatPool = false
	if v, ok := d.GetOk("existing_snat_pool"); ok {
		tcpJson.SnatPoolName = v.(string)
		tcpJson.SnatAutomap = false
		tcpJson.MakeSnatPool = false
	}
	if s, ok := d.GetOk("fast_create_snat_pool_address"); ok {
		tcpJson.SnatAutomap = false
		tcpJson.MakeSnatPool = true
		var snatAdd []string
		for _, addr := range s.([]interface{}) {
			snatAdd = append(snatAdd, addr.(string))
		}
		tcpJson.SnatAddresses = snatAdd
	}

	tcpJson.PoolEnable = false
	if v, ok := d.GetOk("exist_pool_name"); ok {
		tcpJson.PoolEnable = true
		tcpJson.PoolName = v.(string)
		tcpJson.MakePool = false
	}
	if p, ok := d.GetOk("fast_create_pool_members"); ok {
		tcpJson.PoolEnable = true
		tcpJson.MakePool = true
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
		tcpJson.PoolMembers = members
	}

	tcpJson.MonitorEnable = false
	tcpJson.MakeMonitor = false
	if v, ok := d.GetOk("existing_monitor"); ok {
		tcpJson.MonitorEnable = true
		tcpJson.TCPMonitor = v.(string)
	}
	if s, ok := d.GetOk("fast_create_monitor"); ok {
		tcpJson.MonitorEnable = true
		tcpJson.MakeMonitor = true
		v := s.([]interface{})[0].(map[string]interface{})
		interval := v["interval"]
		tcpJson.MonitorInterval = interval.(int)
	}
	if p, ok := d.GetOk("load_balancing_mode"); ok {
		tcpJson.LoadBalancingMode = p.(string)
	}
	if p, ok := d.GetOk("slow_ramp_time"); ok {
		tcpJson.SlowRampTime = p.(int)
	}

	data, err := json.Marshal(tcpJson)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
