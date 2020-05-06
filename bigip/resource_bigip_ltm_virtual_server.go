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
	"regexp"
	"strconv"
	"strings"

	"github.com/f5devcentral/go-bigip"
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
				ValidateFunc: validateF5Name,
			},

			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Listen port for the virtual server",
			},

			"source": {
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "0.0.0.0/0",
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
				Type:     schema.TypeString,
				Optional: true,
				//Default:     "255.255.255.255",
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
				Computed: true,
			},

			"server_profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				Computed: true,
			},

			"persistence_profiles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				Computed: true,
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
				Computed:    true,
				Description: "Enables the virtual server on the VLANs specified by the VLANs option.",
			},
		},
	}
}

func resourceBigipLtmVirtualServerAttrDefaults(d *schema.ResourceData) {
	_, hasMask := d.GetOk("mask")
	_, hasSource := d.GetOk("source")

	// Set default mask if nil
	if !hasMask {
		// looks like IPv6, lets set to /128
		if strings.Contains(d.Get("destination").(string), ":") {
			d.Set("mask", "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff")
		} else { // /32 for IPv4
			d.Set("mask", "255.255.255.255")
		}
	}

	// set default source if nil
	if !hasSource {
		// looks like IPv6, lets set to ::/0
		if strings.Contains(d.Get("destination").(string), ":") {
			d.Set("source", "::/0")
		} else { // 0.0.0.0/0
			d.Set("source", "0.0.0.0/0")
		}
	}
}

func resourceBigipLtmVirtualServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	resourceBigipLtmVirtualServerAttrDefaults(d)

	name := d.Get("name").(string)
	port := d.Get("port").(int)
	TranslateAddress := d.Get("translate_port").(string)
	TranslatePort := d.Get("translate_port").(string)

	log.Println("[INFO] Creating virtual server " + name)
	err := client.CreateVirtualServer(
		name,
		d.Get("destination").(string),
		d.Get("mask").(string),
		d.Get("pool").(string),
		d.Get("vlans_enabled").(bool),
		port,
		TranslateAddress,
		TranslatePort,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to Create Virtual Server  (%s) (%v)", name, err)
		return err
	}

	d.SetId(name)

	err = resourceBigipLtmVirtualServerUpdate(d, meta)
	if err != nil {
		client.DeleteVirtualServer(name)
		return err
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

	vs_dest := vs.Destination
	if strings.Count(vs_dest, ":") >= 2 {
		regex := regexp.MustCompile(`^(\/.+\/)(.*:[^%]*)(?:\%\d+)?(?:\.(\d+))$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		if destination == nil {
			return fmt.Errorf("Unable to extract destination address and port from virtual server destination: " + vs.Destination)
		}
		if err := d.Set("destination", destination[2]); err != nil {
			return fmt.Errorf("[DEBUG] Error saving Destination to state for Virtual Server  (%s): %s", d.Id(), err)
		}
	}
	if strings.Count(vs_dest, ":") < 2 {
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

	d.Set("protocol", vs.IPProtocol)
	d.Set("name", name)
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

	if strings.Count(vs_dest, ":") < 2 {
		regex := regexp.MustCompile(`\:(\d+)`)
		port := regex.FindStringSubmatch(vs.Destination)
		if len(port) < 2 {
			return fmt.Errorf("Unable to extract service port from virtual server destination: %s", vs.Destination)
		}
		parsedPort, _ := strconv.Atoi(port[1])
		d.Set("port", parsedPort)
	}
	if strings.Count(vs_dest, ":") >= 2 {
		regex := regexp.MustCompile(`^(\/.+\/)(.*:[^%]*)(?:\%\d+)?(?:\.(\d+))$`)
		destination := regex.FindStringSubmatch(vs.Destination)
		parsedPort, _ := strconv.Atoi(destination[3])
		d.Set("port", parsedPort)
	}
	d.Set("irules", makeStringList(&vs.Rules))
	d.Set("ip_protocol", vs.IPProtocol)
	d.Set("description", vs.Description)
	if vs.Enabled {
		d.Set("state", "enabled")
	} else {
		d.Set("state", "disabled")
	}
	d.Set("source_address_translation", vs.SourceAddressTranslation.Type)
	if err := d.Set("snatpool", vs.SourceAddressTranslation.Pool); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Snatpool to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	if err := d.Set("policies", vs.Policies); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Policies to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	d.Set("vlans", vs.Vlans)
	if err := d.Set("translate_address", vs.TranslateAddress); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TranslateAddress to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	if err := d.Set("translate_port", vs.TranslatePort); err != nil {
		return fmt.Errorf("[DEBUG] Error saving TranslatePort to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	profile_names := schema.NewSet(schema.HashString, make([]interface{}, 0, len(vs.PersistenceProfiles)))
	for _, profile := range vs.PersistenceProfiles {
		FullProfileName := "/" + profile.Partition + "/" + profile.Name
		profile_names.Add(FullProfileName)
	}
	if profile_names.Len() > 0 {
		d.Set("persistence_profiles", profile_names)
	}
	if err := d.Set("fallback_persistence_profile", vs.FallbackPersistenceProfile); err != nil {
		return fmt.Errorf("[DEBUG] Error saving FallbackPersistenceProfile to state for Virtual Server  (%s): %s", d.Id(), err)
	}
	d.Set("vlans_enabled", vs.VlansEnabled)
	profiles, err := client.VirtualServerProfiles(name)
	if err != nil {
		return err
	}

	if profiles != nil && len(profiles.Profiles) > 0 {
		profile_names := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		client_profile_names := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		server_profile_names := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
		for _, profile := range profiles.Profiles {
			switch profile.Context {
			case bigip.CONTEXT_CLIENT:
				client_profile_names.Add(profile.FullPath)
				break
			case bigip.CONTEXT_SERVER:
				server_profile_names.Add(profile.FullPath)
				break
			default:
				profile_names.Add(profile.FullPath)
			}
		}
		if profile_names.Len() > 0 {
			d.Set("profiles", profile_names)
		}
		if client_profile_names.Len() > 0 {
			d.Set("client_profiles", client_profile_names)
		}
		if server_profile_names.Len() > 0 {
			d.Set("server_profiles", server_profile_names)
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
	if cfg_rules, ok := d.GetOk("irules"); ok {
		rules = listToStringSlice(cfg_rules.([]interface{}))
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
		VlansEnabled:     d.Get("vlans_enabled").(bool),
	}
	if d.Get("state").(string) == "disabled" {
		vs.Disabled = true
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
