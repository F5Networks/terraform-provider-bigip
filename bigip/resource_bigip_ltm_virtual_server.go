/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var cidr = map[string]string{
	"0":  "0.0.0.0",
	"1":  "128.0.0.0",
	"2":  "192.0.0.0",
	"3":  "224.0.0.0",
	"4":  "240.0.0.0",
	"5":  "248.0.0.0",
	"6":  "252.0.0.0",
	"7":  "254.0.0.0",
	"8":  "255.0.0.0",
	"9":  "255.128.0.0",
	"10": "255.192.0.0",
	"11": "255.224.0.0",
	"12": "255.240.0.0",
	"13": "255.248.0.0",
	"14": "255.252.0.0",
	"15": "255.254.0.0",
	"16": "255.255.0.0",
	"17": "255.255.128.0",
	"18": "255.255.192.0",
	"19": "255.255.224.0",
	"20": "255.255.240.0",
	"21": "255.255.248.0",
	"22": "255.255.252.0",
	"23": "255.255.254.0",
	"24": "255.255.255.0",
	"25": "255.255.255.128",
	"26": "255.255.255.192",
	"27": "255.255.255.224",
	"28": "255.255.255.240",
	"29": "255.255.255.248",
	"30": "255.255.255.252",
	"31": "255.255.255.254",
	"32": "255.255.255.255",
}

func resourceBigipLtmVirtualServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBigipLtmVirtualServerCreate,
		ReadContext:   resourceBigipLtmVirtualServerRead,
		UpdateContext: resourceBigipLtmVirtualServerUpdate,
		DeleteContext: resourceBigipLtmVirtualServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the virtual server",
				ForceNew:     true,
				ValidateFunc: validateF5NameWithDirectory,
			},
			"port": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"trafficmatching_criteria"},
				Description:   "Listen port for the virtual server",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Source IP and mask for the virtual server",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validateEnabledDisabled,
				Description:  "Specifies whether the virtual server and its resources are available for load balancing. The default is Enabled",
			},
			"destination": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Specifies destination IP address information to which the virtual server sends traffic",
				ConflictsWith: []string{"trafficmatching_criteria"},
			},
			"trafficmatching_criteria": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Description:   "Specifies destination traffic matching information to which the virtual server sends traffic",
				ConflictsWith: []string{"destination", "port"},
			},
			"pool": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Default pool for this virtual server",
				ValidateFunc: validateF5NameWithDirectory,
			},
			"mask": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				StateFunc: func(val interface{}) string {
					if strings.Contains(val.(string), ".") {
						return val.(string)
					} else {
						return cidr[val.(string)]
					}
				},
				Description: "subnet mask",
			},
			"profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				Computed: true,
			},
			"client_profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				//Computed: true,
			},
			"server_profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				//Computed: true,
			},
			"persistence_profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				//Computed: true,
			},
			"default_persistence_profile": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fallback_persistence_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Fallback persistence profile",
			},
			"irules": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"security_log_profiles": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"per_flow_request_access_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_address_translation": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "none, automap, snat",
			},
			"snatpool": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Name of the snatpool to use. Requires source_address_translation to be set to 'snat'.",
			},
			"ip_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "tcp",
				Description: "Specifies a network protocol name you want the system to use to direct traffic on this virtual server. The default is TCP. The Protocol setting is not available when you select Performance (HTTP) as the Type.",
			},
			"policies": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Optional:    true,
				Description: "Specifies the policies for the virtual server",
			},
			"vlans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"translate_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "Specifies, when checked (enabled), that the system translates the address of the virtual server. When cleared (disabled), specifies that the system uses the address without translation. This option is useful when the system is load balancing devices that have the same IP address. The default is enabled",
			},
			"translate_port": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "enabled"}, false),
				Description:  "Specifies, when checked (enabled), that the system translates the port of the virtual server. When cleared (disabled), specifies that the system uses the port without translation. Turning off port translation for a virtual server is useful if you want to use the virtual server to load balance connections to any service. The default is enabled.",
			},
			"vlans_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables the virtual server on the VLANs specified by the VLANs option. By default it is set to false",
			},
			"firewall_enforced_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Applies the specified AFM policy to the virtual in an enforcing way,when creating a new virtual, if this parameter is not specified, the enforced is disabled.this should be in full path ex: `/Common/afm-test-policy`",
			},
		},
	}
}

