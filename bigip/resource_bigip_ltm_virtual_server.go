package bigip

import (
	"log"
	"regexp"
	"fmt"

	"github.com/DealerDotCom/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
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
				Type:     schema.TypeString,
				Required: true,
				Description: "Name of the virtual server",
			},

			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				Description: "Listen port for the virtual server",
			},

			"destination": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Description: "Default pool for this virtual server",
			},

			"mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default: "255.255.255.255",
				Description: "Mask can either be in CIDR notation or decimal, i.e.: \"24\" or \"255.255.255.0\". A CIDR mask of \"0\" is the same as \"0.0.0.0\"",
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
	d.Set("name", vs.Name)
	d.Set("pool", pool[len(pool) - 1])
	d.Set("mask", vs.Mask)
	d.Set("port", vs.SourcePort)

	return nil;
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

	vs := &bigip.VirtualServer{
		Destination: d.Get("destination").(string),
		Pool: d.Get("pool").(string),
		Mask: d.Get("mask").(string),
		SourcePort: fmt.Sprintf("%d", d.Get("port").(int)),
	}

	return client.ModifyVirtualServer(name, vs)
}

func resourceBigipLtmVirtualServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Deleting virtual server " + name)

	return client.DeleteVirtualServer(name)
}
