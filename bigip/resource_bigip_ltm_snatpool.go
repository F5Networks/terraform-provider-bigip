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
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "SNAT Pool list Name, format /partition/name. e.g. /Common/snat_pool",
				ValidateFunc: validateF5Name,
			},

			"members": {
				Type:        schema.TypeSet,
				Set:         schema.HashString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required     true,
				Description: "Specifies a translation address to add to or delete from a SNAT pool, at least one IP is required.",
			},
		},
	}
}

func resourceBigipLtmSnatpoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Get("name").(string)
	members := setToStringSlice(d.Get("members").(*schema.Set))

	log.Println("[INFO] Creating SNAT Pool " + name)

	err := client.CreateSnatPool(name, members)
	if err != nil {
		return err
	}

	d.SetId(name)

	return resourceBigipLtmSnatpoolRead(d, meta)
}

func resourceBigipLtmSnatpoolUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Updating SNAT Pool " + name)

	r := &bigip.SnatPool{
		Name:    d.Get("name").(string),
		Members: setToStringSlice(d.Get("members").(*schema.Set)),
	}

	err := client.ModifySnatPool(name, r)
	if err != nil {
		return err
	}

	return resourceBigipLtmSnatpoolRead(d, meta)
}

func resourceBigipLtmSnatpoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	log.Println("[INFO] Fetching SNAT Pool " + name)

	snatpool, err := client.GetSnatPool(name)
	if err != nil {
		return err
	}
	if snatpool == nil {
		log.Printf("[WARN] SNAT Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	d.Set("name", name)
	if err := d.Set("members", snatpool.Members); err != nil {
		return fmt.Errorf("[DEBUG] Error saving Members to state for SNAT Pool (%s): %s", d.Id(), err)
	}

	return nil

}

func resourceBigipLtmSnatpoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*bigip.BigIP)

	name := d.Id()

	err := client.DeleteSnatPool(name)
	if err == nil {
		log.Printf("[WARN] SNAT Pool (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}
