package bigip

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/scottdware/go-bigip"
)

func resourceBigipLtmSnat() *schema.Resource {
	log.Println("Resource schema")
	return &schema.Resource{
		Create: resourceBigipLtmSnatCreate,
		Update: resourceBigipLtmSnatUpdate,
		Read:   resourceBigipLtmSnatRead,
		Delete: resourceBigipLtmSnatDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSnatImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Snat list Name",
				//	ValidateFunc: validateF5Name,
			},
			"partition": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which partition on BIG-IP",
			},

			"autolasthop": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP autolasthop",
			},
			"mirror": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"sourceport": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"translation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"snatpool": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "BIG-IP password",
			},
			"vlansdisabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "BIG-IP password",
			},

			"origins": &schema.Schema{
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Origin IP addresses",
			},
		},
	}
}

func resourceBigipLtmSnatCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	Name := d.Get("name").(string)
	Partition := d.Get("partition").(string)
	AutoLasthop := d.Get("autolasthop").(string)
	Mirror := d.Get("mirror").(bool)
	SourcePort := d.Get("sourceport").(string)
	Translation := d.Get("translation").(string)
	Snatpool := d.Get("snatpool").(string)
	VlansDisabled := d.Get("vlansdisabled").(bool)
	Origins := setToStringSlice(d.Get("origins").(*schema.Set))
	log.Println("[INFO] Creating Snat ")

	err := client.CreateSnat(
		Name,
		Partition,
		AutoLasthop,
		SourcePort,
		Translation,
		Snatpool,
		Mirror,
		VlansDisabled,
		Origins,
	)

	if err != nil {
		return err
	}
	d.SetId(Name)
	return nil
}

func resourceBigipLtmSnatUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SNAT " + name)

	r := &bigip.Snat{
		Name:          d.Get("name").(string),
		Partition:     d.Get("partition").(string),
		AutoLasthop:   d.Get("autolasthop").(string),
		Mirror:        d.Get("mirror").(bool),
		SourcePort:    d.Get("sourceport").(string),
		Translation:   d.Get("translation").(string),
		Snatpool:      d.Get("snatpool").(string),
		VlansDisabled: d.Get("vlansdisabled").(bool),
		Origins:       setToStringSlice(d.Get("origins").(*schema.Set)),
	}

	return client.ModifySnat(r)
}

func resourceBigipLtmSnatRead(d *schema.ResourceData, meta interface{}) error {
	/*client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching snatlist " + name)

	snat, err := client.GetSnat(name)
	if err != nil {
		return err
	}
	d.Set("origins", snat.Origins)
	d.Set("name", name)
	*/
	return nil
}

func resourceBigipLtmSnatDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteSnat(name)
	//return nil
}

func resourceBigipLtmSnatImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
