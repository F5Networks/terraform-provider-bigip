package bigip

import (
	"fmt"
	"log"
	"regexp"

	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmVirtualServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmVirtualServerCreate,
		Read:   resourceBigipLtmVirtualServerRead,
		Update: resourceBigipLtmVirtualServerUpdate,
		Delete: resourceBigipLtmVirtualServerDelete,
		Exists: resourceBigipLtmVirtualServerExists,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the virtual server",
			},

			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Listen port for the virtual server",
			},

			"source": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "0.0.0.0/0",
				Description: "Source IP and mask for the virtual server",
			},

			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"pool": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Default pool for this virtual server",
			},

			"mask": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "255.255.255.255",
				Description: "Mask can either be in CIDR notation or decimal, i.e.: \"24\" or \"255.255.255.0\". A CIDR mask of \"0\" is the same as \"0.0.0.0\"",
			},

			"profiles": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
				Computed: true,
			},

			"irules": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				Optional: true,
			},

			"source_address_translation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "none, automap, snat",
			},

			"ip_protocol": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "all, tcp, udp",
			},
		},
	}
}

func resourceBigipLtmVirtualServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	port := d.Get("port").(int)

	log.Println("[INFO] Creating virtual server " + name)
	err := client.CreateVirtualServer(
		name,
		d.Get("destination").(string),
		d.Get("mask").(string),
		d.Get("pool").(string),
		port,
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

	pool := strings.Split(vs.Pool, "/")
	d.Set("destination", destination[2])
	d.Set("source", vs.Source)
	d.Set("protocol", vs.IPProtocol)
	d.Set("name", vs.Name)
	d.Set("pool", pool[len(pool)-1])
	d.Set("mask", vs.Mask)
	d.Set("port", vs.SourcePort)
	d.Set("irules", makeStringSet(&vs.Rules))
	d.Set("ip_protocol", vs.IPProtocol)
	d.Set("source_address_translation", vs.SourceAddressTranslation.Type)

	profiles, err := client.VirtualServerProfiles(vs.Name)
	if err != nil {
		return err
	}
	profile_names := schema.NewSet(schema.HashString, make([]interface{}, 0, len(profiles.Profiles)))
	for _, profile := range profiles.Profiles {
		profile_names.Add(profile.Name)
	}
	d.Set("profiles", profile_names)

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
			profiles = append(profiles, bigip.Profile{Name: profile.(string)})
		}
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
		IPProtocol:  d.Get("ip_protocol").(string),
		SourceAddressTranslation: struct {
			Type string `json:"type,omitempty"`
		}{Type: d.Get("source_address_translation").(string)},
	}

	err := client.ModifyVirtualServer(name, vs)
	if err != nil {
		return err
	}

	return nil
}

func resourceBigipLtmVirtualServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Deleting virtual server " + name)

	return client.DeleteVirtualServer(name)
}