func ltmVirtualServerAttrDefaults(d *schema.ResourceData) {
	_, hasMask := d.GetOk("mask")
	_, hasSource := d.GetOk("source")

	// Set default mask if nil
	if !hasMask {
		// looks like IPv6, lets set to /128
		if strings.Contains(d.Get("destination").(string), ":") {
			_ = d.Set("mask", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
		} else { // /32 for IPv4
			_ = d.Set("mask", "255.255.255.255")
		}
	}
	// set default source if nil
	if !hasSource {
		// looks like IPv6, lets set to ::/0
		if strings.Contains(d.Get("destination").(string), ":") {
			_ = d.Set("source", "::/0")
		} else { // 0.0.0.0/0
			_ = d.Set("source", "0.0.0.0/0")
		}
	}
}

func resourceBigipLtmVirtualServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating virtual server " + name)
	pss := &bigip.VirtualServer{
		Name: name,
	}
	config := getVirtualServerConfig(d, pss)
	err := client.CreateVirtualServer(config)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Virtual Server  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId(name)
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
		err = teemDevice.Report(f, "bigip_ltm_virtual_server", tsVer[3])
		if err != nil {
			log.Printf("[ERROR]Sending Telemetry data failed:%v", err)
		}
	}
	return resourceBigipLtmVirtualServerRead(ctx, d, meta)
}

func resourceBigipLtmVirtualServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Fetching virtual server " + name)

	vs, err := client.GetVirtualServer(name)
	log.Printf("[DEBUG]virtual Server Details:%+v", vs)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Server  (%s) (%v)", name, err)
		d.SetId("")
		return diag.FromErr(err)
	}
	if vs == nil {
		log.Printf("[WARN] VirtualServer (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	vsDest := vs.Destination
	log.Printf("[DEBUG]vsDest :%+v", vsDest)
	if vsDest != ":0" && strings.Count(vsDest, ":") >= 2 {
		log.Printf("[DEBUG] Matched one:%+v", vsDest)
		regex := regexp.MustCompile(`^(/.+/)(.*:[^%]*)(?:%\d+)?(?:\.(\d+))$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		if destination == nil {
			return diag.FromErr(fmt.Errorf("Unable to extract destination address and port from virtual server destination: " + vs.Destination))
		}
		if len(destination) > 3 {
			_ = d.Set("destination", destination[2])
		} else {
			_ = d.Set("destination", vsDest)
		}
	}
	if vsDest != ":0" && vsDest != "0" && strings.Count(vsDest, ":") < 2 {
		log.Printf("[DEBUG] Matched two:%+v", vsDest)
		regex := regexp.MustCompile(`(/.+/)((?:[0-9]{1,3}\.){3}[0-9]{1,3})(%\d+)?(:\d+)`)
		destination := regex.FindStringSubmatch(vs.Destination)
		if len(destination) > 3 {
			parsedDestination := destination[2] + destination[3]
			_ = d.Set("destination", parsedDestination)
		} else {
			_ = d.Set("destination", vsDest)
		}
	}

	_ = d.Set("trafficmatching_criteria", vs.TrafficMatchingCriteria)
	_ = d.Set("source", vs.Source)
	_ = d.Set("ip_protocol", vs.IPProtocol)
	_ = d.Set("name", name)
	_ = d.Set("pool", vs.Pool)
	_ = d.Set("mask", vs.Mask)

	//	/* Service port is provided by the API in the destination attribute "/partition_name/virtual_server_address[%route_domain]:(port)"
	//	   so we need to extract it
	//	*/
	//	regex = regexp.MustCompile(`\:(\d+)`)
	//	port := regex.FindStringSubmatch(vs.Destination)
	//	if len(port) < 2 {
	//		return diag.FromErr(fmt.Errorf("Unable to extract service port from virtual server destination: %s", vs.Destination)
	//	}
	//	parsedPort, _ := strconv.Atoi(port[1])

	if strings.Count(vsDest, ":") < 2 {
		regex := regexp.MustCompile(`:(\d+)`)
		port := regex.FindStringSubmatch(vs.Destination)
		log.Printf("[DEBUG] Matched for port-1:%+v", port)
		if len(port) < 2 {
			return diag.FromErr(fmt.Errorf("Unable to extract service port from virtual server destination: %s ", vs.Destination))
		}
		parsedPort, _ := strconv.Atoi(port[1])
		_ = d.Set("port", parsedPort)
	}
	if strings.Count(vsDest, ":") >= 2 {
		regex := regexp.MustCompile(`[:.](\d+)$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		if len(destination) > 1 {
			parsedPort, _ := strconv.Atoi(destination[1])
			_ = d.Set("port", parsedPort)
		}
	}
	_ = d.Set("irules", makeStringList(&vs.Rules))
	_ = d.Set("security_log_profiles", makeStringList(&vs.SecurityLogProfiles))
	_ = d.Set("per_flow_request_access_policy", vs.PerFlowRequestAccessPolicy)
	_ = d.Set("description", vs.Description)
	if vs.Enabled {
		_ = d.Set("state", "enabled")
	} else {
		_ = d.Set("state", "disabled")
	}
	_ = d.Set("source_address_translation", vs.SourceAddressTranslation.Type)

	_ = d.Set("snatpool", vs.SourceAddressTranslation.Pool)
	_ = d.Set("policies", vs.Policies)
	_ = d.Set("vlans", vs.Vlans)
	_ = d.Set("translate_address", vs.TranslateAddress)
	_ = d.Set("translate_port", vs.TranslatePort)
	_ = d.Set("firewall_enforced_policy", vs.FwEnforcedPolicy)

	profileNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(vs.PersistenceProfiles)))
	for _, profile := range vs.PersistenceProfiles {
		FullProfileName := "/" + profile.Partition + "/" + profile.Name
		profileNames.Add(FullProfileName)
	}
	if profileNames.Len() > 0 {
		if _, ok := d.GetOk("persistence_profiles"); ok {
			_ = d.Set("persistence_profiles", profileNames)
			_ = d.Set("fallback_persistence_profile", vs.FallbackPersistenceProfile)
		}
	}

	_ = d.Set("source_port", vs.SourcePort)
	_ = d.Set("vlans_enabled", vs.VlansEnabled)
	profiles, err := client.VirtualServerProfiles(name)
	if err != nil {
		return diag.FromErr(err)
	}
	if profiles != nil && len(profiles.Profiles) > 0 {
		profileNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		clientProfileNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		serverProfileNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		for _, profile := range profiles.Profiles {
			switch profile.Context {
			case bigip.CONTEXT_CLIENT:
				clientProfileNames.Add(profile.FullPath)
			case bigip.CONTEXT_SERVER:
				serverProfileNames.Add(profile.FullPath)
			default:
				profileNames.Add(profile.FullPath)
			}
		}
		if profileNames.Len() > 0 {
			_ = d.Set("profiles", profileNames)
		}
		if clientProfileNames.Len() > 0 {
			_ = d.Set("client_profiles", clientProfileNames)
		}
		if serverProfileNames.Len() > 0 {
			_ = d.Set("server_profiles", serverProfileNames)
		}
	}
	return nil
}

func resourceBigipLtmVirtualServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	pss := &bigip.VirtualServer{
		Name: name,
	}
	log.Println("[INFO] Updating virtual server " + name)
	config := getVirtualServerConfig(d, pss)
	err := client.ModifyVirtualServer(name, config)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceBigipLtmVirtualServerRead(ctx, d, meta)
}

func resourceBigipLtmVirtualServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting virtual server " + name)

	err := client.DeleteVirtualServer(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Virtual Server  (%s) (%v)", name, err)
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func getVirtualServerConfig(d *schema.ResourceData, config *bigip.VirtualServer) *bigip.VirtualServer {
	ltmVirtualServerAttrDefaults(d)
	port := d.Get("port").(int)
	mask := d.Get("mask").(string)
	destination := d.Get("destination").(string)

	config.Name = d.Get("name").(string)
	config.Pool = d.Get("pool").(string)
	config.TranslatePort = d.Get("translate_port").(string)
	config.TranslateAddress = d.Get("translate_address").(string)
	config.SourcePort = d.Get("source_port").(string)
	config.FwEnforcedPolicy = d.Get("firewall_enforced_policy").(string)
	config.Source = d.Get("source").(string)
	if strings.Contains(destination, ":") {
		subnetMask := mask
		config.Destination = fmt.Sprintf("%s.%d", destination, port)
		config.Mask = subnetMask
	} else {
		var subnetMask string
		if strings.Contains(mask, ".") {
			subnetMask = mask
		} else {
			subnetMask = cidr[mask]
		}
		config.Destination = fmt.Sprintf("%s:%d", destination, port)
		config.Mask = subnetMask
	}

	var profiles []bigip.Profile
	if p, ok := d.GetOk("profiles"); ok {
		for _, profile := range p.(*schema.Set).List() {
			profiles = append(profiles, bigip.Profile{Name: profile.(string), Context: bigip.CONTEXT_ALL})
		}
	}
	if p, ok := d.GetOk("client_profiles"); ok {
		for _, profile := range p.(*schema.Set).List() {
			profiles = append(profiles, bigip.Profile{Name: profile.(string), Context: bigip.CONTEXT_CLIENT})
		}
	}
	if p, ok := d.GetOk("server_profiles"); ok {
		for _, profile := range p.(*schema.Set).List() {
			profiles = append(profiles, bigip.Profile{Name: profile.(string), Context: bigip.CONTEXT_SERVER})
		}
	}
	var persistenceProfiles []bigip.Profile
	if p, ok := d.GetOk("persistence_profiles"); ok {
		for _, profile := range p.(*schema.Set).List() {
			if profile == d.Get("default_persistence_profile").(string) {
				persistenceProfiles = append(persistenceProfiles, bigip.Profile{Name: profile.(string), TmDefault: "yes"})
			} else {
				persistenceProfiles = append(persistenceProfiles, bigip.Profile{Name: profile.(string)})
			}
		}
		config.FallbackPersistenceProfile = d.Get("fallback_persistence_profile").(string)
	}
	var policies []string
	if p, ok := d.GetOk("policies"); ok {
		policies = setToStringSlice(p.(*schema.Set))
	}
	var vlans []string
	if v, ok := d.GetOk("vlans"); ok {
		vlans = setToStringSlice(v.(*schema.Set))
	}
	var rules []string
	if cfgRules, ok := d.GetOk("irules"); ok {
		rules = listToStringSlice(cfgRules.([]interface{}))
	}

	var securityLogProfiles []string
	if cfgLogProfiles, ok := d.GetOk("security_log_profiles"); ok {
		securityLogProfiles = listToStringSlice(cfgLogProfiles.([]interface{}))
	}

	config.SecurityLogProfiles = securityLogProfiles
	config.PerFlowRequestAccessPolicy = d.Get("per_flow_request_access_policy").(string)
	config.Description = d.Get("description").(string)
	config.Rules = rules
	config.PersistenceProfiles = persistenceProfiles
	config.Profiles = profiles
	config.Policies = policies
	config.Vlans = vlans
	config.IPProtocol = d.Get("ip_protocol").(string)
	config.TrafficMatchingCriteria = d.Get("trafficmatching_criteria").(string)
	srcAddrsTrans := struct {
		Type string `json:"type,omitempty"`
		Pool string `json:"pool,omitempty"`
	}{
		Type: d.Get("source_address_translation").(string),
		Pool: d.Get("snatpool").(string),
	}
	config.SourceAddressTranslation = srcAddrsTrans

	if d.Get("vlans_enabled").(bool) {
		config.VlansEnabled = d.Get("vlans_enabled").(bool)
	} else {
		config.VlansDisabled = true
	}
	if d.Get("state").(string) == "disabled" {
		config.Disabled = true
	}
	if d.Get("state").(string) == "enabled" {
		config.Enabled = true
	}
	return config
}
