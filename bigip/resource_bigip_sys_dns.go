package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceBigipSysDns() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipSysDnsCreate,
		Update: resourceBigipSysDnsUpdate,
		Read:   resourceBigipSysDnsRead,
		Delete: resourceBigipSysDnsDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipSysDnsImporter,
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

func resourceBigipSysDnsCreate(d *schema.ResourceData, meta interface{}) error {
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

	return resourceBigipSysDnsRead(d, meta)
}

func resourceBigipSysDnsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Updating DNS " + description)

	r := &bigip.DNS{
		Description:  description,
		NameServers:  setToStringSlice(d.Get("name_servers").(*schema.Set)),
		NumberOfDots: d.Get("numberof_dots").(int),
		Search:       setToStringSlice(d.Get("search").(*schema.Set)),
	}

	err := client.ModifyDNS(r)
	if err != nil {
		return err
	}
	return resourceBigipSysDnsRead(d, meta)
}

func resourceBigipSysDnsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	description := d.Id()

	log.Println("[INFO] Reading DNS " + description)

	dns, err := client.DNSs()
	if err != nil {
		return err
	}
	if dns == nil {
		log.Printf("[WARN] DNS (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("description", dns.Description)
	d.Set("name_servers", dns.NameServers)

	if err := d.Set("name_servers", dns.NameServers); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Name Servers to state for Name Servers  (%s): %s", d.Id(), err)
	}

	d.Set("numberof_dots", dns.NumberOfDots)
	if err := d.Set("numberof_dots", dns.NumberOfDots); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Numbers of dot to state for Number of Dots  (%s): %s", d.Id(), err)
	}
	d.Set("search", dns.Search)
	if err := d.Set("search", dns.Search); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Search  to state for Search  (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceBigipSysDnsDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no Delete API for this operation

	return nil
}

func resourceBigipSysDnsImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
