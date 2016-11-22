package bigip

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmVirtualAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmVirtualAddressCreate,
		Read:   resourceBigipLtmVirtualAddressRead,
		Update: resourceBigipLtmVirtualAddressUpdate,
		Delete: resourceBigipLtmVirtualAddressDelete,
		Exists: resourceBigipLtmVirtualAddressExists,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmVirtualAddressImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the virtual address",
				ValidateFunc: validateF5Name,
			},

			"arp": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable or disable ARP for the virtual address",
				Default:     true,
			},

			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Automatically delete the virtual address with the virtual server",
			},

			"conn_limit": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Max number of connections for virtual address",
			},

			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable or disable the virtual address",
			},

			"icmp_echo": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable/Disable ICMP response to the virtual address",
			},

			"advertize_route": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enabled dynamic routing of the address",
			},

			"traffic_group": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/Common/traffic-group-1",
				Description:  "Specify the partition and traffic group",
				ValidateFunc: validateF5Name,
			},
		},
	}
}

func resourceBigipLtmVirtualAddressCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	log.Println("[INFO] Creating virtual address " + name)

	client.CreateVirtualAddress(name, hydrateVirtualAddress(d))

	d.SetId(name)
	return resourceBigipLtmVirtualAddressRead(d, meta)
}

func resourceBigipLtmVirtualAddressRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching virtual address " + name)

	var va bigip.VirtualAddress
	vas, err := client.VirtualAddresses()
	if err != nil {
		return err
	}
	for _, va = range vas.VirtualAddresses {
		if va.FullPath == name {
			break
		}
	}
	if va.FullPath != name {
		return fmt.Errorf("virtual address %s not found", name)
	}

	d.Set("name", name)
	d.Set("arp", va.ARP)
	d.Set("auto_delete", va.AutoDelete)
	d.Set("conn_limit", va.ConnectionLimit)
	d.Set("enabled", va.Enabled)
	d.Set("icmp_echo", va.ICMPEcho)
	d.Set("advertize_route", va.RouteAdvertisement)
	d.Set("traffic_group", va.TrafficGroup)

	return nil
}

func resourceBigipLtmVirtualAddressExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*bigip.BigIP)

	name := d.Id()
	log.Println("[INFO] Fetching virtual address " + name)

	var va *bigip.VirtualAddress
	vas, err := client.VirtualAddresses()
	if err != nil {
		return false, err
	}
	for _, cand := range vas.VirtualAddresses {
		if cand.FullPath == name {
			va = &cand
			break
		}
	}

	if &va == nil {
		d.SetId("")
	}

	return va != nil, nil
}

func resourceBigipLtmVirtualAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	va := hydrateVirtualAddress(d)

	err := client.ModifyVirtualAddress(name, va)
	if err != nil {
		return err
	}

	return nil
}

func hydrateVirtualAddress(d *schema.ResourceData) *bigip.VirtualAddress {
	return &bigip.VirtualAddress{
		Name:               d.Id(),
		ARP:                d.Get("arp").(bool),
		ConnectionLimit:    d.Get("conn_limit").(int),
		Enabled:            d.Get("enabled").(bool),
		ICMPEcho:           d.Get("icmp_echo").(bool),
		RouteAdvertisement: d.Get("advertize_route").(bool),
		TrafficGroup:       d.Get("traffic_group").(string),
		AutoDelete:         d.Get("auto_delete").(bool),
	}
}

func resourceBigipLtmVirtualAddressDelete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting virtual address " + name)
	client := meta.(*bigip.BigIP)
	return client.DeleteVirtualAddress(name)
}

func resourceBigipLtmVirtualAddressImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
