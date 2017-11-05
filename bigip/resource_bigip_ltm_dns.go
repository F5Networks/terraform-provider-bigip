package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmDns() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmDnsCreate,
		Update: resourceBigipLtmDnsUpdate,
		Read:   resourceBigipLtmDnsRead,
		Delete: resourceBigipLtmDnsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmDnsImporter,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the Dns Servers",
				ValidateFunc: validateF5Name,
			},

			"name_servers": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers Address",
			},

			"numberof_dots": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "how many DNS Servers",
			},

			"search": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Servers search domain",
			},
		},
	}

}

func resourceBigipLtmDnsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Get("description").(string)
	nameservers := setToStringSlice(d.Get("name_servers").(*schema.Set))
	numberofdots := d.Get("numberof_dots").(int)
	search := setToStringSlice(d.Get("search").(*schema.Set))

	log.Println("[INFO] Creating Dns ")

	err := client.CreateDNS(
		description,
		nameservers,
		numberofdots,
		search,
	)

	if err != nil {
		return err
	}
	d.SetId(description)

	return resourceBigipLtmDnsRead(d, meta)
}

func resourceBigipLtmDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating DNS " + description)

	r := &bigip.DNS{
		Description:  description,
		NameServers:  setToStringSlice(d.Get("name_servers").(*schema.Set)),
		NumberOfDots: d.Get("numberof_dots").(int),
		Search:       setToStringSlice(d.Get("search").(*schema.Set)),
	}

	return client.ModifyDNS(r)
}

func resourceBigipLtmDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading DNS " + description)

	dns, err := client.DNSs()
	if err != nil {
		return err
	}

	d.Set("description", dns.Description)
	d.Set("name_servers", dns.NameServers)
	d.Set("numberof_dots", dns.NumberOfDots)
	d.Set("search", dns.Search)

	return nil
}

func resourceBigipLtmDnsDelete(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceBigipLtmDnsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
