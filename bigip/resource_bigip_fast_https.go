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

func resourceBigipFastHTTPSApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipFastHTTPSAppCreate,
		Read:   resourceBigipFastHTTPSAppRead,
		Update: resourceBigipFastHTTPSAppUpdate,
		Delete: resourceBigipFastHTTPSAppDelete,
		Exists: resourceBigipFastHTTPSAppExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of FAST HTTPS application tenant.",
			},
			"application": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of FAST HTTPS application.",
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
			"tls_server_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:"/Common/clientssl",
				Description:   "Select an existing TLS server profile",
				ConflictsWith: []string{"create_tls_server_profile"},
			},
			"create_tls_server_profile": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tls_cert_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Select an existing BIG-IP SSL certificate",
						},
						"tls_key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Select an existing BIG-IP SSL key",
						},
					},
				},
				ConflictsWith: []string{"tls_server_profile_name"},
			},
			"exist_pool_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Select an existing BIG-IP Pool",
				ConflictsWith: []string{"fast_create_pool_members", "existing_monitor", "fast_create_monitor"},
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
				Default:     "least-connections-member",
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

func resourceBigipFastHTTPSAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastJson, err := getFastHTTPSConfig(d)
	if err != nil {
		return err
	}
	m.Lock()
	defer m.Unlock()
	log.Printf("[INFO] Creating HTTPS FastApp config")
	log.Printf("[DEBUG]: FAST JSON:%+v", fastJson)
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
		err = teemDevice.Report(f, "bigip_fast_https_app", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipFastHttpAppRead(d, meta)
}

func resourceBigipFastHTTPSAppRead(d *schema.ResourceData, meta interface{}) error {
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
	err = setFastHTTPSData(d, fastHttp)
	if err != nil {
		return err
	}
	return nil
}

func resourceBigipFastHTTPSAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	fastJson, e := getFastHTTPSConfig(d)
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

func resourceBigipFastHTTPSAppDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceBigipFastHTTPSAppExists(d *schema.ResourceData, meta interface{}) (bool, error) {
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

func setFastHTTPSData(d *schema.ResourceData, data bigip.FastHttpJson) error {
	log.Printf("[DEBUG]: FAST JSON:%+v", data)
	_ = d.Set("tenant", data.Tenant)
	_ = d.Set("application", data.Application)
	log.Printf("[DEBUG]: VirtualAddress :%+v", data.VirtualAddress)
	if err := d.Set("virtual_server", []interface{}{flattenFastVirtualServers(data)}); err != nil {
		return fmt.Errorf("error setting virtual_server: %w", err)
	}
	_ = d.Set("existing_snat_pool", data.SnatPoolName)
	_ = d.Set("fast_create_snat_pool_address", data.SnatAddresses)
	_ = d.Set("exist_pool_name", data.PoolName)
	_ = d.Set("tls_server_profile_name", data.TlsServerProfileName)
	if _, ok := d.GetOk("create_tls_server_profile"); ok {
		if err := d.Set("create_tls_server_profile", []interface{}{flattenFastTlsServerProfile(data)}); err != nil {
			return fmt.Errorf("error setting create_tls_server_profile: %w", err)
		}
	}
	members := flattenFastPoolMembers(data.PoolMembers)
	log.Printf("[DEBUG]: Pool Members :%+v", members)
	_ = d.Set("fast_create_pool_members", members)
	_ = d.Set("load_balancing_mode", data.LoadBalancingMode)
	_ = d.Set("slow_ramp_time", data.SlowRampTime)
	_ = d.Set("existing_monitor", data.HTTPMonitor)
	if _, ok := d.GetOk("fast_create_monitor"); ok {
		if err := d.Set("fast_create_monitor", []interface{}{flattenFastMonitor(data)}); err != nil {
			return fmt.Errorf("error setting fast_create_monitor: %w", err)
		}
	}
	return nil
}

func flattenFastTlsServerProfile(data bigip.FastHttpJson) map[string]interface{} {
	tfMap := map[string]interface{}{}
	tfMap["tls_cert_name"] = data.TlsCertName
	tfMap["tls_key_name"] = data.TlsKeyName
	return tfMap
}

func getFastHTTPSConfig(d *schema.ResourceData) (string, error) {
	httpJson := &bigip.FastHttpJson{
		Tenant:      d.Get("tenant").(string),
		Application: d.Get("application").(string),
	}
	httpJson.TlsServerEnable = true
	httpJson.TlsServerProfileCreate = false

	if v, ok := d.GetOk("virtual_server"); ok {
		vL := v.([]interface{})
		for _, v := range vL {
			httpJson.VirtualAddress = v.(map[string]interface{})["ip"].(string)
			httpJson.VirtualPort = v.(map[string]interface{})["port"]
		}
	}
	if v, ok := d.GetOk("tls_server_profile_name"); ok {
		httpJson.TlsServerProfileName = v.(string)
	}
	if v, ok := d.GetOk("create_tls_server_profile"); ok {
		httpJson.TlsServerProfileCreate = true
		httpJson.TlsServerProfileName = ""
		tlsSer := v.([]interface{})
		for _, v := range tlsSer {
			httpJson.TlsCertName = v.(map[string]interface{})["tls_cert_name"].(string)
			httpJson.TlsKeyName = v.(map[string]interface{})["tls_key_name"].(string)
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
		var members []bigip.FastHttpPool
		for _, r := range p.(*schema.Set).List() {
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
