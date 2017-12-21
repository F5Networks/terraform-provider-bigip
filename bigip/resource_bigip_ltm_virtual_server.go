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
			State: resourceBigipLtmVirtualServerImporter,
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

			"irules": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
		port,
		TranslateAddress,
		TranslatePort,
	)
	if err != nil {
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
		return err
	}

	// /Common/virtual_server_name:80
	regex := regexp.MustCompile("(/\\w+/)?([\\w._-]+)(:\\d+)?")
	destination := regex.FindStringSubmatch(vs.Destination)
	if len(destination) < 4 {
		return fmt.Errorf("Unknown virtual server destination: " + vs.Destination)
	}

	d.Set("destination", destination[2])
	d.Set("source", vs.Source)
	d.Set("protocol", vs.IPProtocol)
	d.Set("name", name)
	d.Set("pool", vs.Pool)
	d.Set("mask", vs.Mask)
	d.Set("port", vs.SourcePort)
	d.Set("irules", makeStringSet(&vs.Rules))
	d.Set("ip_protocol", vs.IPProtocol)
	d.Set("source_address_translation", vs.SourceAddressTranslation.Type)
	d.Set("snatpool", vs.SourceAddressTranslation.Pool)
	d.Set("policies", vs.Policies)
	d.Set("vlans", vs.Vlans)
	d.Set("translate_address", vs.TranslateAddress)
	d.Set("translate_port", vs.TranslatePort)

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
		rules = setToStringSlice(cfg_rules.(*schema.Set))
	}

	vs := &bigip.VirtualServer{
		Destination: fmt.Sprintf("%s:%d", d.Get("destination").(string), d.Get("port").(int)),
		Source:      d.Get("source").(string),
		Pool:        d.Get("pool").(string),
		Mask:        d.Get("mask").(string),
		Rules:       rules,
		Profiles:    profiles,
		Policies:    policies,
		Vlans:       vlans,
		IPProtocol:  d.Get("ip_protocol").(string),
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

	err := client.ModifyVirtualServer(name, vs)
	if err != nil {
		return err
	}

	return nil
}

func resourceBigipLtmVirtualServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Deleting virtual server " + name)

	return client.DeleteVirtualServer(name)
}

func resourceBigipLtmVirtualServerImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
