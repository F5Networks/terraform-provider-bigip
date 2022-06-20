/*
Copyright 2022 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"os"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

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
				Description: "Name of FAST HTTP application tenant.",
			},
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of FAST HTTP application.",
			},
			"virtual_server": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Computed:    true,
							Description: "foo",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Computed:    true,
							Description: "foo",
						},
					},
				},
			},
			"snat": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "foo",
						},
						"automap": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"snat_addresses", "existing_snat_pool"},
						},
						"existing_snat_pool": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"snat_addresses", "automap"},
						},
						"snat_addresses": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							Elem:          &schema.Schema{Type: schema.TypeString},
							Description:   "foo",
							ConflictsWith: []string{"existing_snat_pool", "automap"},
						},
					},
				},
			},
			"pool": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "foo",
						},
						"existing_pool": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"pool_members"},
						},
						"pool_members": {
							Type:          schema.TypeList,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_pool"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses": {
										Type:        schema.TypeList,
										Required:    true,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "foo",
									},
									"port": {
										Type:        schema.TypeInt,
										Required:    true,
										Computed:    true,
										Description: "foo",
									},
									"connection_limit": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "foo",
									},
									"priority_group": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "foo",
									},
									"share_nodes": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "foo",
									},
								},
							},
						},
					},
				},
			},
			"load_balancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Description: "none",
			},
			"monitor": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "foo",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "foo",
						},
						"existing_monitor": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"monitor_auth", "username", "password", "interval", "send_string", "response"},
						},
						"monitor_auth": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
						"username": {
							Type:          schema.TypeString,
							Optional:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
						"password": {
							Type:          schema.TypeString,
							Optional:      true,
							Sensitive:     true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
						"interval": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
						"send_string": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
						"response": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							Description:   "foo",
							ConflictsWith: []string{"existing_monitor"},
						},
					},
				},
			},
		},
	}
}

func resourceBigipFastHttpAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastTmpl := "bigip-fast-templates/http"
	fastJson, err := getFastHttpConfig(d)
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
		}
		tsVer := strings.Split(client.UserAgent, "/")
		err = teemDevice.Report(f, "bigip_fast_application", tsVer[3])
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
	_ = d.Set("virtual_server.ip", data.VirtualAddress)
	_ = d.Set("virtual_server.port", data.VirtualPort)
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
	_ = d.Set("monitor.existing_monitor", data.ExistingMonitor)
	_ = d.Set("monitor.monitor_auth", data.MonitorAuth)
	_ = d.Set("monitor.username", data.MonitorUsername)
	_ = d.Set("monitor.password", data.MonitorPassword)
	_ = d.Set("monitor.interval", data.MonitorInterval)
	_ = d.Set("monitor.send_string", data.MonitorSendString)
	_ = d.Set("monitor.response", data.MonitorResponse)
	return nil
}

func flattenFastPoolMembers(members []bigip.FastHttpPool) []interface{} {
	att := make([]interface{}, len(members))
	for i, v := range members {
		obj := make(map[string]interface{})

		if len(v.ServerAddresses) > 0 {
			obj["adresses"] = v.ServerAddresses
		}
		obj["port"] = v.ServicePort
		obj["connection_limit"] = v.ConnectionLimit
		obj["priority_group"] = v.PriorityGroup
		obj["share_nodes"] = v.ShareNodes
		att[i] = obj
	}
	return att
}

func getFastHttpConfig(d *schema.ResourceData) (string, error) {
	httpJson := &bigip.FastHttpJson{
		Tenant:      d.Get("tenant").(string),
		Application: d.Get("application").(string),
	}
	if v, ok := d.GetOk("virtual_server"); ok {
		virtual := v.(map[string]interface{})
		httpJson.VirtualAddress = virtual["ip"].(string)
		httpJson.VirtualPort = virtual["port"].(int)
	}
	if s, ok := d.GetOk("snat"); ok {
		snat := s.(map[string]interface{})
		if en, ok := snat["enable"].(bool); ok {
			if en == false {
				if _, ok := snat["existing_snat_pool"]; ok {
					return "", fmt.Errorf("cannot use 'existing_snat_pool' when 'enable' is: %t", en)
				}
				if _, ok := snat["snat_addresses"]; ok {
					return "", fmt.Errorf("cannot use 'snat_addresses' when 'enable' is: %t", en)
				}
				if _, ok := snat["automap"]; ok {
					return "", fmt.Errorf("cannot use 'automap' when 'enable' is: %t", en)
				}
			}
			httpJson.SnatEnable = en
		}
		if sn, ok := snat["existing_snat_pool"].(string); ok {
			httpJson.SnatPoolName = sn
		}
		if sn, ok := snat["snat_addresses"].([]string); ok {
			httpJson.SnatAddresses = sn
		}
		if sn, ok := snat["automap"].(bool); ok {
			httpJson.SnatAutomap = sn
		}
	}
	if p, ok := d.GetOk("pool"); ok {
		pool := p.(map[string]interface{})
		if en, ok := pool["enable"].(bool); ok {
			if en == false {
				if _, ok := pool["existing_pool"]; ok {
					return "", fmt.Errorf("cannot use 'existing_pool' when 'enable' is: %t", en)
				}
				if _, ok := pool["pool_members"]; ok {
					return "", fmt.Errorf("cannot use 'pool_members' when 'enable' is: %t", en)
				}
			}
			httpJson.PoolEnable = en
		}
		if pl, ok := pool["existing_pool"].(string); ok {
			httpJson.PoolName = pl
		}
		if mem, ok := pool["pool_members"].([]map[string]interface{}); ok {
			members := make([]bigip.FastHttpPool, 1)
			for _, r := range mem {
				memberConfig := bigip.FastHttpPool{}
				memberConfig.ServerAddresses = r["addresses"].([]string)
				memberConfig.ServicePort = r["port"].(int)
				if s, ok := r["connection_limit"].(int); ok {
					memberConfig.ConnectionLimit = s
				}
				if s, ok := r["priority_group"].(int); ok {
					memberConfig.PriorityGroup = s
				}
				if s, ok := r["share_nodes"].(bool); ok {
					memberConfig.ShareNodes = s
				}
				members = append(members, memberConfig)
			}
			httpJson.PoolMembers = members
		}
	}
	if p, ok := d.GetOk("load_balancing_mode"); ok {
		httpJson.LoadBalancingMode = p.(string)
	}
	if p, ok := d.GetOk("slow_ramp_time"); ok {
		httpJson.SlowRampTime = p.(int)
	}
	if m, ok := d.GetOk("monitor"); ok {
		mon := m.(map[string]interface{})
		if en, ok := mon["enable"].(bool); ok {
			if en == false {
				if _, ok := mon["existing_monitor"]; ok {
					return "", fmt.Errorf("cannot use 'existing_monitor' when 'enable' is: %t", en)
				}
				if _, ok := mon["monitor_auth"]; ok {
					return "", fmt.Errorf("cannot use 'monitor_auth' when 'enable' is: %t", en)
				}
				if _, ok := mon["username"]; ok {
					return "", fmt.Errorf("cannot use 'username' when 'enable' is: %t", en)
				}
				if _, ok := mon["password"]; ok {
					return "", fmt.Errorf("cannot use 'password' when 'enable' is: %t", en)
				}
				if _, ok := mon["interval"]; ok {
					return "", fmt.Errorf("cannot use 'interval' when 'enable' is: %t", en)
				}
				if _, ok := mon["send_string"]; ok {
					return "", fmt.Errorf("cannot use 'send_string' when 'enable' is: %t", en)
				}
				if _, ok := mon["response"]; ok {
					return "", fmt.Errorf("cannot use 'response' when 'enable' is: %t", en)
				}
			}
			httpJson.MonitorEnable = en
		}
		if e, ok := mon["existing_monitor"].(string); ok {
			httpJson.ExistingMonitor = e
		}
		if auth, ok := mon["monitor_auth"].(bool); ok {
			if auth == true {
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
	data, err := json.Marshal(httpJson)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
