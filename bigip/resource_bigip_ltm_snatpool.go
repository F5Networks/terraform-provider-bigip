package bigip

import (
	"log"

	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceBigipLtmSnatpool() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigipLtmSnatpoolCreate,
		Update: resourceBigipLtmSnatpoolUpdate,
		Read:   resourceBigipLtmSnatpoolRead,
		Delete: resourceBigipLtmSnatpoolDelete,
		Importer: &schema.ResourceImporter{
			State: resourceBigipLtmSnatpoolImporter,
		},

		Schema: map[string]*schema.Schema{

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Snatpool list Name",
				//	ValidateFunc: validateF5Name,
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Which partition on BIG-IP",
			},

			"members": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Origin IP addresses",
			},
		},
	}
}

func resourceBigipLtmSnatpoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	Name := d.Get("name").(string)
	Partition := d.Get("partition").(string)
	Members := setToStringSlice(d.Get("members").(*schema.Set))
	log.Println("[INFO] Creating Snatpool ")

	err := client.CreateSnatpool(
		Name,
		Partition,
		Members,
	)

	if err != nil {
		return err
	}
	d.SetId(Name)
	return nil
}

func resourceBigipLtmSnatpoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating Snatpool " + name)

	r := &bigip.Snatpool{
		Name:      d.Get("name").(string),
		Partition: d.Get("partition").(string),
		Members:   setToStringSlice(d.Get("members").(*schema.Set)),
	}

	return client.ModifySnatpool(r)
}

func resourceBigipLtmSnatpoolRead(d *schema.ResourceData, meta interface{}) error {
	/*client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Snatpoollist " + name)

	Snatpool, err := client.GetSnatpool(name)
	if err != nil {
		return err
	}
	d.Set("origins", Snatpool.Origins)
	d.Set("name", name)
	*/
	return nil
}

func resourceBigipLtmSnatpoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	return client.DeleteSnatpool(name)
	//return nil
}

func resourceBigipLtmSnatpoolImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
