/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const DEFAULT_PARTITION = "Common"

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name/IP of the BigIP",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_HOST", nil),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username with API access to the BigIP",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's password",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_PASSWORD", nil),
			},
			"token_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_TOKEN_AUTH", nil),
			},
			"login_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tmos",
				Description: "Login reference for token authentication (see BIG-IP REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_LOGIN_REF", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"bigip_cm_device":                       resourceBigipCmDevice(),
			"bigip_cm_devicegroup":                  resourceBigipCmDevicegroup(),
			"bigip_net_route":                       resourceBigipNetRoute(),
			"bigip_net_selfip":                      resourceBigipNetSelfIP(),
			"bigip_net_vlan":                        resourceBigipNetVlan(),
			"bigip_ltm_irule":                       resourceBigipLtmIRule(),
			"bigip_ltm_datagroup":                   resourceBigipLtmDataGroup(),
			"bigip_ltm_monitor":                     resourceBigipLtmMonitor(),
			"bigip_ltm_node":                        resourceBigipLtmNode(),
			"bigip_ltm_pool":                        resourceBigipLtmPool(),
			"bigip_ltm_pool_attachment":             resourceBigipLtmPoolAttachment(),
			"bigip_ltm_policy":                      resourceBigipLtmPolicy(),
			"bigip_ltm_profile_fasthttp":            resourceBigipLtmProfileFasthttp(),
			"bigip_ltm_profile_fastl4":              resourceBigipLtmProfileFastl4(),
			"bigip_ltm_profile_http2":               resourceBigipLtmProfileHttp2(),
			"bigip_ltm_profile_httpcompress":        resourceBigipLtmProfileHttpcompress(),
			"bigip_ltm_profile_oneconnect":          resourceBigipLtmProfileOneconnect(),
			"bigip_ltm_profile_tcp":                 resourceBigipLtmProfileTcp(),
			"bigip_ltm_profile_http":                resourceBigipLtmProfileHttp(),
			"bigip_ltm_persistence_profile_srcaddr": resourceBigipLtmPersistenceProfileSrcAddr(),
			"bigip_ltm_persistence_profile_dstaddr": resourceBigipLtmPersistenceProfileDstAddr(),
			"bigip_ltm_persistence_profile_ssl":     resourceBigipLtmPersistenceProfileSSL(),
			"bigip_ltm_persistence_profile_cookie":  resourceBigipLtmPersistenceProfileCookie(),
			"bigip_ltm_profile_server_ssl":          resourceBigipLtmProfileServerSsl(),
			"bigip_ltm_profile_client_ssl":          resourceBigipLtmProfileClientSsl(),
			"bigip_ltm_snat":                        resourceBigipLtmSnat(),
			"bigip_ltm_snatpool":                    resourceBigipLtmSnatpool(),
			"bigip_ltm_virtual_address":             resourceBigipLtmVirtualAddress(),
			"bigip_ltm_virtual_server":              resourceBigipLtmVirtualServer(),
			"bigip_sys_dns":                         resourceBigipSysDns(),
			"bigip_sys_iapp":                        resourceBigipSysIapp(),
			"bigip_sys_ntp":                         resourceBigipSysNtp(),
			"bigip_sys_provision":                   resourceBigipSysProvision(),
			"bigip_sys_snmp":                        resourceBigipSysSnmp(),
			"bigip_sys_snmp_traps":                  resourceBigipSysSnmpTraps(),
			"bigip_sys_bigiplicense":                resourceBigipSysBigiplicense(),
			"bigip_as3":                             resourceBigipAs3(),
			"bigip_ssl_certificate":                 resourceBigipSslCertificate(),
			"bigip_ssl_key":                         resourceBigipSslKey(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Address:  d.Get("address").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	if d.Get("token_auth").(bool) {
		config.LoginReference = d.Get("login_ref").(string)
	}

	return config.Client()
}

//Convert slice of strings to schema.TypeSet
func makeStringList(list *[]string) []interface{} {
	ilist := make([]interface{}, len(*list))
	for i, v := range *list {
		ilist[i] = v
	}
	return ilist
}

//Convert slice of strings to schema.Set
func makeStringSet(list *[]string) *schema.Set {
	ilist := make([]interface{}, len(*list))
	for i, v := range *list {
		ilist[i] = v
	}
	return schema.NewSet(schema.HashString, ilist)
}

//Convert schema.TypeList to a slice of strings
func listToStringSlice(s []interface{}) []string {
	list := make([]string, len(s))
	for i, v := range s {
		list[i] = v.(string)
	}
	return list
}

//Convert schema.Set to a slice of strings
func setToStringSlice(s *schema.Set) []string {
	list := make([]string, s.Len())
	for i, v := range s.List() {
		list[i] = v.(string)
	}
	return list
}

//Copy map values into an object where map key == object field name (e.g. map[foo] == &{Foo: ...}
func mapEntity(d map[string]interface{}, obj interface{}) {
	val := reflect.ValueOf(obj).Elem()
	for field := range d {
		f := val.FieldByName(strings.Title(field))
		if f.IsValid() {
			if f.Kind() == reflect.Slice {
				incoming := d[field].([]interface{})
				s := reflect.MakeSlice(f.Type(), len(incoming), len(incoming))
				for i := 0; i < len(incoming); i++ {
					s.Index(i).Set(reflect.ValueOf(incoming[i]))
				}
				f.Set(s)
			} else {
				f.Set(reflect.ValueOf(d[field]))
			}
		} else {
			if field == "http_reply" {
				f := val.FieldByName(strings.Title("httpReply"))
				f.Set(reflect.ValueOf(d[field]))
			}
			log.Printf("[WARN] You probably weren't expecting %s to be an invalid field", field)
		}
	}
}

//Break a string in the format /Partition/name into a Partition / Name object
func parseF5Identifier(str string) (partition, name string) {
	if strings.HasPrefix(str, "/") {
		ary := strings.SplitN(strings.TrimPrefix(str, "/"), "/", 2)
		return ary[0], ary[1]
	}
	return "", str
}
