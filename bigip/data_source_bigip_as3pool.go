/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func dataSourceBigipAs3Pool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3PoolRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Pool",
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional friendly name for this object",
			},
			"monitors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of health monitors (each by name or AS3 pointer)",
			},
			"pool_members": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_port": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  80,
						},
						"server_addresses": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of health monitors (each by name or AS3 pointer)",
						},
						"connection_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rate_limit": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dynamic_ratio": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ratio": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"priority_group": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"enable": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"adminstate": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"address_discovery": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"sharenodes": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
				Description: "Name of Application",
			},
			"loadbalancing_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Load-balancing mode default :round-robin",
			},
			"servicedown_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies connection handling when member is non-responsive,Options:“drop”,“none”,“reselect”, “reset”; Default: \"none\"",
			},
			"minimummembers_active": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Pool is down when fewer than this number of members are up",
			},
			"minimum_monitors": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Member is down when fewer than minimum monitors report it healthy. Specify ‘all’ to require all monitors to be up.",
			},
			"reselect_tries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum number of attempts to find a responsive member for a connection",
			},
			"slowramp_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "AS3 slowly the connection rate to a newly-active member slowly during this interval (seconds)",
			},
		},
	}
}

func dataSourceBigipAs3PoolRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	//idName := d.Id()
	var as3pool = &bigip.As3Pool{}
	as3pool.Class = "Pool"
	if m, ok := d.GetOk("loadbalancing_mode"); ok {
		as3pool.LoadBalancingMode = m.(string)
	}
	as3pool.ServiceDownAction = d.Get("servicedown_action").(string)
	as3pool.MinimumMembersActive = d.Get("minimummembers_active").(int)
	as3pool.ReselectTries = d.Get("reselect_tries").(int)
	as3pool.SlowRampTime = d.Get("slowramp_time").(int)
	as3pool.MinimumMonitors = d.Get("minimum_monitors").(int)
	var monitors []string
	if m, ok := d.GetOk("monitors"); ok {
		for _, vs := range m.([]interface{}) {
			monitors = append(monitors, vs.(string))
		}
	}
	as3pool.Monitors = monitors
	var poolmemberList []bigip.As3PoolMember
	if m, ok := d.GetOk("pool_members"); ok {
		var as3poolMember = bigip.As3PoolMember{}
		//log.Printf("m Struct:%+v\n",reflect.TypeOf(m.(*schema.Set).List()))
		for _, v := range m.(*schema.Set).List() {
			//log.Printf("Map Result:%+v\n",v.(map[string]interface{}))
			as3poolMember.ConnectionLimit = v.(map[string]interface{})["connection_limit"].(int)
			as3poolMember.RateLimit = v.(map[string]interface{})["rate_limit"].(int)
			as3poolMember.DynamicRatio = v.(map[string]interface{})["dynamic_ratio"].(int)
			as3poolMember.ServicePort = v.(map[string]interface{})["service_port"].(int)
			as3poolMember.Ratio = v.(map[string]interface{})["ratio"].(int)
			as3poolMember.PriorityGroup = v.(map[string]interface{})["priority_group"].(int)
			as3poolMember.ShareNodes = v.(map[string]interface{})["sharenodes"].(bool)
			as3poolMember.ServerAddresses = listToStringSlice(v.(map[string]interface{})["server_addresses"].([]interface{}))
			as3poolMember.AdminState = v.(map[string]interface{})["adminstate"].(string)
			as3poolMember.AddressDiscovery = v.(map[string]interface{})["address_discovery"].(string)
		}
		//log.Printf("Inside dataSourceBigipAs3PoolRead : as3poolMember Struct:%+v\n", as3poolMember)
		poolmemberList = append(poolmemberList, as3poolMember)
	}
	as3pool.Members = poolmemberList
	out, err := json.Marshal(as3pool)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	out1, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}
	d.SetId(string(out1))
	log.Printf("Result Map:%+v", d.Get("result_map"))
	return nil
}
