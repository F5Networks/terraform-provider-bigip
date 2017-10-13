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
			"address": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name/IP of the BigIP",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_HOST", nil),
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username with API access to the BigIP",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_USER", nil),
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's password",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_PASSWORD", nil),
			},
			"token_auth": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to use an external authentication source (LDAP, TACACS, etc)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_TOKEN_AUTH", nil),
			},
			"login_ref": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tmos",
				Description: "Login reference for token authentication (see BIG-IP REST docs for details)",
				DefaultFunc: schema.EnvDefaultFunc("BIGIP_LOGIN_REF", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"bigip_ltm_virtual_server":    resourceBigipLtmVirtualServer(),
			"bigip_ltm_node":              resourceBigipLtmNode(),
			"bigip_ltm_pool":              resourceBigipLtmPool(),
			"bigip_ltm_monitor":           resourceBigipLtmMonitor(),
			"bigip_ltm_irule":             resourceBigipLtmIRule(),
			"bigip_ltm_virtual_address":   resourceBigipLtmVirtualAddress(),
			"bigip_ltm_policy":            resourceBigipLtmPolicy(),
			"bigip_ltm_vlan":              resourceBigipLtmVlan(),
			"bigip_ltm_selfip":            resourceBigipLtmSelfIP(),
			"bigip_ntp":                   resourceBigipLtmNtp(),
			"bigip_dns":                   resourceBigipLtmDns(),
			"bigip_license":               resourceBigipLtmLic(),
			"bigip_license_f5bigmstbt10G": resourceBigipLtmULic(),
			"bigip_provision":             resourceBigipLtmProvision(),
			//"bigip_ltm_iapp":          resourceBigipLtmiApp(),
			"bigip_route":                resourceBigipLtmRoute(),
			"bigip_datagroup":            resourceBigipLtmDatagroup(),
			"bigip_ltm_oneconnect":       resourceBigipLtmOneconnect(),
			"bigip_syslog":               resourceBigipLtmSyslog(),
			"bigip_snmp":                 resourceBigipLtmSnmp(),
			"bigip_snmp_traps":           resourceBigipLtmSnmpTraps(),
			"bigip_tcp_profile":          resourceBigipLtmTcp(),
			"bigip_fasthttp_profile":     resourceBigipLtmFasthttp(),
			"bigip_fastl4_profile":       resourceBigipLtmFastl4(),
			"bigip_httpcompress_profile": resourceBigipLtmHttpcompress(),
			"bigip_http2_profile":        resourceBigipLtmHttp2(),
			"bigip_device_name":          resourceBigipLtmDevicename(),
			"bigip_device":               resourceBigipLtmDevice(),
			"bigip_devicegroup":          resourceBigipLtmDevicegroup(),
			"bigip_snat":                 resourceBigipLtmSnat(),
			"bigip_snatpool":             resourceBigipLtmSnatpool(),
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

//Convert slice of strings to schema.Set
func makeStringSet(list *[]string) *schema.Set {
	ilist := make([]interface{}, len(*list))
	for i, v := range *list {
		ilist[i] = v
	}
	return schema.NewSet(schema.HashString, ilist)
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
	for field, _ := range d {
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
