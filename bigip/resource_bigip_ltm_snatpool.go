package bigip

import (
	"fmt"
	"github.com/f5devcentral/go-bigip"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
				Description: "Specifies a translation address to add to or delete from a SNAT pool.",
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
	return resourceBigipLtmSnatpoolRead(d, meta)
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

	err := client.ModifySnatpool(r)
	if err != nil {
		return err
	}
	return resourceBigipLtmSnatpoolRead(d, meta)
}

func resourceBigipLtmSnatpoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching Snatpoollist " + name)

	snatpool, err := client.GetSnatpool(name)
	if err != nil {
		return err
	}
	if snatpool == nil {
		log.Printf("[WARN] Snatpool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	d.Set("partition", snatpool.Partition)
	if err := d.Set("members", snatpool.Members); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Members to state for Snatpool  (%s): %s", d.Id(), err)
	}

	return nil

}

func resourceBigipLtmSnatpoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)
	name := d.Id()
	err := client.DeleteSnatpool(name)
	if err == nil {
		log.Printf("[WARN] Snat pool  (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	return nil

}

func resourceBigipLtmSnatpoolImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, nil
}
