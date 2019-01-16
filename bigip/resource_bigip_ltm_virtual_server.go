package bigip

import (
	"fmt"
	"log"
	"regexp"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
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
				ValidateFunc: validateF5Name,
			},

			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Listen port for the virtual server",
			},

			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.0.0/0",
				Description: "Source IP and mask for the virtual server",
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
				Default:     "255.255.255.255",
				Description: "Mask can either be in CIDR notation or decimal, i.e.: \"24\" or \"255.255.255.0\". A CIDR mask of \"0\" is the same as \"0.0.0.0\"",
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

func resourceBigipLtmVirtualServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

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
	// Extract destination address from "/partition_name/(virtual_server_address)[%route_domain]:port"
	regex := regexp.MustCompile(`(\/.+\/)((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?(\:\d+)`)
	destination := regex.FindStringSubmatch(vs.Destination)
	if len(destination) < 3 {
		return fmt.Errorf("Unable to extract destination address from virtual server destination: " + vs.Destination)
	}
	if err := d.Set("destination", destination[2]); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Destination to state for Virtual Server  (%s): %s", d.Id(), err)
	}

	// Extract source address from "(source_address)[%route_domain](/mask)" groups 1 + 2
	regex = regexp.MustCompile(`((?:[0-9]{1,3}\.){3}[0-9]{1,3})(?:\%\d+)?(\/\d+)`)
	source := regex.FindStringSubmatch(vs.Source)
	parsedSource := source[1] + source[2]
	if err := d.Set("source", parsedSource); err != nil {
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
	d.Set("port", vs.SourcePort)
	d.Set("irules", makeStringList(&vs.Rules))
	d.Set("ip_protocol", vs.IPProtocol)
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
	d.Set("persistence_profiles", vs.PersistenceProfiles)
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
			persistenceProfiles = append(persistenceProfiles, bigip.Profile{Name: profile.(string)})
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

	vs := &bigip.VirtualServer{
		Destination:                fmt.Sprintf("%s:%d", d.Get("destination").(string), d.Get("port").(int)),
		FallbackPersistenceProfile: d.Get("fallback_persistence_profile").(string),
		Source:                     d.Get("source").(string),
		Pool:                       d.Get("pool").(string),
		Mask:                       d.Get("mask").(string),
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
