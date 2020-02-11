/*
Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"encoding/json"
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	//"reflect"
)

func dataSourceBigipAs3Service() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBigipAs3ServiceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of Service",
			},
			"virtual_addresses": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The virtual IP address you want clients to use to access resources behind the BIG-IP",
			},
			"pool_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "name of the pool",
			},
			"virtual_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "virtual server port",
			},
			"layer4": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The L4 protocol type for this virtual server",
			},
			"snat": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of built-in SNAT method or AS3 pointer to SNAT pool",
			},
			"translate_serveraddress": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true (default), make server-side connection to server address (otherwise, treat server as gateway to virtual-server address)",
			},
			"translate_serverport": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true (default), make server-side connection to server port (otherwise, connect to server on virtual-server port)",
			},
			"persistence_methods": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Default cookie is generally good. Use persistenceMethods: [] for no persistence.",
			},
			"profile_http": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "HTTP profile; name of built-in or else AS3 pointer",
			},
			"profile_tcp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TCP profile; name of built-in or else AS3 pointer",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Virtual server handles traffic only when enabled (default)",
			},
			"max_connections": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of concurrent connections you want to allow for the virtual server",
			},
			"address_status": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether the virtual server will contribute to the operational status of the associated virtual address",
			},
			"mirroring": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Controls connection-mirroring for high-availability",
			},
			"lasthop": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of built-in last-hop method or AS3 pointer to last-hop pool (default 'default' means use system setting)",
			},
			"translate_clientport": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, hide client port number from server (default false)",
			},
			"nat64enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, translate IPv6 traffic into IPv4 (default false)",
			},
			"server_tls": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "AS3 pointer to TLS Server declaration",
			},
			"service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

}
func dataSourceBigipAs3ServiceRead(d *schema.ResourceData, meta interface{}) error {
	//client := meta.(*bigip.BigIP)
	var as3service = &bigip.As3Service{}
	globalServiceType := d.Get("service_type").(string)
	if globalServiceType != "http" && globalServiceType != "https" && globalServiceType != "tcp" && globalServiceType != "udp" && globalServiceType != "l4" {
		//return errors.New(fmt.Sprintf("Incorrect Service Type"))
		return fmt.Errorf("[DEBUG] Incorrect Service Type")
	}
	if d.Get("service_type").(string) == "https" {
		as3service.Class = "Service_HTTPS"
		as3service.ServerTLS = d.Get("server_tls").(string)
	} else {
		as3service.Class = "Service_HTTP"
	}
	var virtualaddresses []string
	if m, ok := d.GetOk("virtual_addresses"); ok {
		for _, vs := range m.([]interface{}) {
			virtualaddresses = append(virtualaddresses, vs.(string))
		}
	}
	as3service.VirtualAddresses = virtualaddresses
	as3service.Pool = d.Get("pool_name").(string)
	as3service.VirtualPort = d.Get("virtual_port").(int)
	var persistencemethods []string
	if x, ok := d.GetOk("persistence_methods"); ok {
		for _, y := range x.([]interface{}) {
			persistencemethods = append(persistencemethods, y.(string))
		}
	}
	as3service.PersistenceMethods = persistencemethods
	as3service.ProfileHTTP = d.Get("profile_http").(string)
	as3service.Layer4 = d.Get("layer4").(string)
	as3service.ProfileTCP = d.Get("profile_tcp").(string)
	as3service.Enable = d.Get("enable").(bool)
	as3service.MaxConnections = d.Get("max_connections").(int)
	as3service.Snat = d.Get("snat").(string)
	as3service.Mirroring = d.Get("mirroring").(string)
	as3service.LastHop = d.Get("lasthop").(string)
	as3service.AddressStatus = d.Get("address_status").(bool)
	as3service.TranslateClientPort = d.Get("translate_clientport").(bool)
	as3service.TranslateServerAddress = d.Get("translate_serveraddress").(bool)
	as3service.TranslateServerPort = d.Get("translate_serverport").(bool)
	as3service.Nat64Enabled = d.Get("nat64enabled").(bool)
	out, err := json.Marshal(as3service)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)
	resultMap := make(map[string]interface{})
	resultMap[name] = string(out)
	resultMap["service_type"] = globalServiceType
	out1, err := json.Marshal(resultMap)
	if err != nil {
		return err
	}
	d.SetId(string(out1))
	log.Printf("[DEBUG] Service Class string :%+v", string(out1))
	return nil
}
