/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"log"

	"fmt"
	"reflect"
	"strings"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

var CONTROLS = schema.NewSet(schema.HashString, []interface{}{"caching", "compression", "classification", "forwarding", "request-adaptation", "response-adpatation", "server-ssl"})
var REQUIRES = schema.NewSet(schema.HashString, []interface{}{"client-ssl", "ssl-persistence", "tcp", "http"})

func resourceBigipLtmPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmPolicyCreate,
		Read:   resourceBigipLtmPolicyRead,
		Update: resourceBigipLtmPolicyUpdate,
		Delete: resourceBigipLtmPolicyDelete,
		Exists: resourceBigipLtmPolicyExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Policy",
				ForceNew:    true,
			},
			"published_copy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Publish the Policy",
				ForceNew:    true,
			},

			"controls": {
				Type:     schema.TypeSet,
				Set:      schema.HashString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"requires": {
				Type:     schema.TypeSet,
				Set:      schema.HashString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},

			"strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "/Common/first-match",
				Description: "Policy Strategy (i.e. /Common/first-match)",
			},

			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Rule name",
						},

						"action": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"app_service": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"application": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"asm": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"avr": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cache": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"carp": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"classify": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"clone_pool": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"code": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"compress": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"content": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"cookie_hash": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cookie_insert": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cookie_passive": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cookie_rewrite": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"decompress": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"defer": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"destination_address": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"disable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"enable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"expiry": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"expiry_secs": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"extension": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"facility": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"forward": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"from_profile": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"hash": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"host": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"http": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_basic_auth": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_cookie": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_header": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_host": {
										Type:     schema.TypeBool,
										Optional: true,
										//Computed: true,
									},
									"http_referer": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_reply": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_set_cookie": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_uri": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ifile": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"insert": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"internal_virtual": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"ip_address": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"l7dos": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"length": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"location": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"log": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ltm_policy": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"member": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"tm_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"netmask": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"nexthop": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"node": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"offset": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"pem": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"persist": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"pin": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"policy": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"pool": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"priority": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"profile": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"query_string": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"rateclass": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"redirect": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"remove": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"replace": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"request": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"request_adapt": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"reset": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"response": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"response_adapt": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"scheme": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"script": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"select": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"server_ssl": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"set_variable": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"snat": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"snatpool": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"source_address": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_client_hello": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_server_handshake": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_server_hello": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_session_id": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"tcl": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"tcp_nagle": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"text": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"timeout": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"uie": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"universal": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"virtual": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vlan": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"vlan_id": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"wam": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"write": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"all": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"app_service": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"browser_type": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"browser_version": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"case_insensitive": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"case_sensitive": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cipher": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cipher_bits": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"client_ssl": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"code": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"common_name": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"contains": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"continent": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"country_code": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"country_name": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"cpu_usage": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"device_make": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"device_model": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ends_with": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"equals": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"expiry": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"extension": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"external": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"geoip": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"greater": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"greater_or_equal": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"host": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_basic_auth": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_cookie": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_header": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_host": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_method": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_referer": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_set_cookie": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_status": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_uri": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_user_agent": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"http_version": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"index": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"internal": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"isp": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"last_15secs": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"last_1min": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"last_5mins": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"less": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"less_or_equal": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"local": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"major": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"matches": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"minor": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"missing": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"mss": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"tm_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"not": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"org": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"password": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"path": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"path_segment": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"present": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"protocol": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"query_parameter": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"query_string": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"region_code": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"remote": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"request": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"response": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"route_domain": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"rtt": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"scheme": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"server_name": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_cert": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_client_hello": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_extension": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_server_handshake": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"ssl_server_hello": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"starts_with": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"tcp": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"text": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"unnamed_query_parameter": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"user_agent_token": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"username": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"version": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"vlan": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
									"vlan_id": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceBigipLtmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Get("name").(string)
	log.Println("[INFO] Creating Policy " + name)

	p := dataToPolicy(name, d)

	d.SetId(name)
	err := client.CreatePolicy(&p)
	if err != nil {
		return err
	}
	published_copy := d.Get("published_copy").(string)
	if published_copy == "" {
		published_copy = "Drafts/" + name
	}
	t := client.PublishPolicy(name, published_copy)
	if t != nil {
		return t
	}
	return resourceBigipLtmPolicyRead(d, meta)
}

func resourceBigipLtmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()

	log.Println("[INFO] Fetching policy " + name)
	p, err := client.GetPolicy(name)

	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Policy   (%s) (%v) ", name, err)
		return err
	}

	if p == nil {
		log.Printf("[WARN] Policy  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return policyToData(p, d)
}

func resourceBigipLtmPolicyExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching policy " + name)
	p, err := client.GetPolicy(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Policy   (%s) (%v) ", name, err)
		return false, err
	}
	if p == nil {
		log.Printf("[WARN] Policy  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return false, nil
	}
	return p != nil, nil
}

func resourceBigipLtmPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Updating  Policy " + name)
	p := dataToPolicy(name, d)
	err := client.UpdatePolicy(name, &p)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Policy   (%s) (%v) ", name, err)
		return err
	}
	return resourceBigipLtmPolicyRead(d, meta)
}

func resourceBigipLtmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeletePolicy(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Policy   (%s) (%v) ", name, err)
		return err
	}
	d.SetId("")
	return nil
}

func dataToPolicy(name string, d *schema.ResourceData) bigip.Policy {
	var p bigip.Policy
	values := []string{}
	values = append(values, "Drafts/")
	values = append(values, name)
	// Join three strings into one.
	result := strings.Join(values, "")
	p.Name = result
	p.Strategy = d.Get("strategy").(string)
	p.Controls = setToStringSlice(d.Get("controls").(*schema.Set))
	p.Requires = setToStringSlice(d.Get("requires").(*schema.Set))
	ruleCount := d.Get("rule.#").(int)
	p.Rules = make([]bigip.PolicyRule, 0, ruleCount)
	for i := 0; i < ruleCount; i++ {
		var r bigip.PolicyRule
		prefix := fmt.Sprintf("rule.%d", i)
		r.Name = d.Get(prefix + ".name").(string)

		actionCount := d.Get(prefix + ".action.#").(int)
		r.Actions = make([]bigip.PolicyRuleAction, actionCount, actionCount)
		for x := 0; x < actionCount; x++ {
			var a bigip.PolicyRuleAction
			mapEntity(d.Get(fmt.Sprintf("%s.action.%d", prefix, x)).(map[string]interface{}), &a)
			r.Actions[x] = a
		}

		conditionCount := d.Get(prefix + ".condition.#").(int)
		r.Conditions = make([]bigip.PolicyRuleCondition, conditionCount, conditionCount)
		for x := 0; x < conditionCount; x++ {
			var c bigip.PolicyRuleCondition
			mapEntity(d.Get(fmt.Sprintf("%s.condition.%d", prefix, x)).(map[string]interface{}), &c)
			r.Conditions[x] = c
		}
		p.Rules = append(p.Rules, r)
	}

	return p
}

func policyToData(p *bigip.Policy, d *schema.ResourceData) error {
	if err := d.Set("strategy", p.Strategy); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Strategy   state for Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("controls", makeStringSet(&p.Controls)); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Controls  state for Policy (%s): %s", d.Id(), err)
	}
	if err := d.Set("requires", makeStringSet(&p.Requires)); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Requires  state for Policy (%s): %s", d.Id(), err)
	}

	for i, r := range p.Rules {
		rule := fmt.Sprintf("rule.%d", i)
		d.Set(fmt.Sprintf("%s.name", rule), r.FullPath)
		for x, a := range r.Actions {
			action := fmt.Sprintf("%s.action.%d", rule, x)
			interfaceToResourceData(a, d, action)
		}

		for x, c := range r.Conditions {
			condition := fmt.Sprintf("%s.condition.%d", rule, x)
			interfaceToResourceData(c, d, condition)
		}
	}
	return nil
}

func interfaceToResourceData(obj interface{}, d *schema.ResourceData, prefix string) {
	v := reflect.ValueOf(obj)
	for fi := 0; fi < v.NumField(); fi++ {
		fn := v.Type().Field(fi).Name
		if fn != "Name" && fn != "Generation" {
			f := v.Field(fi)
			if (f.Kind() == reflect.Slice && f.Interface() != nil) || f.Interface() != reflect.Zero(f.Type()).Interface() {
				d.Set(fmt.Sprintf("%s.%s%s", prefix, strings.ToLower(fn[0:1]), fn[1:]), f.Interface())
			}
		}
	}
}
