/*
Copyright 2022 F5 Networks Inc.
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

var fastTmpl = "bigip-fast-templates/http"

func resourceBigipHttpFastApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastHttpAppCreate,
		Read:   resourceBigipFastHttpAppRead,
		Update: resourceBigipFastHttpAppUpdate,
		Delete: resourceBigipFastHttpAppDelete,
		Exists: resourceBigipFastHttpAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of FAST HTTP application tenant.",
			},
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of FAST HTTP application.",
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
							Description: "foo",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "foo",
						},
					},
				},
			},
			"existing_snat_pool": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "name of an existing BIG-IP SNAT pool",
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
				Description:   "Select an existing BIG-IP Pool",
				ConflictsWith: []string{"fast_create_pool_members"},
			},
			"fast_create_pool_members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addresses": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "foo",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     80,
							Description: "foo",
						},
						"connection_limit": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "foo",
						},
						"priority_group": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "foo",
						},
						"share_nodes": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "foo",
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
				Description: "none",
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
				Description: "foo",
				MaxItems:    1,

				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "foo",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "foo",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "foo",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Set the time between health checks, in seconds.",
						},
						"send_string": {
							Type:     schema.TypeString,
							Optional: true,
							//Default:"GET / HTTP/1.1\r\nHost: example.com\r\nConnection: Close\r\n\r\n",
							Description: "foo",
						},
						"response": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "foo",
						},
					},
				},
				ConflictsWith: []string{"existing_monitor", "exist_pool_name"},
			},
		},
	}
}

func resourceBigipFastHttpAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastJson, err := getFastHttpConfig(d)
	if err != nil {
		return err
	}
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Creating HTTP FastApp config")
	userAgent := fmt.Sprintf("?userAgent=%s/%s", client.UserAgent, fastTmpl)
	tenant, app, err := client.PostFastAppBigip(fastJson, fastTmpl, userAgent)
	if err != nil {
		return err
	}
	_ = d.Set("tenant", tenant)
	_ = d.Set("application", app)
	log.Printf("[DEBUG] ID for resource :%+v", app)
	d.SetId(d.Get("application").(string))
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
			"tenant":            tenant,
			"application":       app,
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_fast_http_app", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipFastHttpAppRead(d, meta)
}

func resourceBigipFastHttpAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	var fastHttp bigip.FastHttpJson
	log.Printf("[INFO] Reading FastApp config")
	name := d.Id()
	tenant := d.Get("tenant").(string)
	log.Printf("[DEBUG] FAST HTTP application get call : %s", name)
	fastJson, err := client.GetFastApp(tenant, name)
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
	_ = d.Set("fast_json", fastJson)
	err = json.Unmarshal([]byte(fastJson), &fastHttp)
	if err != nil {
		return err
	}
	err = setFastHttpData(d, fastHttp)
	if err != nil {
		return err
	}
	return nil
}

func resourceBigipFastHttpAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastJson, e := getFastHttpConfig(d)
	if e != nil {
		return e
	}
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Updating FastApp Config :%s", fastJson)
	name := d.Id()
	tenant := d.Get("tenant").(string)
	err := client.ModifyFastAppBigip(fastJson, tenant, name)
	if err != nil {
		return err
	}
	return resourceBigipFastAppRead(d, meta)
}

func resourceBigipFastHttpAppDelete(d *schema.ResourceData, meta interface{}) error {
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
	return nil
}

func resourceBigipFastHttpAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)
	log.Printf("[INFO] Checking if FastApp config exists in BIGIP")
	name := d.Id()
	tenant := d.Get("tenant").(string)
	fastJson, err := client.GetFastApp(tenant, name)
	if err != nil {
		log.Printf("[ERROR] Unable to retrieve json ")
		if err.Error() == "unexpected end of JSON input" {
			log.Printf("[ERROR] %v", err)
			d.SetId("")
			return false, nil
		}
		return false, err
	}
	log.Printf("[INFO] FAST response Body:%+v", fastJson)
	if fastJson == "" {
		log.Printf("[WARN] Json (%s) not found, removing from state", d.Id())
		return false, nil
	}
	return true, nil
}

func setFastHttpData(d *schema.ResourceData, data bigip.FastHttpJson) error {

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
	_ = d.Set("existing_monitor", data.HTTPMonitor)
	_ = d.Set("fast_create_monitor.0.monitor_auth", data.MonitorAuth)
	_ = d.Set("fast_create_monitor.0.username", data.MonitorUsername)
	_ = d.Set("fast_create_monitor.0.password", data.MonitorPassword)
	_ = d.Set("fast_create_monitor.0.interval", data.MonitorInterval)
	_ = d.Set("fast_create_monitor.0.send_string", data.MonitorSendString)
	_ = d.Set("fast_create_monitor.0.response", data.MonitorResponse)
	return nil
}

func flattenFastVirtualServers(data bigip.FastHttpJson) map[string]interface{} {
	tfMap := map[string]interface{}{}
	tfMap["ip"] = data.VirtualAddress
	tfMap["port"] = data.VirtualPort
	return tfMap
}

func flattenFastMonitor(data bigip.FastHttpJson) map[string]interface{} {
	tfMap := map[string]interface{}{}
	if data.MonitorAuth {
		tfMap["monitor_auth"] = data.MonitorAuth
	}
	tfMap["username"] = data.MonitorUsername
	tfMap["password"] = data.MonitorPassword
	if data.MonitorInterval > 0 {
		tfMap["interval"] = data.MonitorInterval
	}
	tfMap["send_string"] = data.MonitorSendString
	tfMap["response"] = data.MonitorResponse
	return tfMap
}

func flattenFastPoolMembers(members []bigip.FastHttpPool) []interface{} {
	att := make([]interface{}, len(members))
	for i, v := range members {
		obj := make(map[string]interface{})
		if len(v.ServerAddresses) > 0 {
			obj["adresses"] = v.ServerAddresses
		}
		obj["port"] = v.ServicePort
		if v.ConnectionLimit > 0 {
			obj["connection_limit"] = v.ConnectionLimit
		}
		obj["priority_group"] = v.PriorityGroup
		if v.ShareNodes {
			obj["share_nodes"] = v.ShareNodes
		}
		att[i] = obj
	}
	return att
}

func getFastHttpConfig(d *schema.ResourceData) (string, error) {
	httpJson := &bigip.FastHttpJson{
		Tenant:      d.Get("tenant").(string),
		Application: d.Get("application").(string),
	}
	httpJson.TlsServerEnable = false
	httpJson.TlsServerProfileCreate = false
	if v, ok := d.GetOk("virtual_server"); ok {
		vL := v.([]interface{})
		for _, v := range vL {
			httpJson.VirtualAddress = v.(map[string]interface{})["ip"].(string)
			httpJson.VirtualPort = v.(map[string]interface{})["port"]
		}
	}

	httpJson.PoolEnable = false
	if v, ok := d.GetOk("exist_pool_name"); ok {
		httpJson.PoolEnable = true
		httpJson.PoolName = v.(string)
		httpJson.MakePool = false
	}
	if p, ok := d.GetOk("fast_create_pool_members"); ok {
		httpJson.PoolEnable = true
		httpJson.MakePool = true
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
		httpJson.PoolMembers = members
	}
	httpJson.SnatEnable = true
	httpJson.SnatAutomap = true
	httpJson.MakeSnatPool = false
	if v, ok := d.GetOk("existing_snat_pool"); ok {
		httpJson.SnatPoolName = v.(string)
		httpJson.SnatEnable = true
		httpJson.SnatAutomap = false
		httpJson.MakeSnatPool = false
	}
	if s, ok := d.GetOk("fast_create_snat_pool_address"); ok {
		httpJson.SnatEnable = true
		httpJson.SnatAutomap = false
		httpJson.MakeSnatPool = true
		var snatAdd []string
		for _, addr := range s.([]interface{}) {
			snatAdd = append(snatAdd, addr.(string))
		}
		httpJson.SnatAddresses = snatAdd
	}
	httpJson.MonitorEnable = false
	httpJson.MakeMonitor = false
	if v, ok := d.GetOk("existing_monitor"); ok {
		httpJson.HTTPMonitor = v.(string)
		httpJson.MonitorEnable = true
		httpJson.MakeMonitor = false
	}
	if s, ok := d.GetOk("fast_create_monitor"); ok {
		log.Printf("[DEBUG] fast_create_monitor:%+v", s)
		httpJson.MonitorEnable = true
		httpJson.MakeMonitor = true
		sL := s.([]interface{})
		for _, v := range sL {
			mon := v.(map[string]interface{})
			if auth, ok := mon["monitor_auth"].(bool); ok {
				if auth {
					if u, ok := mon["username"].(string); ok {
						httpJson.MonitorUsername = u
					} else {
						return "", fmt.Errorf("the 'username' must be specified if 'monitor_auth' is %t", auth)
					}
					if p, ok := mon["password"].(string); ok {
						httpJson.MonitorPassword = p
					} else {
						return "", fmt.Errorf("the 'password' must be specified if 'monitor_auth' is %t", auth)
					}
				}
				httpJson.MonitorAuth = auth
			}
			if e, ok := mon["interval"].(int); ok {
				httpJson.MonitorInterval = e
			}
			if e, ok := mon["send_string"].(string); ok {
				httpJson.MonitorSendString = e
			}
			if e, ok := mon["response"].(string); ok {
				httpJson.MonitorResponse = e
			}
		}
	}
	if p, ok := d.GetOk("load_balancing_mode"); ok {
		httpJson.LoadBalancingMode = p.(string)
	}
	if p, ok := d.GetOk("slow_ramp_time"); ok {
		httpJson.SlowRampTime = p.(int)
	}
	data, err := json.Marshal(httpJson)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
