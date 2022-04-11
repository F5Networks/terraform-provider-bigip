/*
Original work from https://github.com/DealerDotCom/terraform-provider-bigip
Modifications Copyright 2019 F5 Networks Inc.
This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.
*/
package bigip

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	bigip "github.com/f5devcentral/go-bigip"
	"github.com/f5devcentral/go-bigip/f5teem"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigipLtmVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmVirtualServerCreate,
		Read:   resourceBigipLtmVirtualServerRead,
		Update: resourceBigipLtmVirtualServerUpdate,
		Delete: resourceBigipLtmVirtualServerDelete,
		Exists: resourceBigipLtmVirtualServerExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Listen port for the virtual server",
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
				Type:     schema.TypeString,
				Required: true,
			},

			"pool": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Default pool for this virtual server",
				ValidateFunc: validateF5Name,
			},

			"mask": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
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
				Computed:    true,
				Description: "all, tcp, udp",
			},

			"policies": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},

			"vlans": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},
			"translate_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable Address translation",
			},
			"translate_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "To enable _ disable port translation",
			},
			"vlans_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables the virtual server on the VLANs specified by the VLANs option. By default it is set to false",
			},
		},
	}
}

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

func resourceBigipLtmVirtualServerAttrDefaults(d *schema.ResourceData) {
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

func resourceBigipLtmVirtualServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	resourceBigipLtmVirtualServerAttrDefaults(d)

	name := d.Get("name").(string)
	port := d.Get("port").(int)
	mask := d.Get("mask").(string)
	pool := d.Get("pool").(string)
	destination := d.Get("destination").(string)
	//        VlansEnabled := d.Get("vlans_enabled").(bool)
	TranslateAddress := d.Get("translate_port").(string)
	TranslatePort := d.Get("translate_port").(string)

	log.Println("[INFO] Creating virtual server " + name)

	if strings.Contains(destination, ":") {
		subnetMask := mask

		config := &bigip.VirtualServer{
			Name:             name,
			Destination:      fmt.Sprintf("%s.%d", destination, port),
			Mask:             subnetMask,
			Pool:             pool,
			TranslateAddress: TranslateAddress,
			TranslatePort:    TranslatePort,
			Source:           d.Get("source").(string),
		}
		err := client.CreateVirtualServer(config)
		if err != nil {
			log.Printf("[ERROR] Unable to Create Virtual Server  (%s) (%v)", name, err)
			return err
		}

	} else {

		subnetMask := cidr[mask]

		config := &bigip.VirtualServer{
			Name:             name,
			Destination:      fmt.Sprintf("%s:%d", destination, port),
			Mask:             subnetMask,
			Pool:             pool,
			TranslateAddress: TranslateAddress,
			TranslatePort:    TranslatePort,
			Source:           d.Get("source").(string),
		}
		err := client.CreateVirtualServer(config)
		if err != nil {
			log.Printf("[ERROR] Unable to Create Virtual Server  (%s) (%v)", name, err)
			return err
		}
	}

	d.SetId(name)

	err := resourceBigipLtmVirtualServerUpdate(d, meta)
	if err != nil {
		_ = client.DeleteVirtualServer(name)
		return err
	}
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
	return resourceBigipLtmVirtualServerRead(d, meta)
}

func resourceBigipLtmVirtualServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	log.Println("[INFO] Fetching virtual server " + name)

	vs, err := client.GetVirtualServer(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Server  (%s) (%v)", name, err)
		return err
	}
	if vs == nil {
		log.Printf("[WARN] VirtualServer (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	vsDest := vs.Destination
	if strings.Count(vsDest, ":") >= 2 {
		regex := regexp.MustCompile(`^(\/.+\/)(.*:[^%]*)(?:\%\d+)?(?:\.(\d+))$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		if destination == nil {
			return fmt.Errorf("Unable to extract destination address and port from virtual server destination: " + vs.Destination)
		}
		if err := d.Set("destination", destination[2]); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Destination to state for Virtual Server  (%s): %s", d.Id(), err)
		}
	}
	if strings.Count(vsDest, ":") < 2 {
		regex := regexp.MustCompile(`(\/.+\/)((?:[0-9]{1,3}\.){3}[0-9]{1,3})(\%\d+)?(\:\d+)`)
		destination := regex.FindStringSubmatch(vs.Destination)
		parsedDestination := destination[2] + destination[3]
		if len(destination) < 3 {
			return fmt.Errorf("Unable to extract destination address from virtual server destination: " + vs.Destination)
		}
		if err := d.Set("destination", parsedDestination); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Destination to state for Virtual Server  (%s): %s", d.Id(), err)
		}
	}
	if err := d.Set("source", vs.Source); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Source to state for Virtual Server  (%s): %s", d.Id(), err)
	}

	_ = d.Set("protocol", vs.IPProtocol)
	_ = d.Set("name", name)
	if err := d.Set("pool", vs.Pool); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Pool to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	if err := d.Set("mask", vs.Mask); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Mask to state for Virtual Server  (%s): %s", d.Id(), err)
	}

	//	/* Service port is provided by the API in the destination attribute "/partition_name/virtual_server_address[%route_domain]:(port)"
	//	   so we need to extract it
	//	*/
	//	regex = regexp.MustCompile(`\:(\d+)`)
	//	port := regex.FindStringSubmatch(vs.Destination)
	//	if len(port) < 2 {
	//		return fmt.Errorf("Unable to extract service port from virtual server destination: %s", vs.Destination)
	//	}
	//	parsedPort, _ := strconv.Atoi(port[1])

	if strings.Count(vsDest, ":") < 2 {
		regex := regexp.MustCompile(`\:(\d+)`)
		port := regex.FindStringSubmatch(vs.Destination)
		if len(port) < 2 {
			return fmt.Errorf("Unable to extract service port from virtual server destination: %s", vs.Destination)
		}
		parsedPort, _ := strconv.Atoi(port[1])
		_ = d.Set("port", parsedPort)
	}
	if strings.Count(vsDest, ":") >= 2 {
		regex := regexp.MustCompile(`^(\/.+\/)(.*:[^%]*)(?:\%\d+)?(?:\.(\d+))$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		parsedPort, _ := strconv.Atoi(destination[3])
		_ = d.Set("port", parsedPort)
	}
	_ = d.Set("irules", makeStringList(&vs.Rules))
	_ = d.Set("security_log_profiles", makeStringList(&vs.SecurityLogProfiles))
	_ = d.Set("per_flow_request_access_policy", vs.PerFlowRequestAccessPolicy)
	_ = d.Set("ip_protocol", vs.IPProtocol)
	_ = d.Set("description", vs.Description)
	if vs.Enabled {
		_ = d.Set("state", "enabled")
	} else {
		_ = d.Set("state", "disabled")
	}
	_ = d.Set("source_address_translation", vs.SourceAddressTranslation.Type)
	if err := d.Set("snatpool", vs.SourceAddressTranslation.Pool); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Snatpool to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	if err := d.Set("policies", vs.Policies); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Policies to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	_ = d.Set("vlans", vs.Vlans)
	if err := d.Set("translate_address", vs.TranslateAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TranslateAddress to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	if err := d.Set("translate_port", vs.TranslatePort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TranslatePort to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	profileNames := schema.NewSet(schema.HashString, make([]interface{}, 0, len(vs.PersistenceProfiles)))
	for _, profile := range vs.PersistenceProfiles {
		FullProfileName := "/" + profile.Partition + "/" + profile.Name
		profileNames.Add(FullProfileName)
	}
	if profileNames.Len() > 0 {
		_ = d.Set("persistence_profiles", profileNames)
	}
	if err := d.Set("fallback_persistence_profile", vs.FallbackPersistenceProfile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FallbackPersistenceProfile to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	_ = d.Set("vlans_enabled", vs.VlansEnabled)
	profiles, err := client.VirtualServerProfiles(name)
	if err != nil {
		return err
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

func resourceBigipLtmVirtualServerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching virtual server " + name)

	vs, err := client.GetVirtualServer(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Retrieve Virtual Server  (%s) (%v)", name, err)
		return false, err
	}

	if vs == nil {
		d.SetId("")
	}

	return vs != nil, nil
}

func resourceBigipLtmVirtualServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	resourceBigipLtmVirtualServerAttrDefaults(d)

	name := d.Id()

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

	destPort := fmt.Sprintf("%s:%d", d.Get("destination").(string), d.Get("port").(int))
	if strings.Contains(d.Get("destination").(string), ":") {
		destPort = fmt.Sprintf("%s.%d", d.Get("destination").(string), d.Get("port").(int))
	}

	vs := &bigip.VirtualServer{
		Destination:                destPort,
		FallbackPersistenceProfile: d.Get("fallback_persistence_profile").(string),
		Source:                     d.Get("source").(string),
		Pool:                       d.Get("pool").(string),
		Mask:                       d.Get("mask").(string),
		Description:                d.Get("description").(string),
		Rules:                      rules,
		SecurityLogProfiles:        securityLogProfiles,
		PerFlowRequestAccessPolicy: d.Get("per_flow_request_access_policy").(string),
		PersistenceProfiles:        persistenceProfiles,
		Profiles:                   profiles,
		Policies:                   policies,
		Vlans:                      vlans,
		IPProtocol:                 d.Get("ip_protocol").(string),
		SourceAddressTranslation: struct {
			Type string `json:"type,omitempty"`
			Pool string `json:"pool,omitempty"`
		}{
			Type: d.Get("source_address_translation").(string),
			Pool: d.Get("snatpool").(string),
		},
		TranslatePort:    d.Get("translate_port").(string),
		TranslateAddress: d.Get("translate_address").(string),
	}
	if d.Get("vlans_enabled").(bool) {
		vs.VlansEnabled = d.Get("vlans_enabled").(bool)
	} else {
		vs.VlansDisabled = true
	}

	if d.Get("state").(string) == "disabled" {
		vs.Disabled = true
	}
	if d.Get("state").(string) == "enabled" {
		vs.Enabled = true
	}
	err := client.ModifyVirtualServer(name, vs)
	if err != nil {
		return err
	}

	return resourceBigipLtmVirtualServerRead(d, meta)
}

func resourceBigipLtmVirtualServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting virtual server " + name)

	err := client.DeleteVirtualServer(name)
	if err != nil {
		log.Printf("[ERROR] Unable to Delete Virtual Server  (%s) (%v)", name, err)
		return err
	}
	d.SetId("")
	return nil
}
